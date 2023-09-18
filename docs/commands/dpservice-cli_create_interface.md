## dpservice-cli create interface

Create an interface

```
dpservice-cli create interface <--id> [<--ip>] <--vni> <--device> [flags]
```

### Examples

```
dpservice-cli create interface --id=vm4 --ip=10.200.1.4 --ip=2000:200:1::4 --vni=200 --device=net_tap5
```

### Options

```
      --device string          Device to allocate.
  -h, --help                   help for interface
      --id string              ID of the interface.
      --ip addrSlice           IP to assign to the interface. (default [])
      --pxe-file-name string   PXE boot file name.
      --pxe-server string      PXE next server.
      --vni uint32             VNI to add the interface to.
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