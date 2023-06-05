# dpservice-cli development guide

This page is intended as a general overview for all development-oriented topics.
[Cobra](https://github.com/spf13/cobra) framework is used for generating and handling commands in this project.
<br />


## Running dpservice

Before using dpservice-cli client, you need to have dpservice instance running.

Please refer to this guide [net-dpservice](https://github.com/onmetal/net-dpservice/blob/osc/grpc_docs/docs/development/building.md) on how to build dpservice from source.

You can then run python script **/test/dp_service.py** that will start the dpservice with preloaded config.
```bash
sudo ./test/dp_service.py
```
If there is error about number of hugepages run this as root:
```bash
echo 2048 > /sys/devices/system/node/node0/hugepages/hugepages-2048kB/nr_hugepages 
```
<br />


## Running dpservice-cli

Go version 1.18 or newer is needed. \"make\" tool is also needed to utilize the Makefile.  
To run the dpservice-cli client build the binary first and then use it with commands and flags:
```bash
make build
./bin/dpservice-cli -h
```
When you are running dpservice on the same VM you don't need to specify the address and defaults are used (localhost:1337).  
If dpservice is running on different machine or you changed the default settings, use **--address \<string\>** flag:
```bash
./bin/dpservice-cli --address <IP:port> [command] [flags]
```
To change the output format of commands you can use **-o, --output** flag with one of **json | yaml | table | name**

  -  **json**   - shows output in json (you can use **--pretty** flag to show formatted json)
  -  **yaml**   - shows output in yaml
  -  **table**  - shows output in predefined table format (you can use **-w, --wide** for more information)
  -  **name**   - shows only short output with type/name
<br />


## Adding new type

Basic steps when implementing new type (similar to Interface, Route, LoadBalancer, ...):
- Create new type in [/dpdk/api/types.go](/dpdk/api/types.go):
    - create structs and methods
	- at the bottom add new \<type\>Kind variable
- Create new [add|get|list|delete]\<type\>.go file in /cmd/ folder and implement the logic
- Add new command function to subcommands of matching parent command in /cmd/[add|get|list|delete].go
- If needed add aliases for \<type\> at the bottom of [/cmd/common.go](/cmd/common.go)
- Add new function to [/dpdk/api/client.go](/dpdk/api/client.go):
    - add function to Client interface
    - implement the function
- Add new \<type\> to DefaultScheme in [/dpdk/api/register.go](/dpdk/api/register.go)
- If needed create new conversion function(s) between dpdk struct and local struct in [/dpdk/api/conversion.go](/dpdk/api/conversion.go)
- Add new function to show \<type\> as table in [/renderer/renderer.go](/renderer/renderer.go)
    - add new \<type\> to ConvertToTable method
    - implement function to show new \<type\>
<br />


## gRPC

This client uses golang bindings from repo [net-dpservice-go](https://github.com/onmetal/net-dpservice-go).

Definition go files in [proto](https://github.com/onmetal/net-dpservice-go/tree/main/proto) folder are auto-generated from [dpdk.proto](https://github.com/onmetal/net-dpservice/blob/osc/main/proto/dpdk.proto) file in [net-dpservice](https://github.com/onmetal/net-dpservice/) repo.

In case of upgrade of net-dpservice-go, it is needed to pull the latest version:
```
go get github.com/onmetal/net-dpservice-go@[<version>|latest]
```
In order to pull from private repository set the **GOPRIVATE** variable:
```
export GOPRIVATE=github.com/onmetal/*
```

More info about gRPC can be found [here](https://grpc.io/docs/what-is-grpc/introduction/).