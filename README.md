# dpservice-cli

CLI for [net-dpservice](https://github.com/onmetal/net-dpservice).

## Installation

To build the CLI binary, run

```shell
make build
```

This will build the binary at `bin/dpservice-cli`.

To install it on a system where `GOBIN` is part of the `PATH`,
run

```shell
make install
```

## Usage

```bash
Usage:
  dpservice-cli [flags]
  dpservice-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Creates one of [interface prefix route virtualip]
  delete      Deletes one of [interface prefix route virtualip]
  get         Gets/Lists one of [interface prefix route virtualip]
  help        Help about any command

Flags:
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
  -h, --help                       help for dpservice-cli

Use "dpservice-cli [command] --help" for more information about a command.
```