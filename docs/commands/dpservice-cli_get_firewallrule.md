## dpservice-cli get firewallrule

Get firewall rule

```
dpservice-cli get firewallrule <--rule-id> <--interface-id> [flags]
```

### Examples

```
dpservice-cli get fwrule --rule-id=1 --interface-id=vm1
```

### Options

```
  -h, --help                  help for firewallrule
      --interface-id string   Interface ID where is firewall rule.
      --rule-id string        Rule ID to get.
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