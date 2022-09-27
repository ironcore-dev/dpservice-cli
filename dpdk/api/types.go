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
	"net/netip"
	"reflect"
)

type TypeMeta struct {
	Kind string
}

type RouteList struct {
	TypeMeta `json:",inline"`
	Items    []Route `json:"items"`
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

type Prefix struct {
	TypeMeta   `json:",inline"`
	PrefixMeta `json:"metadata"`
	Spec       PrefixSpec `json:"spec"`
}

type PrefixMeta struct {
	InterfaceID string       `json:"interfaceID"`
	Prefix      netip.Prefix `json:"prefix"`
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

var (
	InterfaceKind     = reflect.TypeOf(Interface{}).Name()
	InterfaceListKind = reflect.TypeOf(InterfaceList{}).Name()
	PrefixKind        = reflect.TypeOf(Prefix{}).Name()
	PrefixListKind    = reflect.TypeOf(PrefixList{}).Name()
	VirtualIPKind     = reflect.TypeOf(VirtualIP{}).Name()
	RouteKind         = reflect.TypeOf(Route{}).Name()
	RouteListKind     = reflect.TypeOf(RouteList{}).Name()
)
