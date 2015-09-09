package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"gopkg.in/yaml.v2"
)

type KeySpace struct {
	RawData   map[interface{}]interface{}
	KeyValues map[string]string
}

func (kh *KeySpace) sendToEtcd() error {
	machines := []string{"http://192.168.99.100:2379"}
	client := etcd.NewClient(machines)

	for key, value := range kh.KeyValues {
		if _, err := client.Set(key, value, 0); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	keySpace := &KeySpace{}

	yamlData, err := ioutil.ReadFile("./test.yml")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	keySpace.RawData = make(map[interface{}]interface{})
	keySpace.KeyValues = make(map[string]string)
	err = yaml.Unmarshal(yamlData, &keySpace.RawData)

	keySpace.GenerateKeyValues()
	err = keySpace.PrintKeys()
	err = keySpace.sendToEtcd()
}

func (kh *KeySpace) GenerateKeyValues() (map[string]string, error) {
	keys := make(map[string]string)

	for key := range kh.RawData {
		kh.traverse(strings.Join([]string{"", key.(string)}, "/"), kh.RawData[key])
	}

	return keys, nil
}

func (kh *KeySpace) traverse(root string, data interface{}) {
	switch data.(type) {
	case map[interface{}]interface{}:
		for _, k := range traverseMapForKeys(data.(map[interface{}]interface{})) {
			kh.traverse(strings.Join([]string{root, k}, "/"), data.(map[interface{}]interface{})[k])
		}
	case []interface{}:
		for i, item := range data.([]interface{}) {
			kh.traverse(strings.Join([]string{root, strconv.Itoa(i)}, "/"), item)
		}
	case string:
		kh.KeyValues[root] = data.(string)
	case int:
		kh.KeyValues[root] = strconv.Itoa(data.(int))
	default:
		fmt.Printf("default: %s\n", reflect.TypeOf(data))
	}
}

func (kh *KeySpace) PrintKeys() error {
	for key := range kh.KeyValues {
		fmt.Printf("%s\n", strings.Join([]string{key, kh.KeyValues[key]}, "="))
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
