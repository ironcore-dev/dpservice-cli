## dpservice-cli



```
dpservice-cli [flags]
```

### Options

```
      --address string             dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the dpservice. (default 4s)
  -h, --help                       help for dpservice-cli
  -o, --output string              Output format. [json|yaml|table|name]
      --pretty                     Whether to render pretty output.
  -w, --wide                       Whether to render more info in table output.
```

### SEE ALSO

* [dpservice-cli completion](dpservice-cli_completion.md)	 - Generate completion script
* [dpservice-cli create](dpservice-cli_create.md)	 - Creates one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]
* [dpservice-cli delete](dpservice-cli_delete.md)	 - Deletes one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]
* [dpservice-cli get](dpservice-cli_get.md)	 - Gets one of [interface virtualip loadbalancer nat firewallrule vni version init]
* [dpservice-cli init](dpservice-cli_init.md)	 - Initial set up of the DPDK app
* [dpservice-cli list](dpservice-cli_list.md)	 - Lists one of [firewallrules interfaces prefixes lbprefixes routes lbtargets nats]
* [dpservice-cli reset](dpservice-cli_reset.md)	 - Resets one of [vni]

###### Auto generated by spf13/cobra on 25-Sep-2023
