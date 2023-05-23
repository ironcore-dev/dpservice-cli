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

	proto "github.com/onmetal/net-dpservice-go/proto"
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

type Status struct {
	Error   int32  `json:"error"`
	Message string `json:"message"`
}

type ServerError struct {
	ServerError Status `json:"serverError"`
}

func (m *ServerError) GetKind() string {
	return "ServerError"
}

func (m *ServerError) GetName() string {
	return fmt.Sprintf("%d", m.ServerError.Error)
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
	Status    Status    `json:"status"`
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
	Status     Status     `json:"status"`
}

type PrefixMeta struct {
	InterfaceID string       `json:"interfaceID"`
	Prefix      netip.Prefix `json:"prefix"`
}

func (m *PrefixMeta) GetName() string {
	return m.Prefix.String()
}

type PrefixSpec struct {
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

type VirtualIP struct {
	TypeMeta      `json:",inline"`
	VirtualIPMeta `json:"metadata"`
	Spec          VirtualIPSpec `json:"spec"`
	Status        Status        `json:"status"`
}

type VirtualIPMeta struct {
	InterfaceID string     `json:"interfaceID"`
	IP          netip.Addr `json:"ip"`
}

func (m *VirtualIPMeta) GetName() string {
	return m.IP.String()
}

type VirtualIPSpec struct {
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

// LoadBalancer section
type LoadBalancer struct {
	TypeMeta         `json:",inline"`
	LoadBalancerMeta `json:"metadata"`
	Spec             LoadBalancerSpec `json:"spec"`
	Status           Status           `json:"status"`
}

type LoadBalancerMeta struct {
	ID string `json:"id"`
}

func (m *LoadBalancerMeta) GetName() string {
	return m.ID
}

type LoadBalancerSpec struct {
	VNI           uint32      `json:"vni,omitempty"`
	LbVipIP       *netip.Addr `json:"lbVipIP,omitempty"`
	Lbports       []LBPort    `json:"lbports,omitempty"`
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

type LBPort struct {
	Protocol uint32 `json:"protocol,omitempty"`
	Port     uint32 `json:"port,omitempty"`
}

type LoadBalancerList struct {
	TypeMeta `json:",inline"`
	Items    []LoadBalancer `json:"items"`
}

func (l *LoadBalancerList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

type LoadBalancerTarget struct {
	TypeMeta               `json:",inline"`
	LoadBalancerTargetMeta `json:"metadata"`
	Spec                   LoadBalancerTargetSpec `json:"spec"`
	Status                 Status                 `json:"status"`
}

type LoadBalancerTargetMeta struct {
	ID string `json:"id"`
}

func (m *LoadBalancerTargetMeta) GetName() string {
	return m.ID
}

type LoadBalancerTargetSpec struct {
	TargetIP *netip.Addr `json:"targetIP,omitempty"`
}

type LoadBalancerTargetList struct {
	TypeMeta `json:",inline"`
	Items    []LoadBalancerTarget `json:"items"`
}

// Interface section
type Interface struct {
	TypeMeta      `json:",inline"`
	InterfaceMeta `json:"metadata"`
	Spec          InterfaceSpec `json:"spec"`
	Status        Status        `json:"status"`
}

type InterfaceMeta struct {
	ID string `json:"id"`
}

func (m *InterfaceMeta) GetName() string {
	return m.ID
}

type InterfaceSpec struct {
	VNI             uint32           `json:"vni,omitempty"`
	Device          string           `json:"device,omitempty"`
	IPs             []netip.Addr     `json:"ips,omitempty"`
	UnderlayRoute   *netip.Addr      `json:"underlayRoute,omitempty"`
	VirtualFunction *VirtualFunction `json:"virtualFunction,omitempty"`
}

type VirtualFunction struct {
	Name     string `json:"vfName,omitempty"`
	Domain   uint32 `json:"vfDomain,omitempty"`
	Bus      uint32 `json:"vfBus,omitempty"`
	Slot     uint32 `json:"vfSlot,omitempty"`
	Function uint32 `json:"vfFunction,omitempty"`
}

func (vf *VirtualFunction) String() string {
	return fmt.Sprintf("Name: %s, Domain: %d, Bus: %d, Slot: %d, Function: %d", vf.Name, vf.Domain, vf.Bus, vf.Slot, vf.Function)
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

// NAT section
type Nat struct {
	TypeMeta `json:",inline"`
	NatMeta  `json:"metadata"`
	Spec     NatSpec `json:"spec"`
	Status   *Status `json:"status,omitempty"`
}

type NatMeta struct {
	InterfaceID string `json:"interfaceID,omitempty"`
}

func (m *NatMeta) GetName() string {
	return m.InterfaceID
}

type NatSpec struct {
	NatVIPIP      *netip.Addr `json:"natVIPIP,omitempty"`
	MinPort       uint32      `json:"minPort,omitempty"`
	MaxPort       uint32      `json:"maxPort,omitempty"`
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

type NatList struct {
	TypeMeta `json:",inline"`
	Items    []Nat `json:"items"`
}

func (l *NatList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

type NeighborNat struct {
	TypeMeta        `json:",inline"`
	NeighborNatMeta `json:"metadata"`
	Spec            NeighborNatSpec `json:"spec"`
	Status          Status          `json:"status"`
}

type NeighborNatMeta struct {
	NatVIPIP *netip.Addr `json:"natVIPIP"`
}

func (m *NeighborNatMeta) GetName() string {
	return m.NatVIPIP.String()
}

type NeighborNatSpec struct {
	Vni           uint32      `json:"vni,omitempty"`
	MinPort       uint32      `json:"minPort,omitempty"`
	MaxPort       uint32      `json:"maxPort,omitempty"`
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

// FirewallRule section
type FirewallRule struct {
	TypeMeta         `json:",inline"`
	FirewallRuleMeta `json:"metadata"`
	Spec             FirewallRuleSpec `json:"spec"`
	Status           Status           `json:"status"`
}

type FirewallRuleMeta struct {
	InterfaceID string `json:"interfaceID"`
	RuleID      string `json:"ruleID"`
}

func (m *FirewallRuleMeta) GetName() string {
	return m.InterfaceID + "/" + m.RuleID
}

type FirewallRuleSpec struct {
	TrafficDirection  string                `json:"trafficDirection,omitempty"`
	FirewallAction    string                `json:"firewallAction,omitempty"`
	Priority          uint32                `json:"priority,omitempty"`
	IpVersion         string                `json:"ipVersion,omitempty"`
	SourcePrefix      *netip.Prefix         `json:"sourcePrefix,omitempty"`
	DestinationPrefix *netip.Prefix         `json:"destinationPrefix,omitempty"`
	ProtocolFilter    *proto.ProtocolFilter `json:"protocolFilter,omitempty"`
}

type FirewallRuleList struct {
	TypeMeta `json:",inline"`
	Items    []FirewallRule `json:"items"`
}

func (l *FirewallRuleList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

var (
	InterfaceKind              = reflect.TypeOf(Interface{}).Name()
	InterfaceListKind          = reflect.TypeOf(InterfaceList{}).Name()
	LoadBalancerKind           = reflect.TypeOf(LoadBalancer{}).Name()
	LoadBalancerTargetKind     = reflect.TypeOf(LoadBalancerTarget{}).Name()
	LoadBalancerTargetListKind = reflect.TypeOf(LoadBalancerTargetList{}).Name()
	PrefixKind                 = reflect.TypeOf(Prefix{}).Name()
	PrefixListKind             = reflect.TypeOf(PrefixList{}).Name()
	VirtualIPKind              = reflect.TypeOf(VirtualIP{}).Name()
	RouteKind                  = reflect.TypeOf(Route{}).Name()
	RouteListKind              = reflect.TypeOf(RouteList{}).Name()
	NatKind                    = reflect.TypeOf(Nat{}).Name()
	NatListKind                = reflect.TypeOf(NatList{}).Name()
	NeighborNatKind            = reflect.TypeOf(NeighborNat{}).Name()
	FirewallRuleKind           = reflect.TypeOf(FirewallRule{}).Name()
	FirewallRuleListKind       = reflect.TypeOf(FirewallRuleList{}).Name()
)
