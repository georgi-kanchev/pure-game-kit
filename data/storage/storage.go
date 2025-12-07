// Wraps the JSON & XML formats, as well as some byte compression/decompression functionalities to make
// them more digestible and clarify their API.
package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"pure-game-kit/debug"
)

func FromFileJSON(path string, structInstance any) {
	var file, err = os.Open(path)
	if err != nil {
		debug.LogError("Failed to open JSON file: \"", path, "\"\n", err)
		return
	}
	defer file.Close()

	var decoder = json.NewDecoder(file)
	var err2 = decoder.Decode(structInstance)
	if err2 != nil {
		debug.LogError("Failed to decode JSON file: \"", path, "\"\n", err2)
	}
}
func FromFileXML(path string, structInstance any) {
	var file, err = os.Open(path)
	if err != nil {
		debug.LogError("Failed to open XML file: \"", path, "\"\n", err)
		return
	}
	defer file.Close()

	var decoder = xml.NewDecoder(file)
	var err2 = decoder.Decode(structInstance)
	if err2 != nil {
		debug.LogError("Failed to decode XML file: \"", path, "\"\n", err2)
	}
}

func FromJSON(jsonData string, structInstance any) {
	var err = json.Unmarshal([]byte(jsonData), structInstance)
	if err != nil {
		debug.LogError("Failed to populate struct instance with JSON data!\n", err)
	}
}
func FromXML(xmlData string, structInstance any) {
	var err = xml.Unmarshal([]byte(xmlData), structInstance)
	if err != nil {
		debug.LogError("Failed to populate struct instance with XML data!\n", err)
	}
}

func ToJSON(structPointer any) string {
	var data, err = json.MarshalIndent(structPointer, "", "  ") // pretty print
	if err != nil {
		debug.LogError("Failed to create JSON data from struct instance!\n", err)
		return ""
	}
	return string(data)
}
func ToXML(structPointer any) string {
	var data, err = xml.MarshalIndent(structPointer, "", "  ") // pretty print
	if err != nil {
		debug.LogError("Failed to create XML data from struct instance!\n", err)
		return ""
	}
	return string(data)
}

func CompressZLIB(data []byte) []byte {
	var buf bytes.Buffer
	var gw = zlib.NewWriter(&buf)
	var _, err = gw.Write(data)
	var err2 = gw.Close()

	if err != nil || err2 != nil {
		debug.LogError("Failed to compress data with ZLIB!\n", err)
		return data
	}
	return buf.Bytes()
}
func DecompressZLIB(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = zlib.NewReader(buf)
	var result, err2 = io.ReadAll(gr)
	var err3 = gr.Close()

	if err != nil || err2 != nil || err3 != nil {
		debug.LogError("Failed to decompress data with ZLIB!\n", err)
		return data
	}
	return result

}
func CompressGZIP(data []byte) []byte {
	var buf bytes.Buffer
	var gw = gzip.NewWriter(&buf)
	var _, err = gw.Write(data)
	var err2 = gw.Close()

	if err != nil || err2 != nil {
		debug.LogError("Failed to compress data with GZIP!\n", err)
		return data
	}
	return buf.Bytes()
}
func DecompressGZIP(data []byte) []byte {
	var buf = bytes.NewReader(data)
	var gr, err = gzip.NewReader(buf)
	var result, err2 = io.ReadAll(gr)
	var err3 = gr.Close()

	if err != nil || err2 != nil || err3 != nil {
		debug.LogError("Failed to compress data with GZIP!\n", err)
		return data
	}
	return result

}
