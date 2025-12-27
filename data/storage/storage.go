/*
Wraps the JSON & XML formats, as well as some byte convertion/compression/decompression
functionalities to make them more digestible and clarify their API.
*/
package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"pure-game-kit/debug"
)

func FromFileJSON(path string, objectPointer any) {
	var file, err = os.Open(path)
	if err != nil {
		debug.LogError("Failed to open JSON file: \"", path, "\"\n", err)
		return
	}
	defer file.Close()

	var decoder = json.NewDecoder(file)
	var err2 = decoder.Decode(objectPointer)
	if err2 != nil {
		debug.LogError("Failed to decode JSON file: \"", path, "\"\n", err2)
	}
}
func FromFileXML(path string, objectPointer any) {
	var file, err = os.Open(path)
	if err != nil {
		debug.LogError("Failed to open XML file: \"", path, "\"\n", err)
		return
	}
	defer file.Close()

	var decoder = xml.NewDecoder(file)
	var err2 = decoder.Decode(objectPointer)
	if err2 != nil {
		debug.LogError("Failed to decode XML file: \"", path, "\"\n", err2)
	}
}

func FromJSON(jsonData string, objectPointer any) {
	var err = json.Unmarshal([]byte(jsonData), objectPointer)
	if err != nil {
		debug.LogError("Failed to populate struct instance with JSON data!\n", err)
	}
}
func FromXML(xmlData string, objectPointer any) {
	var err = xml.Unmarshal([]byte(xmlData), objectPointer)
	if err != nil {
		debug.LogError("Failed to populate struct instance with XML data!\n", err)
	}
}

func ToJSON(objectPointer any) string {
	var data, err = json.MarshalIndent(objectPointer, "", "  ") // pretty print
	if err != nil {
		debug.LogError("Failed to create JSON data from struct instance!\n", err)
		return ""
	}
	return string(data)
}
func ToXML(objectPointer any) string {
	var data, err = xml.MarshalIndent(objectPointer, "", "  ") // pretty print
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

/*
Converts any object (struct, slice, map) into bytes that can be used later to populate the object back with:

	storage.FromBytes(...)

To register interface implementations, optional types are required when the object contains interfaces.
Does not work when interfaces with the same names are provided from different packages.
Useful for saving to a file.
*/
func ToBytes(objectPointer any, registerTypes ...any) []byte {
	for _, t := range registerTypes {
		gob.Register(t)
	}

	var buf bytes.Buffer
	var err = gob.NewEncoder(&buf).Encode(objectPointer)
	if err != nil {
		debug.LogError("Failed to convert struct instance into binary data!\n", err)
		return nil
	}
	return CompressZLIB(buf.Bytes())
}

/*
Populates an object (struct, slice, map) from bytes produced by:

	storage.ToBytes(...)

To populate the object correctly, optional types are required when the data (that was converted to bytes)
contains interfaces. Does not work when interfaces with the same names are provided from different packages.
Useful for loading from a file.
*/
func FromBytes(data []byte, objectPointer any, registerTypes ...any) {
	for _, t := range registerTypes {
		gob.Register(t)
	}

	const msg = "Failed to populate struct instance from binary data!\n"
	if len(data) == 0 {
		debug.LogError(msg, "Bytes data is empty.")
		return
	}
	var buf = bytes.NewBuffer(DecompressZLIB(data))
	var dec = gob.NewDecoder(buf)
	var err = dec.Decode(objectPointer)
	if err != nil {
		debug.LogError(msg, err)
	}
}
