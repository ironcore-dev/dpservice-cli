## dpservice-cli get nat

Get NAT on interface

```
dpservice-cli get nat <--interface-id> [flags]
```

### Examples

```
dpservice-cli get nat --interface-id=vm1
```

### Options

```
  -h, --help                  help for nat
      --interface-id string   Interface ID of the NAT.
```

### Options inherited from parent commands

```
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
  -o, --output string              Output format.
      --pretty                     Whether to render pretty output.
```

### SEE ALSO

* [dpservice-cli get](dpservice-cli_get.md)	 - Gets one of [interface virtualip loadbalancer lbtarget nat natinfo firewallrule]

###### Auto generated by spf13/cobra on 17-May-2023