/*
* Copyright (c) 2025 System233
*
* This software is released under the MIT License.
* https://opensource.org/licenses/MIT
 */
package utils

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadYamlData(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
func LoadYamlFile(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return LoadYamlData(data, v)
}
func DumpYamlData(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
func DumpYamlFile(file string, v interface{}) error {
	data, err := DumpYamlData(v)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0755)
}

func LoadJsonData(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
func LoadJsonFile(file string, v interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return LoadJsonData(data, v)
}
func DumpJsonData(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func DumpJsonFile(file string, v interface{}) error {
	data, err := DumpJsonData(v)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0755)
}
