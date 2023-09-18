## dpservice-cli create lbprefix

Create a loadbalancer prefix

```
dpservice-cli create lbprefix <--prefix> <--interface-id> [flags]
```

### Examples

```
dpservice-cli create lbprefix --prefix=10.10.10.0/24 --interface-id=vm1
```

### Options

```
  -h, --help                  help for lbprefix
      --interface-id string   ID of the interface to create the prefix for.
      --prefix ipprefix       Prefix to add to the interface. (default invalid Prefix)
```

### Options inherited from parent commands

```
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
  -o, --output string              Output format. [json|yaml|table|name] (default "name")
      --pretty                     Whether to render pretty output.
  -w, --wide                       Whether to render more info in table output.
```

### SEE ALSO

* [dpservice-cli create](dpservice-cli_create.md)	 - Creates one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]

###### Auto generated by spf13/cobra on 26-Jul-2023