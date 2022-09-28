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

package flag

import (
	"net/netip"
	"strings"

	"github.com/spf13/pflag"
)

type addrValue netip.Addr

func newAddrValue(val netip.Addr, p *netip.Addr) *addrValue {
	*p = val
	return (*addrValue)(p)
}

func (v *addrValue) String() string {
	return netip.Addr(*v).String()
}

func (v *addrValue) Set(s string) error {
	addr, err := netip.ParseAddr(strings.TrimSpace(s))
	if err != nil {
		return err
	}

	*v = addrValue(addr)
	return nil
}

func (v *addrValue) Type() string {
	return "ip"
}

func AddrVar(f *pflag.FlagSet, p *netip.Addr, name string, value netip.Addr, usage string) {
	f.VarP(newAddrValue(value, p), name, "", usage)
}

func AddrVarP(f *pflag.FlagSet, p *netip.Addr, name, shorthand string, value netip.Addr, usage string) {
	f.VarP(newAddrValue(value, p), name, shorthand, usage)
}
