package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	conf := "data/testconfig.json"
	data, err := os.ReadFile(conf)
	check(err)

	var config Config
	err = json.Unmarshal(data, &config)
	check(err)

	expected := Config{
		"localhost",
		1234,
		Protocol{"tcp", true},
		60,
		"/tmp",
		true,
	}

	if !reflect.DeepEqual(expected, config) {
		log.Fatalf("got: %+v, expected: %+v", config, expected)
	}

}
