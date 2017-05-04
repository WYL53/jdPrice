package main

import (
	"encoding/json"
	"io/ioutil"
)

var conf *Config

type Config struct {
	Port           float64 `json:"port"`
	FrequencyOfDay float64 `json:"frequencyOfDay"`
}

func init() {
	fileBytes, err := ioutil.ReadFile("conf.txt")
	if err != nil {
		panic(err)
	}
	c := new(Config)
	err = json.Unmarshal(fileBytes, c)
	if err != nil {
		panic(err)
	}
	conf = c
}
