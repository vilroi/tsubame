package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"time"
)

//go:embed ash
//go:embed config.json
var fs embed.FS

type Config struct {
	Addr     string `json:"address"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Timeout  int64  `json:"timeout"`
}

var (
	ErrorTimeOut     = errors.New("Connection timed out")
	ErrorConnUnavail = errors.New("Connection terminated or unavailable")
)

type NetLineReader struct {
	Conn    net.Conn
	Scanner *bufio.Scanner
	Timeout int64
}

func readConfig() Config {
	data, err := fs.ReadFile("config.json")
	check(err)

	var config Config
	err = json.Unmarshal(data, &config)
	check(err)

	return config
}

func newNetLineReader(conn net.Conn, timeout int64) NetLineReader {
	return NetLineReader{
		conn,
		bufio.NewScanner(conn),
		timeout,
	}
}

func (n *NetLineReader) Readline() ([]byte, error) {
	dur := time.Now().Add(time.Second * time.Duration(n.Timeout))
	n.Conn.SetReadDeadline(dur)

	if !n.Scanner.Scan() {
		if errors.Is(n.Scanner.Err(), os.ErrDeadlineExceeded) {
			return nil, ErrorTimeOut
		}
		return nil, ErrorConnUnavail
	}

	line := append(n.Scanner.Bytes(), byte('\n'))
	return line, nil
}

func loadShell() string {
	data, err := fs.ReadFile("ash")
	check(err)

	path := "/tmp/ash"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	check(err)
	defer f.Close()

	_, err = f.Write(data)
	check(err)

	return f.Name()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
