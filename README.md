# dpservice-cli

CLI for [net-dpservice](https://github.com/onmetal/net-dpservice).
<br />

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
<br />

## Autocompletion

To generate autocompletion use:

```shell
dpservice-cli completion [bash|zsh|fish|powershell]
```

Or use -h to get more info and examples for specific shell:

```shell
dpservice-cli completion -h
```
<br />

## Usage

```bash
Usage:
  dpservice-cli [flags]
  dpservice-cli [command]

Available Commands:
  add         Creates one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]
  completion  Generate completion script
  delete      Deletes one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]
  get         Gets one of [interface virtualip loadbalancer lbtarget nat natinfo firewallrule]
  help        Help about any command
  init        Initial set up of the DPDK app
  initialized Indicates if the DPDK app has been initialized already
  list        Lists one of [firewallrules interfaces prefixes lbprefixes routes]

Flags:
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
  -h, --help                       help for dpservice-cli

Use "dpservice-cli [command] --help" for more information about a command.
```
All commands are in [docs](/docs/dpservice-cli.md)

<br />
