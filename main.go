package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type KeyHole struct {
	Global    map[interface{}]interface{} `yaml:"global,omitempty"`
	Services  map[interface{}]interface{} `yaml:"services,omitempty"`
	KeyValues map[string]string           `yaml:"keyvalues,omitempty"`
}

func main() {
	m := &KeyHole{}

	yamlData, err := ioutil.ReadFile("./test.yml")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	err = yaml.Unmarshal(yamlData, m)
	m.KeyValues = make(map[string]string)
	m.GenerateKeyValues()
	err = m.PrintKeys()
}

func (kh *KeyHole) GenerateKeyValues() (map[string]string, error) {
	keys := make(map[string]string)

	for key := range kh.Global {
		kh.traverse(strings.Join([]string{"", key.(string)}, "/"), kh.Global[key])
	}

	return keys, nil
}

func (kh *KeyHole) traverse(root string, data interface{}) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		for _, k := range traverseMapForKeys(data.(map[interface{}]interface{})) {
			kh.traverse(strings.Join([]string{root, k}, "/"), data.(map[interface{}]interface{})[k])
		}
	case reflect.Slice:
		for i, item := range data.([]interface{}) {
			kh.KeyValues[strings.Join([]string{root, strconv.Itoa(i)}, "/")] = item.(string)
			//fmt.Printf("Slice Root: %v\n", strings.Join([]string{root, strconv.Itoa(i)}, "/"))
			//fmt.Printf("Item: %v\n", item)
		}
	default:
		fmt.Printf("default: %s\n", reflect.TypeOf(data).Kind())
	}
}

func (kh *KeyHole) PrintKeys() error {
	for key := range kh.KeyValues {
		fmt.Printf("%s\n", strings.Join([]string{key, kh.KeyValues[key]}, " = "))
	}
	return nil
}

func traverseMapForKeys(data map[interface{}]interface{}) []string {
	rData := []string{}
	for key := range data {
		rData = append(rData, key.(string))
	}
	return rData
}
