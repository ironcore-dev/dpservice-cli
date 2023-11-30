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

package flag

import (
	"net/netip"
	"strings"

	"github.com/spf13/pflag"
)

type prefixValue netip.Prefix

func newPrefixValue(val netip.Prefix, p *netip.Prefix) *prefixValue {
	*p = val
	return (*prefixValue)(p)
}

func (v *prefixValue) String() string {
	return netip.Prefix(*v).String()
}

func (v *prefixValue) Set(s string) error {
	prefix, err := netip.ParsePrefix(strings.TrimSpace(s))
	if err != nil {
		return err
	}

	*v = prefixValue(prefix)
	return nil
}

func (v *prefixValue) Type() string {
	return "ipprefix"
}

func PrefixVar(f *pflag.FlagSet, p *netip.Prefix, name string, value netip.Prefix, usage string) {
	f.VarP(newPrefixValue(value, p), name, "", usage)
}

func PrefixVarP(f *pflag.FlagSet, p *netip.Prefix, name, shorthand string, value netip.Prefix, usage string) {
	f.VarP(newPrefixValue(value, p), name, shorthand, usage)
}
