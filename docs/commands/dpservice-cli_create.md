## dpservice-cli create

Creates one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]

### Synopsis

Creates one of [interface prefix route virtualip loadbalancer lbprefix lbtarget nat neighbornat firewallrule]

```
dpservice-cli create [flags]
```

### Options

```
  -f, --filename strings   Filename, directory, or URL to file to use to create the resource
  -h, --help               help for create
  -o, --output string      Output format. [json|yaml|table|name] (default "name")
      --pretty             Whether to render pretty output.
  -w, --wide               Whether to render more info in table output.
```

### Options inherited from parent commands

```
      --address string             net-dpservice address. (default "localhost:1337")
      --connect-timeout duration   Timeout to connect to the net-dpservice. (default 4s)
```

### SEE ALSO

* [dpservice-cli](dpservice-cli.md)	 - 
* [dpservice-cli create firewallrule](dpservice-cli_create_firewallrule.md)	 - Create a FirewallRule on interface
* [dpservice-cli create interface](dpservice-cli_create_interface.md)	 - Create an interface
* [dpservice-cli create lbprefix](dpservice-cli_create_lbprefix.md)	 - Create a loadbalancer prefix
* [dpservice-cli create lbtarget](dpservice-cli_create_lbtarget.md)	 - Create a loadbalancer target
* [dpservice-cli create loadbalancer](dpservice-cli_create_loadbalancer.md)	 - Create a loadbalancer
* [dpservice-cli create nat](dpservice-cli_create_nat.md)	 - Create a NAT on interface
* [dpservice-cli create neighbornat](dpservice-cli_create_neighbornat.md)	 - Create a Neighbor NAT
* [dpservice-cli create prefix](dpservice-cli_create_prefix.md)	 - Create a prefix on interface.
* [dpservice-cli create route](dpservice-cli_create_route.md)	 - Create a route
* [dpservice-cli create virtualip](dpservice-cli_create_virtualip.md)	 - Create a virtual IP on interface.

###### Auto generated by spf13/cobra on 26-Jul-2023