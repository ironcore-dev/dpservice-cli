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

package api

import (
	"fmt"
	"net/netip"
	"reflect"
)

type Object interface {
	GetKind() string
	GetName() string
}

type List interface {
	GetItems() []Object
}

type TypeMeta struct {
	Kind string `json:"kind"`
}

func (m *TypeMeta) GetKind() string {
	return m.Kind
}

type RouteList struct {
	TypeMeta `json:",inline"`
	Items    []Route `json:"items"`
}

func (l *RouteList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

type Route struct {
	TypeMeta  `json:",inline"`
	RouteMeta `json:"metadata"`
	Spec      RouteSpec `json:"spec"`
}

type RouteMeta struct {
	VNI     uint32       `json:"vni"`
	Prefix  netip.Prefix `json:"prefix"`
	NextHop RouteNextHop `json:"nextHop"`
}

func (m *RouteMeta) GetName() string {
	return fmt.Sprintf("%s-%d:%s", m.Prefix, m.NextHop.VNI, m.NextHop.IP)
}

type RouteSpec struct {
}

type RouteNextHop struct {
	VNI uint32     `json:"vni"`
	IP  netip.Addr `json:"ip"`
}

type PrefixList struct {
	TypeMeta `json:",inline"`
	Items    []Prefix `json:"items"`
}

func (l *PrefixList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

type Prefix struct {
	TypeMeta   `json:",inline"`
	PrefixMeta `json:"metadata"`
	Spec       PrefixSpec `json:"spec"`
}

type PrefixMeta struct {
	InterfaceID string       `json:"interfaceID"`
	Prefix      netip.Prefix `json:"prefix"`
}

func (m *PrefixMeta) GetName() string {
	return m.Prefix.String()
}

type PrefixSpec struct {
}

type VirtualIP struct {
	TypeMeta      `json:",inline"`
	VirtualIPMeta `json:"metadata"`
	Spec          VirtualIPSpec `json:"spec"`
}

type VirtualIPMeta struct {
	InterfaceID string     `json:"interfaceID"`
	IP          netip.Addr `json:"ip"`
}

func (m *VirtualIPMeta) GetName() string {
	return m.IP.String()
}

type VirtualIPSpec struct {
}

type Interface struct {
	TypeMeta      `json:",inline"`
	InterfaceMeta `json:"metadata"`
	Spec          InterfaceSpec   `json:"spec"`
	Status        InterfaceStatus `json:"status"`
}

type InterfaceMeta struct {
	ID string `json:"id"`
}

func (m *InterfaceMeta) GetName() string {
	return m.ID
}

type InterfaceSpec struct {
	VNI    uint32       `json:"vni"`
	Device string       `json:"device"`
	IPs    []netip.Addr `json:"ips"`
}

type InterfaceStatus struct {
	UnderlayIP netip.Addr `json:"underlayIP"`
}

type InterfaceList struct {
	TypeMeta `json:",inline"`
	Items    []Interface `json:"items"`
}

func (l *InterfaceList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

var (
	InterfaceKind     = reflect.TypeOf(Interface{}).Name()
	InterfaceListKind = reflect.TypeOf(InterfaceList{}).Name()
	PrefixKind        = reflect.TypeOf(Prefix{}).Name()
	PrefixListKind    = reflect.TypeOf(PrefixList{}).Name()
	VirtualIPKind     = reflect.TypeOf(VirtualIP{}).Name()
	RouteKind         = reflect.TypeOf(Route{}).Name()
	RouteListKind     = reflect.TypeOf(RouteList{}).Name()
)
