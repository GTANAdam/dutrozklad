// Package util ..
package util

import (
	"dutrozkladbot/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

// UpcaseInitial ..
func UpcaseInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+2:]
	}
	return ""
}

// Mapkey ..
func Mapkey(m map[int]interface{}, value string) string {
	for k, v := range m {
		if v == value {
			return fmt.Sprint(k)
		}
	}

	panic("KEY NOT FOUND! WTF!?")
}

// LoadConfig ..
func LoadConfig(conf *model.Config) {
	log.Println("Loading config...")

	file := filepath.FromSlash("data/config.json")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		SaveToJSON(conf, file)
		log.Println("Config created..")
		return
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		panic(err)
	}

	log.Println("Config loaded..")
}

// SaveToJSON ..
func SaveToJSON(data interface{}, filename string) {
	file, _ := json.MarshalIndent(data, "", " ")
	if err := ioutil.WriteFile(filepath.FromSlash(filename), file, 0644); err != nil {
		log.Println("Error saving to json file.", err)
		return
	}
}
