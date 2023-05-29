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
	GetStatus() int32
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

func (status *Status) String() string {
	return fmt.Sprintf("Error: %d, Message: %s", status.Error, status.Message)
}

type RouteList struct {
	TypeMeta      `json:",inline"`
	RouteListMeta `json:"metadata"`
	Status        Status  `json:"status"`
	Items         []Route `json:"items"`
}

type RouteListMeta struct {
	VNI uint32 `json:"vni"`
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

func (m *Route) GetStatus() int32 {
	return m.Status.Error
}

type RouteSpec struct {
}

type RouteNextHop struct {
	VNI uint32      `json:"vni"`
	IP  *netip.Addr `json:"ip,omitempty"`
}

type PrefixList struct {
	TypeMeta       `json:",inline"`
	PrefixListMeta `json:"metadata"`
	Status         Status   `json:"status"`
	Items          []Prefix `json:"items"`
}

type PrefixListMeta struct {
	InterfaceID string `json:"interfaceID"`
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

func (m *Prefix) GetStatus() int32 {
	return m.Status.Error
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

func (m *VirtualIP) GetStatus() int32 {
	return m.Status.Error
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

func (m *LoadBalancer) GetStatus() int32 {
	return m.Status.Error
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

type LoadBalancerTarget struct {
	TypeMeta               `json:",inline"`
	LoadBalancerTargetMeta `json:"metadata"`
	Spec                   LoadBalancerTargetSpec `json:"spec"`
	Status                 Status                 `json:"status"`
}

type LoadBalancerTargetMeta struct {
	LoadbalancerID string `json:"loadbalancerId"`
}

func (m *LoadBalancerTargetMeta) GetName() string {
	return m.LoadbalancerID
}

func (m *LoadBalancerTarget) GetStatus() int32 {
	return m.Status.Error
}

type LoadBalancerTargetSpec struct {
	TargetIP *netip.Addr `json:"targetIP,omitempty"`
}

type LoadBalancerTargetList struct {
	TypeMeta                   `json:",inline"`
	LoadBalancerTargetListMeta `json:"metadata"`
	Status                     Status               `json:"status"`
	Items                      []LoadBalancerTarget `json:"items"`
}

type LoadBalancerTargetListMeta struct {
	LoadBalancerID string `json:"loadbalancerID"`
}

func (l *LoadBalancerTargetList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

// Interface section
type Interface struct {
	TypeMeta      `json:",inline"`
	InterfaceMeta `json:"metadata"`
	Spec          InterfaceSpec `json:"spec"`
	Status        Status        `json:"status"`
}

type InterfaceMeta struct {
	ID  string `json:"id"`
	PXE PXE    `json:"pxe"`
}

type PXE struct {
	Server   string `json:"server"`
	FileName string `json:"fileName"`
}

func (m *InterfaceMeta) GetName() string {
	return m.ID
}

func (m *Interface) GetStatus() int32 {
	return m.Status.Error
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
	TypeMeta          `json:",inline"`
	InterfaceListMeta `json:"metadata"`
	Status            Status      `json:"status"`
	Items             []Interface `json:"items"`
}

type InterfaceListMeta struct {
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
	Status   Status  `json:"status"`
}

type NatMeta struct {
	InterfaceID string `json:"interfaceID,omitempty"`
}

func (m *NatMeta) GetName() string {
	return m.InterfaceID
}

func (m *Nat) GetStatus() int32 {
	return m.Status.Error
}

type NatSpec struct {
	NatVIPIP      *netip.Addr `json:"natVIPIP,omitempty"`
	MinPort       uint32      `json:"minPort,omitempty"`
	MaxPort       uint32      `json:"maxPort,omitempty"`
	UnderlayRoute *netip.Addr `json:"underlayRoute,omitempty"`
}

type NatList struct {
	TypeMeta    `json:",inline"`
	NatListMeta `json:"metadata"`
	Status      Status `json:"status"`
	Items       []Nat  `json:"items"`
}

type NatListMeta struct {
	NatIP       *netip.Addr `json:"natIp,omitempty"`
	NatInfoType string      `json:"infoType,omitempty"`
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

func (m *NeighborNat) GetStatus() int32 {
	return m.Status.Error
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

func (m *FirewallRule) GetStatus() int32 {
	return m.Status.Error
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
	TypeMeta             `json:",inline"`
	FirewallRuleListMeta `json:"metadata"`
	Status               Status         `json:"status"`
	Items                []FirewallRule `json:"items"`
}

type FirewallRuleListMeta struct {
	InterfaceID string `json:"interfaceID"`
}

func (l *FirewallRuleList) GetItems() []Object {
	res := make([]Object, len(l.Items))
	for i := range l.Items {
		res[i] = &l.Items[i]
	}
	return res
}

type Init struct {
	TypeMeta `json:",inline"`
	InitMeta `json:"metadata"`
	Spec     InitSpec `json:"spec"`
	Status   Status   `json:"status"`
}

type InitMeta struct {
}

type InitSpec struct {
}

func (m *InitMeta) GetName() string {
	return "init"
}

func (m *Init) GetStatus() int32 {
	return m.Status.Error
}

type Initialized struct {
	TypeMeta        `json:",inline"`
	InitializedMeta `json:"metadata"`
	Spec            InitializedSpec `json:"spec"`
	Status          Status          `json:"status"`
}

type InitializedMeta struct {
}

type InitializedSpec struct {
	UUID string `json:"uuid"`
}

func (m *InitializedMeta) GetName() string {
	return "initialized"
}

func (m *Initialized) GetStatus() int32 {
	return 0
}

type Vni struct {
	TypeMeta `json:",inline"`
	VniMeta  `json:"metadata"`
	Spec     VniSpec `json:"spec"`
	Status   Status  `json:"status"`
}

type VniMeta struct {
	VNI     uint32 `json:"vni"`
	VniType uint8  `json:"vniType"`
}

type VniSpec struct {
	InUse bool `json:"inUse"`
}

func (m *VniMeta) GetName() string {
	return "vni"
}

func (m *Vni) GetStatus() int32 {
	return 0
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
	VniKind                    = reflect.TypeOf(Vni{}).Name()
)
