/*
Wraps some well-known text data formats, as well as a binary format.
Also contains a few binary compression/decompression solutions.

Supported data formats:

	| Feature             | JSON     | XML       | YAML      | Binary           |
	|--------=------------|----------|-----------|-----------|------------------|
	| Primary Use         | Config   | Config    | Config    | Game State       |
	| Data Size           | Medium   | High      | Medium    | Very Low         |
	| Parsing Speed       | Medium   | Slow      | Slow      | Fast             |
	| Data Types          | Basic    | Basic     | Moderate  | All              |
	| Struct Inheritance  | No       | No        | No        | Yes              |
	| Mergability (Git)   | Good     | Moderate  | Excellent | Bad              |
	| Breakability        | High     | Medium    | Low       | Critical         |
	| Editability         | High     | Medium    | Very High | None             |
	| Text Features ------------------------------------------------------------|
	| Nesting Quality     | Moderate | Very Good | Excellent |
	| Layout Readability  | Vertical | Ver + Hor | Vertical  |
	| Overall Readability | High     | Medium    | Very High |
	| Syntax Style        | Braces   | Tags      | Indent    |
	| Comments Support    | No       | Yes       | Yes       |
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
	"pure-game-kit/debug"

	"gopkg.in/yaml.v3"
)

func FromJSON(jsonData string, objectPointer any) {
	var err = json.Unmarshal([]byte(jsonData), objectPointer)
	if err != nil {
		debug.LogError("Failed to populate object with JSON data!\n", err)
	}
}
func FromXML(xmlData string, objectPointer any) {
	var err = xml.Unmarshal([]byte(xmlData), objectPointer)
	if err != nil {
		debug.LogError("Failed to populate object with XML data!\n", err)
	}
}
func FromYAML(yamlData string, objectPointer any) {
	var err = yaml.Unmarshal([]byte(yamlData), objectPointer)
	if err != nil {
		debug.LogError("Failed to populate object with YAML data!\n", err)
	}
}

/*
Populates an object (struct, slice, map) from bytes produced by:

	storage.ToBytes(...)

To populate the object correctly, optional types are required when the data (that was converted to bytes)
contains interfaces. Does not work when interfaces with the same names are provided from different packages.
Useful for loading from a file. See package description for features.
*/
func FromBytes(data []byte, objectPointer any, registerTypes ...any) {
	for _, t := range registerTypes {
		gob.Register(t)
	}

	const msg = "Failed to populate object from binary data!\n"
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

func ToJSON(objectPointer any) string {
	var data, err = json.MarshalIndent(objectPointer, "", "  ") // pretty print
	if err != nil {
		debug.LogError("Failed to create JSON data from object!\n", err)
		return ""
	}
	return string(data)
}
func ToXML(objectPointer any) string {
	var data, err = xml.MarshalIndent(objectPointer, "", "  ") // pretty print
	if err != nil {
		debug.LogError("Failed to create XML data from object!\n", err)
		return ""
	}
	return string(data)
}
func ToYAML(objectPointer any) string {
	var data, err = yaml.Marshal(objectPointer)
	if err != nil {
		debug.LogError("Failed to create YAML data from object!\n", err)
		return ""
	}
	return string(data)
}

/*
Converts any object (struct, slice, map) into bytes that can be used later to populate the object back with:

	storage.FromBytes(...)

To register interface implementations, optional types are required when the object contains interfaces.
Does not work when interfaces with the same names are provided from different packages.
Useful for saving to a file. See package description for features.
*/
func ToBytes(objectPointer any, registerTypes ...any) []byte {
	for _, t := range registerTypes {
		gob.Register(t)
	}

	var buf bytes.Buffer
	var err = gob.NewEncoder(&buf).Encode(objectPointer)
	if err != nil {
		debug.LogError("Failed to convert object into binary data!\n", err)
		return nil
	}
	return CompressZLIB(buf.Bytes())
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
