package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func FromFileJSON(path string, structInstance any) {
	var file, err = os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open JSON file: %v\n", err)
		return
	}
	defer file.Close()

	var decoder = json.NewDecoder(file)
	var err2 = decoder.Decode(structInstance)
	if err2 != nil {
		fmt.Printf("Failed to decode JSON file: %v\n", err2)
	}
}
func FromFileXML(path string, structInstance any) {
	var file, err = os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open XML file: %v\n", err)
		return
	}
	defer file.Close()

	var decoder = xml.NewDecoder(file)
	var err2 = decoder.Decode(structInstance)
	if err2 != nil {
		fmt.Printf("Failed to decode XML file: %v\n", err2)
	}
}

func FromJSON(jsonData string, structInstance any) {
	var err = json.Unmarshal([]byte(jsonData), structInstance)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v\n", err)
	}
}
func FromXML(xmlData string, structInstance any) {
	var err = xml.Unmarshal([]byte(xmlData), structInstance)
	if err != nil {
		fmt.Printf("Failed to unmarshal XML: %v\n", err)
	}
}

func ToJSON(structPointer any) string {
	var data, err = json.MarshalIndent(structPointer, "", "  ") // pretty print
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return ""
	}
	return string(data)
}
func ToXML(structPointer any) string {
	var data, err = xml.MarshalIndent(structPointer, "", "  ") // pretty print
	if err != nil {
		fmt.Printf("Failed to marshal XML: %v\n", err)
		return ""
	}
	return string(data)
}

func CompressZLIB(data []byte) []byte {
	var buf bytes.Buffer
	var gw = zlib.NewWriter(&buf)
	var _, err = gw.Write(data)

	if err != nil {
		return data
	}

	if err := gw.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}
func DecompressZLIB(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = zlib.NewReader(buf)

	if err != nil {
		return data
	}
	defer gr.Close()

	var result, err2 = io.ReadAll(gr)
	if err2 != nil {
		return data
	}
	return result

}
func CompressGZIP(data []byte) []byte {
	var buf bytes.Buffer
	var gw = gzip.NewWriter(&buf)
	var _, err = gw.Write(data)

	if err != nil {
		return data
	}

	if err := gw.Close(); err != nil {
		return data
	}
	return buf.Bytes()
}
func DecompressGZIP(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = gzip.NewReader(buf)

	if err != nil {
		return data
	}
	defer gr.Close()

	var result, err2 = io.ReadAll(gr)
	if err2 != nil {
		return data
	}
	return result

}
