package main

import (
	"embed"
	"encoding/json"
	"path"
)

//go:embed data
var fs embed.FS

var DefaultConfigFile = "config.json"
var DefaultShell = "ash"

type Config struct {
	Addr           string   `json:"address"`
	Port           int      `json:"port"`
	Protocol       Protocol `json:"protocol"`
	Timeout        int64    `json:"timeout"`
	Path           string   `json:"shellpath"`
	Debug          bool     `json:"debug"`
	ExtractApplets bool     `json:"extract_applets"`
}

type Protocol struct {
	ConnType string `json:"conn_type"`
	TLS      bool   `json:"tls"`
}

func readConfig() Config {
	data, err := fs.ReadFile(path.Join("data", DefaultConfigFile))
	check(err)

	var config Config
	err = json.Unmarshal(data, &config)
	check(err)

	return config
}
