// Copyright 2022 IronCore authors
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

package runtime

import "github.com/ironcore-dev/dpservice-go/api"

var DefaultScheme = NewScheme()

func init() {
	if err := DefaultScheme.Add(
		&api.Interface{},
		&api.InterfaceList{},
		&api.Prefix{},
		&api.PrefixList{},
		&api.Route{},
		&api.RouteList{},
		&api.VirtualIP{},
		&api.LoadBalancer{},
		&api.LoadBalancerTarget{},
		&api.LoadBalancerPrefix{},
		&api.LoadBalancerTargetList{},
		&api.Nat{},
		&api.NatList{},
		&api.NeighborNat{},
		&api.FirewallRule{},
		&api.Vni{},
	); err != nil {
		panic(err)
	}
}
