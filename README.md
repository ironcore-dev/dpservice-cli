
Grpc client use example:
```golang
package main

import (
	"github.com/onmetal/dpservice-go-library/pkg/client"
)

func main(){
    grpcClient, closer, err := client.New(server)
    if err != nil {
        panic(err)
    }
	defer closer.Close()
	
	res := grpcClient.SomeGrpcCall(...)
```




For build CLI client please run:
```bash
make build
```

```bash
Usage:
  dpservice-go-library [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  machine     
  route       
  vip         

Flags:
  -h, --help            help for dpservice-go-library
      --server string    (default "localhost:1337")

Use "dpservice-go-library [command] --help" for more information about a command.
```