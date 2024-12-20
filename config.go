package main

import (
	"embed"
	"encoding/json"
	"path"
)

// In order to use a different config file, replace `config.json`
// with the new file name.

//go:embed data
var fs embed.FS

var DefaultConfigFile = "config.json"
var DefaultShell = "ash"

type Config struct {
	Addr     string `json:"address"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Timeout  int64  `json:"timeout"`
	Path     string `json:"shellpath"`
	Debug    bool   `json:"debug"`
}

func readConfig() Config {
	data, err := fs.ReadFile(path.Join("data", DefaultConfigFile))
	check(err)

	var config Config
	err = json.Unmarshal(data, &config)
	check(err)

	return config
}
