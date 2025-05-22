# az-dump

A Golang CLI utility export and deploy ARM templates from and to Azure.

## Install

### Go install

Run `go install` to download the binary to the go's binary folder:

```bash
go install github.com/nu12/az-dump@latest
```

Note: go's binary folder (tipically `~/go/bin`) should be added to your PATH.

### From release (x86_64 only)

Download a tagged release binary for your OS (ubuntu, macos, windows) placing it in a folder in your PATH and make it executable (may require elevated permissions):

```bash
wget -O /usr/local/bin/az-dump https://github.com/nu12/az-dump/releases/download/vX.Y.Z/az-dump-linux-amd64.zip
unzip az-dump-linux-amd64.zip
chmod +x az-dump
mv az-dump /usr/local/bin/az-dump
```

Note: replace `X.Y.Z` with a valid version from the repository's releases and `linux-amd64` with the appropriate OS.

### From source

Clone this repo and compile the source code:

```bash
git clone github.com/nu12/az-dump
cd az-dump
go build -o az-dump main.go
```

Move binary to a bin folder in your PATH (may require elevated permissions):
```bash
mv az-dump /usr/local/bin/
```

## Usage

General usage for all commands is `az-dump [command] [flags]`. Find out all available commands with `az-dump`:

```
Export Azure ARM templates to create local backups with backup restore capabilities.

Usage:
  az-dump [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a backup by exporting the ARM templates
  help        Help about any command
  restore     Restore a backup from ARM templates
  version     Show current version

Flags:
  -h, --help   help for az-dump

Use "az-dump [command] --help" for more information about a command.
```

For further information about a specific command, use the `-h` flag (i.e. `az-dump create -h`).

#### Docker

Run the following docker command to run `az-dump` without a local instalation:

```
docker run --rm ghcr.io/nu12/az-dump [command] [flags]
```