# dpservice-cli

Command-line tool for debugging over gRPC for [net-dpservice](https://github.com/onmetal/net-dpservice).

This tool connects directly to a running dp-service and communicates with it (orchestrates it).
<br />

## Installation and developing

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

For more details about developing refer to documentation folder [docs](/docs/development/README.md)
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

Each command or subcommand has help that can be viewed with -h or --help flag.
```shell
dpservice-cli --help
```
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
  -o, --output string              Output format. [json|yaml|table|name]
      --pretty                     Whether to render pretty output.
  -w, --wide                       Whether to render more info in table output.

Use "dpservice-cli [command] --help" for more information about a command.
```
---
Add and Delete commands also support file input with **-f, --filename** flag:
```bash
dpservice-cli [add|delete] -f /<path>/<filename>.[json|yaml]
```
Filename, directory, or URL can be used.
One file can contain multiple objects of any kind, example file:
```bash
{"kind":"VirtualIP","metadata":{"interfaceID":"vm1"},"spec":{"ip":"20.20.20.20"}}
{"kind":"VirtualIP","metadata":{"interfaceID":"vm2"},"spec":{"ip":"20.20.20.21"}}
{"kind":"Prefix","metadata":{"interfaceID":"vm3"},"spec":{"prefix":"20.20.20.0/24"}}
{"kind":"LoadBalancer","metadata":{"id":"4"},"spec":{"vni":100,"lbVipIP":"10.20.30.40","lbports":[{"protocol":6,"port":443},{"protocol":17,"port":53}]}}
{"kind":"LoadBalancerPrefix","metadata":{"interfaceID":"vm1"},"spec":{"prefix":"10.10.10.0/24"}}
```

All commands can be found in [docs](/docs/commands/dpservice-cli.md)

<br />
