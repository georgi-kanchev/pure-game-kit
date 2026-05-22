#!/bin/bash

# Repository 			https://github.com/Chlumsky/msdf-atlas-gen
# Known kerning issue 	https://github.com/Chlumsky/msdf-atlas-gen/issues/4#issuecomment-792912921

# 1. Ask the user for the font file path
echo "=== MSDF Atlas Gen (Wine Wrapper) ==="
read -p "Drag & drop the .ttf file: " INPUT_PATH

# Remove surrounding quotes if the user dragged and dropped the file into the terminal
INPUT_PATH=$(echo "$INPUT_PATH" | sed -e 's/^["'\'']//' -e 's/["'\'']$//')

# 2. Verify the file actually exists
if [ ! -f "$INPUT_PATH" ]; then
    echo "Error: File '$INPUT_PATH' not found. Please check the path and try again."
    echo ""
    read -n 1 -s -r -p "Press any key to exit..."
    echo ""
    exit 1
fi

# 3. Ask for Font Size with a default of 16
read -p "Enter font glyph size [Default: 16]: " INPUT_SIZE
SIZE=${INPUT_SIZE:-16}

# 4. Ask for Smoothness/Pixel Range with a default of 4
read -p "Enter pixel range / smoothness [Default: 4]: " INPUT_RANGE
RANGE=${INPUT_RANGE:-4}

# 5. Ask for Charset Directory
read -p "Enter path to charset directory: " CHARSET_DIR
CHARSET_DIR=$(echo "$CHARSET_DIR" | sed -e 's/^["'\'']//' -e 's/["'\'']$//')

# --- NEW STEP: Extract Unique Characters ---
if [ ! -d "$CHARSET_DIR" ]; then
    echo "Error: Directory '$CHARSET_DIR' not found."
    exit 1
fi

echo "Scanning '$CHARSET_DIR' recursively for unique characters..."

# Gather the unique raw strings into a temporary internal block
RAW_CONTENT=$(find "$CHARSET_DIR" -type f -name "*.txt" -exec cat {} + | LC_ALL=en_US.UTF-8 perl -CS -gpe 's/(.)/$1\n/g' | sort -u)

# 1. Extract raw literal characters for the console display
CHARS_STRING=$(echo "$RAW_CONTENT" | LC_ALL=en_US.UTF-8 perl -CS -ne '
    chomp;
    next if $_ eq "" or $_ eq "\n" or $_ eq "\r" or $_ eq " ";
    print "$_ ";
')

# 2. Extract numeric hexadecimal codes to feed to msdf-atlas-gen safely
HEX_STRING=$(echo "$RAW_CONTENT" | LC_ALL=en_US.UTF-8 perl -CS -ne '
    chomp;
    next if $_ eq "" or $_ eq "\n" or $_ eq "\r" or $_ eq " ";
    my $hex = sprintf("0x%x", ord($_));
    print "$hex ";
')

# Check if we actually found any characters
if [ -z "$HEX_STRING" ]; then
    echo "Warning: No characters extracted. Defaulting to standard ASCII alphanumeric."
    CHARS_STRING="A B C"
    HEX_STRING="0x41 0x42 0x43"
fi
# -------------------------------------------

# 6. Extract the base filename without the path and extension
FONT_NAME=$(basename "$INPUT_PATH" | sed 's/\.[^.]*$//')

echo ""
echo "Processing font: $FONT_NAME"
echo "Glyph Size: $SIZE px"
echo "Pixel Range: $RANGE px"
echo "Unique characters found:"
echo "$CHARS_STRING"
echo "Running msdf-atlas-gen via Wine..."
echo "----------------------------------------"

# 7. Execute the command with the dynamic arguments
# The tool receives the HEX values string, preventing Wine string translation problems completely
WINEDEBUG=-all LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8 wine ./msdf-atlas-gen.exe \
  -font "$INPUT_PATH" \
  -type msdf \
  -format png \
  -imageout "${FONT_NAME}.png" \
  -json "${FONT_NAME}.json" \
  -pxrange "$RANGE" \
  -size "$SIZE" \
  -chars "$HEX_STRING" \
  -yorigin "top"

# 8. Check if Wine executed the tool successfully
echo "----------------------------------------"
if [ $? -eq 0 ]; then
    echo "Success! Created files:"
    echo " - ${FONT_NAME}.png"
    echo " - ${FONT_NAME}.json"
else
    echo "Something went wrong while running the executable through Wine."
fi

# 9. Keep terminal open until a key is pressed
echo ""
read -n 1 -s -r -p "Press any key to exit..."
echo ""