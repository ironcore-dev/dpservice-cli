// Copyright 2022 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/onmetal/dpservice-cli/cmd"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/onmetal/net-dpservice-go/errors"
)

var version = "unknown"

func main() {
	util.BuildVersion = version
	if err := cmd.Command().Execute(); err != nil {
		if strings.Contains(err.Error(), "Unimplemented desc") {
			fmt.Println("Error in gRPC, client and server are probably using different proto version")
			os.Exit(errors.SERVER_ERROR)
		}
		// check if it is Server side error
		if err.Error() == strconv.Itoa(errors.SERVER_ERROR) || strings.Contains(err.Error(), "error code") {
			os.Exit(errors.SERVER_ERROR)
		}
		// else it is Client side error
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		os.Exit(errors.CLIENT_ERROR)
	}
}
