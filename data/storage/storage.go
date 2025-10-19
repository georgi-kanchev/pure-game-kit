package storage

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"pure-game-kit/internal"
)

func FromFileJSON(path string, structInstance any) {
	path = internal.MakeAbsolutePath(path)
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
	path = internal.MakeAbsolutePath(path)
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
