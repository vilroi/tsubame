## Overview

tsubame is a simple reverse shell written with the aim of being a self contained binary which could be deployed on environments where a shell (such as `bash` or `sh`) may not be available.

Upon compilation, tsubame produces a static binary which contains a copy of `busybox`'s `ash` using the [embed package](https://pkg.go.dev/embed). At runtime `tsubame` uses `ash` as a "back-end", passing commands read over the network and feeding it to `ash`. 

For this reason, a copy of `ash` which is compatible with the target architecture must be present in this directory during the build. The default copy of ash has been compiled for x8664, so the resulting binary will not work on an ARM machine, to give an example.

For further details about configuration, please see the relevant section below.

## Usage

## Embedded files
As stated above, `tsubame` embeds some files in the final binary produced. Currently, the following files are included:

- `ash`: A copy of `busybox`'s ash
- `config.json`: The configuration file

If the user would like to replace any of these files, all of the relevant `//go:embed` directives and references to the files must be updated. 

An example of this is shown in the configuration section.

## Configuration

Configuration is done through editing `config.json` in the `./data` directory.

This configuration file is also embedded into the binary along side the shell, and is referenced by `tsubame` at runtime. 

Please do not rename the file, or the program produced will not work.

If you would like to supply a different configuration file or shell, place the files in the `data` directory and update the `Default*` variables in `config.go`.

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

