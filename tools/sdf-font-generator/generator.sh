#!/bin/bash

# --- CONFIGURATION ---
RESULT_DIR="./results"
PACKAGE="msdf-bmfont-xml"
CHARSET_FILE="$(pwd)/.tmp_charset.txt"

echo "--- MSDF Font Generator (Auto-Detect All Characters) ---"

# 1. Global Dependency Check
if ! npm list -g $PACKAGE --depth=0 &> /dev/null; then
    echo "📦 Package '$PACKAGE' missing. Installing..."
    sudo npm install -g $PACKAGE || { echo "❌ Install failed"; exit 1; }
fi

if ! command -v fc-query &> /dev/null; then
    echo "❌ Error: 'fontconfig' is not installed. Please install it to auto-detect characters."
    exit 1
fi

# 2. Input Handling
read -p "Drag & Drop TTF file: " FONT_PATH_RAW
# Clean quotes and whitespace
FONT_PATH=$(echo "$FONT_PATH_RAW" | sed "s/['\"]//g" | xargs)

if [ ! -f "$FONT_PATH" ]; then
    echo "❌ Error: Cannot find file at '$FONT_PATH'"
else
    read -p "Enter font size [48]: " FONT_SIZE
    FONT_SIZE=${FONT_SIZE:-48}
    FONT_NAME=$(basename "$FONT_PATH" | cut -d. -f1)
    
    # 3. Create Folder
    mkdir -p "$RESULT_DIR"
    
    # 4. Extract All Characters from Font
    echo "🔍 Scanning font for all available glyphs..."
    
    # This magic line extracts the hex codes from the font and converts them to UTF-8
    # It filters out control characters and ensures a clean string
    ALL_CHARS=$(fc-query --format='%{charset}\n' "$FONT_PATH" | \
                grep -v '^$' | \
                python3 -c "import sys; 
ranges = sys.stdin.read().split(); 
chars = ''; 
for r in ranges:
    if '-' in r:
        start, end = r.split('-')
        for i in range(int(start, 16), int(end, 16) + 1):
            chars += chr(i)
    else:
        chars += chr(int(r, 16))
print(chars)" 2>/dev/null)

    if [ -z "$ALL_CHARS" ]; then
        echo "❌ Error: Could not extract characters from font."
        exit 1
    fi

    # Save to a temporary file (passing massive strings as arguments can fail)
    echo -n "$ALL_CHARS" > "$CHARSET_FILE"
    
    echo "⏳ Generating MSDF Atlas for $(echo -n "$ALL_CHARS" | wc -m) characters..."

    # 5. Execution
    # -i: Uses the charset file we just created
    # -m: Increased texture size (needed if font has many characters)
    npx -g $PACKAGE \
      -t psdf \
      -f xml \
      -s "$FONT_SIZE" \
      -i "$CHARSET_FILE" \
      -p 0 \
      -b 1 \
      -o "$RESULT_DIR/$FONT_NAME" \
      "$FONT_PATH"
    
    # 6. Rename .fnt to .xml
    if [ -f "$RESULT_DIR/$FONT_NAME.fnt" ]; then
        mv "$RESULT_DIR/$FONT_NAME.fnt" "$RESULT_DIR/$FONT_NAME.xml"
    fi

    # 7. Cleanup
    [ -f "$CHARSET_FILE" ] && rm "$CHARSET_FILE"
    
    # Fix permissions for the result folder
    if [ -d "$RESULT_DIR" ]; then
        chmod -R 666 "$RESULT_DIR"/* 2>/dev/null
    fi

    if [ -f "$RESULT_DIR/$FONT_NAME.xml" ]; then
        echo "--------------------------------------"
        echo "✅ Success!"
        echo "Format: XML"
        echo "Glyphs: $(echo -n "$ALL_CHARS" | wc -m) characters included."
        echo "Files: $RESULT_DIR/$FONT_NAME.xml, $RESULT_DIR/$FONT_NAME.png"
    else
        echo "--------------------------------------"
        echo "❌ Error: Generation failed. Check if the font is valid."
    fi
fi

echo ""
read -n 1 -s -r -p "Press any key to close..."
echo ""