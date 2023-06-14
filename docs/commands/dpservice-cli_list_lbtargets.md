## dpservice-cli list lbtargets

List LoadBalancer Targets

```
dpservice-cli list lbtargets <--lb-id> [flags]
```

### Examples

```
dpservice-cli list lbtargets --lb-id=1
```

### Options

```
  -h, --help           help for lbtargets
      --lb-id string   ID of the loadbalancer to get the targets for.
```

### Options inherited from parent commands

```
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
  -o, --output string              Output format. [json|yaml|table|name] (default "table")
      --pretty                     Whether to render pretty output.
  -w, --wide                       Whether to render more info in table output.
```

### SEE ALSO

* [dpservice-cli list](dpservice-cli_list.md)	 - Lists one of [firewallrules interfaces prefixes lbprefixes routes lbtargets nats]

###### Auto generated by spf13/cobra on 14-Jun-2023