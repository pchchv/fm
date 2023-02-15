[![Go Reference](https://pkg.go.dev/badge/github.com/pchchv/fm.svg)](https://pkg.go.dev/github.com/pchchv/fm)

<div align="center">

# fm

Terminal **f**ile **m**anager

</div>

## Features

* Cross-platform (Linux, macOS, BSDs, Windows)
* Single binary without any runtime dependencies
* Asynchronous IO operations to avoid UI locking
* Extendable and configurable with shell commands
* Customizable keybindings (vi and readline defaults)
* Server/client architecture and remote commands to manage multiple instances
* Fast startup and low memory footprint due to native code and static binaries


Building from the source requires [Go](https://go.dev/).

On Unix (Go version < 1.17):

```bash
env CGO_ENABLED=0 GO111MODULE=on go get -u -ldflags="-s -w" github.com/pchchv/fm
```

On Unix (Go version >= 1.17):

```bash
env CGO_ENABLED=0 go install -ldflags="-s -w" github.com/pchchv/fm
```

On Windows `cmd` (Go version < 1.17):

```cmd
set CGO_ENABLED=0
set GO111MODULE=on
go get -u -ldflags="-s -w" github.com/pchchv/fm
```

On Windows `cmd` (Go version >= 1.17):

```cmd
set CGO_ENABLED=0
go install -ldflags="-s -w" github.com/pchchv/fm
```

On Windows `powershell` (Go version < 1.17):

```powershell
$env:CGO_ENABLED = '0'
$env:GO111MODULE = 'on'
go get -u -ldflags="-s -w" github.com/pchchv/fm
```

On Windows `powershell` (Go version >= 1.17):

```powershell
$env:CGO_ENABLED = '0'
go install -ldflags="-s -w" github.com/pchchv/fm
```

## Usage

After the installation `fm` command should start the application in the current directory.

Run `fm -help` to see command line options.

Run `fm -doc` to see the [documentation](https://pkg.go.dev/github.com/pchchv/fm).
