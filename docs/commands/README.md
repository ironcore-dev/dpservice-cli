# dpservice-cli commands:

You can browse help for all commands starting in main command [here](/docs/commands/dpservice-cli.md)

# Available commands:

Most of the validation is done on server side (dpservice).
All parameters are validated based on their type (see below).
In some cases there is validation also on client side (dpservice-cli) user is then notified with proper usage.

## Initialization/check for initialized service and generating auto-completion:
```
init
initialized
completion [bash|zsh|fish|powershell]
```

## Add/delete/list network interfaces:
```
add interface --id=<string> --ip=<netip.Addr> --ip=<netip.Addr> --vni=<uint32> --device=<string>
delete interface --id=<string>
get interface --id=<string>
list interfaces
```

## Add/delete/list routes (ip route equivalents):
```
add route --prefix=<netip.Prefix> --next-hop-vni=<uint32> --next-hop-ip=<netip.Addr> --vni=<uint32>
delete route --prefix=<netip.Prefix> --vni=<uint32>
list routes --vni=<uint32>
```

## Add/delete/list prefixes (to route other IP ranges to a given interface):
```
add prefix --prefix=<netip.Prefix> --interface-id=<string>
delete prefix --prefix=<netip.Prefix> --interface-id=<string>
list prefixes --interface-id=<string>
```

## Create/delete/list loadbalancers:
```
add loadbalancer --id=<string> --vni=<uint32> --vip=<netip.Addr> --lbports=<string>
delete loadbalancer --id=<string>
get loadbalancer --id=<string>
```

## Add/delete/list loadbalancer backing IPs:
```
add lbtarget --target-ip=<netip.Addr> --lb-id=<string>
delete lbtarget --target-ip=<netip.Addr> --lb-id=<string>
list lbtargets --lb-id=<string>
```

## Add/delete/list loadbalancer prefixes (call on loadbalancer targets so the public IP packets can reach them):
```
add lbprefix --prefix=<netip.Prefix> --interface-id=<string>
delete lbprefix --prefix=<netip.Prefix> --interface-id=<string>
list lbprefixes --interface-id=<string>
```

## Add/delete/list a virtual IP for the interface (SNAT):
```
add virtualip --vip=<netip.Addr> --interface-id=<string>
delete virtualip --interface-id=<string>
get virtualip --interface-id=<string>
```

## Add/delete/list NAT IP (with port range) for the interface:
```
add nat --interface-id=<string> --natip=<netip.Addr> --minport=<uint32> --maxport=<uint32>
delete nat --interface-id=<string>
get nat --interface-id=<string>
list nats
```

## Add/delete/list neighbors (dp-services) with the same NAT IP:
```
add neighbornat --natip=<netip.Addr> --vni=<uint32> --minport=<uint32> --maxport=<uint32> --underlayroute=<netip.Addr>
delete neighbornat --natip=<netip.Addr> --vni=<uint32> --minport=<uint32> --maxport=<uint32>
get natinfo --nat-ip=<netip.Addr> --info-type=<string>
```

## Add/delete/list firewall rules:
```
add fwrule --interface-id=<string> --action=<string> --direction=<string> --dst=<netip.Prefix> --ipver=<string> --priority=<uint32> --rule-id=<string> --src=<netip.Prefix> --protocol=<string> --src-port-min=<int32> --src-port-max=<int32> --dst-port-min=<int32> --dst-port-max=<int32> --icmp-type=<int32> --icmp-code=<int32>
delete firewallrule --rule-id=<string> --interface-id=<string>
get fwrule --rule-id=<string> --interface-id=<string>
list firewallrules --interface-id=<string>
```