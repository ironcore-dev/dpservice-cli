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

package client

import (
	"context"
	"fmt"
	"net/netip"
	"strings"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	apierrors "github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/netiputil"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
)

type Client interface {
	GetLoadBalancer(ctx context.Context, id string) (*api.LoadBalancer, error)
	CreateLoadBalancer(ctx context.Context, lb *api.LoadBalancer) (*api.LoadBalancer, error)
	DeleteLoadBalancer(ctx context.Context, id string) error

	ListLoadBalancerPrefixes(ctx context.Context, interfaceID string) (*api.PrefixList, error)
	CreateLoadBalancerPrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error)
	DeleteLoadBalancerPrefix(ctx context.Context, interfaceID string, prefix netip.Prefix) error

	GetLoadBalancerTargets(ctx context.Context, interfaceID string) (*api.LoadBalancerTargetList, error)
	AddLoadBalancerTarget(ctx context.Context, lbtarget *api.LoadBalancerTarget) (*api.LoadBalancerTarget, error)
	DeleteLoadBalancerTarget(ctx context.Context, id string, targetIP netip.Addr) error

	GetInterface(ctx context.Context, id string) (*api.Interface, error)
	ListInterfaces(ctx context.Context) (*api.InterfaceList, error)
	CreateInterface(ctx context.Context, iface *api.Interface) (*api.Interface, error)
	DeleteInterface(ctx context.Context, id string) error

	GetVirtualIP(ctx context.Context, interfaceID string) (*api.VirtualIP, error)
	AddVirtualIP(ctx context.Context, virtualIP *api.VirtualIP) (*api.VirtualIP, error)
	DeleteVirtualIP(ctx context.Context, interfaceID string) error

	ListPrefixes(ctx context.Context, interfaceID string) (*api.PrefixList, error)
	AddPrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error)
	DeletePrefix(ctx context.Context, interfaceID string, prefix netip.Prefix) error

	ListRoutes(ctx context.Context, vni uint32) (*api.RouteList, error)
	AddRoute(ctx context.Context, route *api.Route) (*api.Route, error)
	DeleteRoute(ctx context.Context, vni uint32, prefix netip.Prefix, nextHopVNI uint32, nextHopIP netip.Addr) error

	GetNat(ctx context.Context, interfaceID string) (*api.Nat, error)
	AddNat(ctx context.Context, nat *api.Nat) (*api.Nat, error)
	DeleteNat(ctx context.Context, interfaceID string) error

	AddNeighborNat(ctx context.Context, nat *api.NeighborNat) (*api.NeighborNat, error)
	GetNATInfo(ctx context.Context, natVIPIP netip.Addr, natType string) (*api.NatList, error)
	DeleteNeighborNat(ctx context.Context, neigbhorNat api.NeighborNat) error

	ListFirewallRules(ctx context.Context, interfaceID string) (*api.FirewallRuleList, error)
	AddFirewallRule(ctx context.Context, fwRule *api.FirewallRule) (*api.FirewallRule, error)
	GetFirewallRule(ctx context.Context, interfaceID string, ruleID string) (*api.FirewallRule, error)
	DeleteFirewallRule(ctx context.Context, interfaceID string, ruleID string) (*api.FirewallRule, error)

	Initialized(ctx context.Context) (string, error)
	Init(ctx context.Context, initConfig dpdkproto.InitConfig) error
}

type client struct {
	dpdkproto.DPDKonmetalClient
}

func NewClient(protoClient dpdkproto.DPDKonmetalClient) Client {
	return &client{protoClient}
}

func (c *client) GetLoadBalancer(ctx context.Context, id string) (*api.LoadBalancer, error) {
	res, err := c.DPDKonmetalClient.GetLoadBalancer(ctx, &dpdkproto.GetLoadBalancerRequest{LoadBalancerID: []byte(id)})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	lb, err := api.ProtoLoadBalancerToLoadBalancer(res, id)
	return lb, err
}

func (c *client) CreateLoadBalancer(ctx context.Context, lb *api.LoadBalancer) (*api.LoadBalancer, error) {
	var lbPorts = make([]*dpdkproto.LBPort, 0, len(lb.Spec.Lbports))
	for _, p := range lb.Spec.Lbports {
		lbPort := &dpdkproto.LBPort{Port: p.Port, Protocol: dpdkproto.Protocol(p.Protocol)}
		lbPorts = append(lbPorts, lbPort)
	}
	res, err := c.DPDKonmetalClient.CreateLoadBalancer(ctx, &dpdkproto.CreateLoadBalancerRequest{
		LoadBalancerID: []byte(lb.LoadBalancerMeta.ID),
		Vni:            lb.Spec.VNI,
		LbVipIP:        api.LbipToProtoLbip(*lb.Spec.LbVipIP),
		Lbports:        lbPorts,
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}

	underlayRoute, err := netip.ParseAddr(string(res.GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}
	lb.Spec.UnderlayRoute = &underlayRoute

	return &api.LoadBalancer{
		TypeMeta:         api.TypeMeta{Kind: api.LoadBalancerKind},
		LoadBalancerMeta: lb.LoadBalancerMeta,
		Spec:             lb.Spec,
		Status:           api.ProtoStatusToStatus(res.Status),
	}, nil
}

func (c *client) DeleteLoadBalancer(ctx context.Context, id string) error {
	res, err := c.DPDKonmetalClient.DeleteLoadBalancer(ctx, &dpdkproto.DeleteLoadBalancerRequest{LoadBalancerID: []byte(id)})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) ListLoadBalancerPrefixes(ctx context.Context, interfaceID string) (*api.PrefixList, error) {
	res, err := c.DPDKonmetalClient.ListInterfaceLoadBalancerPrefixes(ctx, &dpdkproto.ListInterfaceLoadBalancerPrefixesRequest{
		InterfaceID: []byte(interfaceID),
	})
	if err != nil {
		return nil, err
	}

	prefixes := make([]api.Prefix, len(res.GetPrefixes()))
	for i, dpdkPrefix := range res.GetPrefixes() {
		prefix, err := api.ProtoPrefixToPrefix(interfaceID, api.ProtoLBPrefixToProtoPrefix(*dpdkPrefix))
		if err != nil {
			return nil, err
		}

		prefixes[i] = *prefix
	}

	return &api.PrefixList{
		TypeMeta: api.TypeMeta{Kind: "LoadBalancerPrefixList"},
		Items:    prefixes,
	}, nil
}

func (c *client) CreateLoadBalancerPrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error) {
	res, err := c.DPDKonmetalClient.CreateInterfaceLoadBalancerPrefix(ctx, &dpdkproto.CreateInterfaceLoadBalancerPrefixRequest{
		InterfaceID: &dpdkproto.InterfaceIDMsg{
			InterfaceID: []byte(prefix.InterfaceID),
		},
		Prefix: &dpdkproto.Prefix{
			IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Prefix.Addr()),
			Address:      []byte(prefix.Prefix.Addr().String()),
			PrefixLength: uint32(prefix.Prefix.Bits()),
		},
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	underlayRoute, err := netip.ParseAddr(string(res.GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}
	return &api.Prefix{
		TypeMeta:   api.TypeMeta{Kind: "LoadBalancerPrefix"},
		PrefixMeta: prefix.PrefixMeta,
		Spec:       api.PrefixSpec{UnderlayRoute: &underlayRoute},
		Status:     api.ProtoStatusToStatus(res.Status),
	}, nil
}

func (c *client) DeleteLoadBalancerPrefix(ctx context.Context, interfaceID string, prefix netip.Prefix) error {
	res, err := c.DPDKonmetalClient.DeleteInterfaceLoadBalancerPrefix(ctx, &dpdkproto.DeleteInterfaceLoadBalancerPrefixRequest{
		InterfaceID: &dpdkproto.InterfaceIDMsg{
			InterfaceID: []byte(interfaceID),
		},
		Prefix: &dpdkproto.Prefix{
			IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Addr()),
			Address:      []byte(prefix.Addr().String()),
			PrefixLength: uint32(prefix.Bits()),
		},
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) GetLoadBalancerTargets(ctx context.Context, loadBalancerID string) (*api.LoadBalancerTargetList, error) {
	res, err := c.DPDKonmetalClient.GetLoadBalancerTargets(ctx, &dpdkproto.GetLoadBalancerTargetsRequest{
		LoadBalancerID: []byte(loadBalancerID),
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}

	lbtargets := make([]api.LoadBalancerTarget, len(res.GetTargetIPs()))
	for i, dpdkLBtarget := range res.GetTargetIPs() {
		var lbtarget api.LoadBalancerTarget
		lbtarget.TypeMeta.Kind = api.LoadBalancerTargetKind
		lbtarget.Spec.TargetIP = api.ProtoLbipToLbip(*dpdkLBtarget)
		lbtarget.LoadBalancerTargetMeta.ID = loadBalancerID

		lbtargets[i] = lbtarget
	}

	return &api.LoadBalancerTargetList{
		TypeMeta: api.TypeMeta{Kind: api.LoadBalancerTargetListKind},
		Items:    lbtargets,
	}, nil
}

func (c *client) AddLoadBalancerTarget(ctx context.Context, lbtarget *api.LoadBalancerTarget) (*api.LoadBalancerTarget, error) {
	res, err := c.DPDKonmetalClient.AddLoadBalancerTarget(ctx, &dpdkproto.AddLoadBalancerTargetRequest{
		LoadBalancerID: []byte(lbtarget.LoadBalancerTargetMeta.ID),
		TargetIP:       api.LbipToProtoLbip(*lbtarget.Spec.TargetIP),
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetMessage())
	}

	return &api.LoadBalancerTarget{
		TypeMeta:               api.TypeMeta{Kind: api.LoadBalancerTargetKind},
		LoadBalancerTargetMeta: lbtarget.LoadBalancerTargetMeta,
		Spec:                   lbtarget.Spec,
		Status:                 api.ProtoStatusToStatus(res),
	}, nil
}

func (c *client) DeleteLoadBalancerTarget(ctx context.Context, id string, targetIP netip.Addr) error {
	res, err := c.DPDKonmetalClient.DeleteLoadBalancerTarget(ctx, &dpdkproto.DeleteLoadBalancerTargetRequest{
		LoadBalancerID: []byte(id),
		TargetIP:       api.LbipToProtoLbip(targetIP),
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) GetInterface(ctx context.Context, name string) (*api.Interface, error) {
	res, err := c.DPDKonmetalClient.GetInterface(ctx, &dpdkproto.InterfaceIDMsg{InterfaceID: []byte(name)})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	return api.ProtoInterfaceToInterface(res.GetInterface())
}

func (c *client) ListInterfaces(ctx context.Context) (*api.InterfaceList, error) {
	res, err := c.DPDKonmetalClient.ListInterfaces(ctx, &dpdkproto.Empty{})
	if err != nil {
		return nil, err
	}

	ifaces := make([]api.Interface, len(res.GetInterfaces()))
	for i, dpdkIface := range res.GetInterfaces() {
		iface, err := api.ProtoInterfaceToInterface(dpdkIface)
		if err != nil {
			return nil, err
		}

		ifaces[i] = *iface
	}

	return &api.InterfaceList{
		TypeMeta: api.TypeMeta{Kind: api.InterfaceListKind},
		Items:    ifaces,
	}, nil
}

func (c *client) CreateInterface(ctx context.Context, iface *api.Interface) (*api.Interface, error) {
	res, err := c.DPDKonmetalClient.CreateInterface(ctx, &dpdkproto.CreateInterfaceRequest{
		InterfaceType: dpdkproto.InterfaceType_VirtualInterface,
		InterfaceID:   []byte(iface.ID),
		Vni:           iface.Spec.VNI,
		Ipv4Config:    api.NetIPAddrToProtoIPConfig(netiputil.FindIPv4(iface.Spec.IPs)),
		Ipv6Config:    api.NetIPAddrToProtoIPConfig(netiputil.FindIPv6(iface.Spec.IPs)),
		DeviceName:    iface.Spec.Device,
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetResponse().GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetResponse().GetStatus().GetMessage())
	}

	underlayRoute, err := netip.ParseAddr(string(res.GetResponse().GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}

	return &api.Interface{
		TypeMeta:      api.TypeMeta{Kind: api.InterfaceKind},
		InterfaceMeta: iface.InterfaceMeta,
		Spec: api.InterfaceSpec{
			VNI:           iface.Spec.VNI,
			Device:        iface.Spec.Device,
			IPs:           iface.Spec.IPs,
			UnderlayRoute: &underlayRoute,
			VirtualFunction: &api.VirtualFunction{
				Name:     res.Vf.Name,
				Domain:   res.Vf.Domain,
				Bus:      res.Vf.Bus,
				Slot:     res.Vf.Slot,
				Function: res.Vf.Function,
			},
		},
		Status: api.ProtoStatusToStatus(res.Response.Status),
	}, nil
}

func (c *client) DeleteInterface(ctx context.Context, name string) error {
	res, err := c.DPDKonmetalClient.DeleteInterface(ctx, &dpdkproto.InterfaceIDMsg{InterfaceID: []byte(name)})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) GetVirtualIP(ctx context.Context, interfaceName string) (*api.VirtualIP, error) {
	res, err := c.DPDKonmetalClient.GetInterfaceVIP(ctx, &dpdkproto.InterfaceIDMsg{
		InterfaceID: []byte(interfaceName),
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}

	return api.ProtoVirtualIPToVirtualIP(interfaceName, res)
}

func (c *client) AddVirtualIP(ctx context.Context, virtualIP *api.VirtualIP) (*api.VirtualIP, error) {
	res, err := c.DPDKonmetalClient.AddInterfaceVIP(ctx, &dpdkproto.InterfaceVIPMsg{
		InterfaceID: []byte(virtualIP.InterfaceID),
		InterfaceVIPIP: &dpdkproto.InterfaceVIPIP{
			IpVersion: api.NetIPAddrToProtoIPVersion(virtualIP.IP),
			Address:   []byte(virtualIP.IP.String()),
		},
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	underlayRoute, err := netip.ParseAddr(string(res.GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}

	return &api.VirtualIP{
		TypeMeta:      api.TypeMeta{Kind: api.VirtualIPKind},
		VirtualIPMeta: virtualIP.VirtualIPMeta,
		Spec:          api.VirtualIPSpec{UnderlayRoute: &underlayRoute},
		Status:        api.ProtoStatusToStatus(res.Status),
	}, nil
}

func (c *client) DeleteVirtualIP(ctx context.Context, interfaceID string) error {
	res, err := c.DPDKonmetalClient.DeleteInterfaceVIP(ctx, &dpdkproto.InterfaceIDMsg{
		InterfaceID: []byte(interfaceID),
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) ListPrefixes(ctx context.Context, interfaceID string) (*api.PrefixList, error) {
	res, err := c.DPDKonmetalClient.ListInterfacePrefixes(ctx, &dpdkproto.InterfaceIDMsg{
		InterfaceID: []byte(interfaceID),
	})
	if err != nil {
		return nil, err
	}

	prefixes := make([]api.Prefix, len(res.GetPrefixes()))
	for i, dpdkPrefix := range res.GetPrefixes() {
		prefix, err := api.ProtoPrefixToPrefix(interfaceID, dpdkPrefix)
		if err != nil {
			return nil, err
		}

		prefixes[i] = *prefix
	}

	return &api.PrefixList{
		TypeMeta: api.TypeMeta{Kind: api.PrefixListKind},
		Items:    prefixes,
	}, nil
}

func (c *client) AddPrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error) {
	res, err := c.DPDKonmetalClient.AddInterfacePrefix(ctx, &dpdkproto.InterfacePrefixMsg{
		InterfaceID: &dpdkproto.InterfaceIDMsg{
			InterfaceID: []byte(prefix.InterfaceID),
		},
		Prefix: &dpdkproto.Prefix{
			IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Prefix.Addr()),
			Address:      []byte(prefix.Prefix.Addr().String()),
			PrefixLength: uint32(prefix.Prefix.Bits()),
		},
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	underlayRoute, err := netip.ParseAddr(string(res.GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}
	return &api.Prefix{
		TypeMeta:   api.TypeMeta{Kind: api.PrefixKind},
		PrefixMeta: prefix.PrefixMeta,
		Spec:       api.PrefixSpec{UnderlayRoute: &underlayRoute},
		Status:     api.ProtoStatusToStatus(res.Status),
	}, nil
}

func (c *client) DeletePrefix(ctx context.Context, interfaceID string, prefix netip.Prefix) error {
	res, err := c.DPDKonmetalClient.DeleteInterfacePrefix(ctx, &dpdkproto.InterfacePrefixMsg{
		InterfaceID: &dpdkproto.InterfaceIDMsg{
			InterfaceID: []byte(interfaceID),
		},
		Prefix: &dpdkproto.Prefix{
			IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Addr()),
			Address:      []byte(prefix.Addr().String()),
			PrefixLength: uint32(prefix.Bits()),
		},
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) AddRoute(ctx context.Context, route *api.Route) (*api.Route, error) {
	res, err := c.DPDKonmetalClient.AddRoute(ctx, &dpdkproto.VNIRouteMsg{
		Vni: &dpdkproto.VNIMsg{Vni: route.VNI},
		Route: &dpdkproto.Route{
			IpVersion: api.NetIPAddrToProtoIPVersion(route.NextHop.IP),
			Weight:    100,
			Prefix: &dpdkproto.Prefix{
				IpVersion:    api.NetIPAddrToProtoIPVersion(route.Prefix.Addr()),
				Address:      []byte(route.Prefix.Addr().String()),
				PrefixLength: uint32(route.Prefix.Bits()),
			},
			NexthopVNI:     route.NextHop.VNI,
			NexthopAddress: []byte(route.NextHop.IP.String()),
		},
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return &api.Route{
		TypeMeta:  api.TypeMeta{Kind: api.RouteKind},
		RouteMeta: route.RouteMeta,
		Spec:      route.Spec,
		Status:    api.ProtoStatusToStatus(res),
	}, nil
}

func (c *client) DeleteRoute(ctx context.Context, vni uint32, prefix netip.Prefix, nextHopVNI uint32, nextHopIP netip.Addr) error {
	res, err := c.DPDKonmetalClient.DeleteRoute(ctx, &dpdkproto.VNIRouteMsg{
		Vni: &dpdkproto.VNIMsg{Vni: vni},
		Route: &dpdkproto.Route{
			IpVersion: api.NetIPAddrToProtoIPVersion(nextHopIP),
			Weight:    100,
			Prefix: &dpdkproto.Prefix{
				IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Addr()),
				Address:      []byte(prefix.Addr().String()),
				PrefixLength: uint32(prefix.Bits()),
			},
			NexthopVNI:     nextHopVNI,
			NexthopAddress: []byte(nextHopIP.String()),
		},
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) ListRoutes(ctx context.Context, vni uint32) (*api.RouteList, error) {
	res, err := c.DPDKonmetalClient.ListRoutes(ctx, &dpdkproto.VNIMsg{
		Vni: vni,
	})
	if err != nil {
		return nil, err
	}

	routes := make([]api.Route, len(res.GetRoutes()))
	for i, dpdkRoute := range res.GetRoutes() {
		route, err := api.ProtoRouteToRoute(vni, dpdkRoute)
		if err != nil {
			return nil, err
		}

		routes[i] = *route
	}

	return &api.RouteList{
		TypeMeta: api.TypeMeta{Kind: api.RouteListKind},
		Items:    routes,
	}, nil
}

func (c *client) GetNat(ctx context.Context, interfaceID string) (*api.Nat, error) {
	res, err := c.DPDKonmetalClient.GetNAT(ctx, &dpdkproto.GetNATRequest{InterfaceID: []byte(interfaceID)})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}
	nat, err := api.ProtoNatToNat(res, interfaceID)
	return nat, err
}

func (c *client) AddNat(ctx context.Context, nat *api.Nat) (*api.Nat, error) {
	res, err := c.DPDKonmetalClient.AddNAT(ctx, &dpdkproto.AddNATRequest{
		InterfaceID: []byte(nat.NatMeta.InterfaceID),
		NatVIPIP: &dpdkproto.NATIP{
			IpVersion: api.NetIPAddrToProtoIPVersion(*nat.Spec.NatVIPIP),
			Address:   []byte(nat.Spec.NatVIPIP.String()),
		},
		MinPort: nat.Spec.MinPort,
		MaxPort: nat.Spec.MaxPort,
	})
	if err != nil {
		return nil, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return nil, apierrors.NewStatusError(errorCode, res.GetStatus().GetMessage())
	}

	underlayRoute, err := netip.ParseAddr(string(res.GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}
	nat.Spec.UnderlayRoute = &underlayRoute
	status := api.ProtoStatusToStatus(res.Status)

	return &api.Nat{
		TypeMeta: api.TypeMeta{Kind: api.NatKind},
		NatMeta:  nat.NatMeta,
		Spec:     nat.Spec,
		Status:   &status,
	}, nil
}

func (c *client) DeleteNat(ctx context.Context, interfaceID string) error {
	res, err := c.DPDKonmetalClient.DeleteNAT(ctx, &dpdkproto.DeleteNATRequest{
		InterfaceID: []byte(interfaceID),
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) AddNeighborNat(ctx context.Context, nNat *api.NeighborNat) (*api.NeighborNat, error) {

	res, err := c.DPDKonmetalClient.AddNeighborNAT(ctx, &dpdkproto.AddNeighborNATRequest{
		NatVIPIP: &dpdkproto.NATIP{
			IpVersion: api.NetIPAddrToProtoIPVersion(*nNat.NeighborNatMeta.NatVIPIP),
			Address:   []byte(nNat.NeighborNatMeta.NatVIPIP.String()),
		},
		Vni:           nNat.Spec.Vni,
		MinPort:       nNat.Spec.MinPort,
		MaxPort:       nNat.Spec.MaxPort,
		UnderlayRoute: []byte(nNat.Spec.UnderlayRoute.String()),
	})
	if err != nil {
		return nil, err
	}

	if res.Error == 0 {
		return &api.NeighborNat{
			TypeMeta:        api.TypeMeta{Kind: api.NeighborNatKind},
			NeighborNatMeta: nNat.NeighborNatMeta,
			Spec:            nNat.Spec,
			Status:          api.ProtoStatusToStatus(res),
		}, nil
	}
	return nil, fmt.Errorf("%d", res.Error)
}

func (c *client) GetNATInfo(ctx context.Context, natVIPIP netip.Addr, natType string) (*api.NatList, error) {
	var nType int32
	switch strings.ToLower(natType) {
	case "local", "1":
		nType = 1
	case "neigh", "2", "neighbor":
		nType = 2
	default:
		return nil, fmt.Errorf("nat info type can be only: Local = 1/Neigh(bor) = 2")
	}

	res, err := c.DPDKonmetalClient.GetNATInfo(ctx, &dpdkproto.GetNATInfoRequest{
		NatVIPIP: &dpdkproto.NATIP{IpVersion: api.NetIPAddrToProtoIPVersion(natVIPIP),
			Address: []byte(natVIPIP.String()),
		},
		NatInfoType: dpdkproto.NATInfoType(nType),
	})
	if err != nil {
		return nil, err
	}

	var nats = make([]api.Nat, len(res.NatInfoEntries))
	var nat api.Nat
	for i, natInfoEntry := range res.GetNatInfoEntries() {

		var underlayRoute, vipIP netip.Addr
		if res.NatInfoType == 2 {
			underlayRoute, err = netip.ParseAddr(string(natInfoEntry.GetUnderlayRoute()))
			if err != nil {
				return nil, fmt.Errorf("error parsing underlay route: %w", err)
			}
			nat.Spec.UnderlayRoute = &underlayRoute
			vipIP, err = netip.ParseAddr(string(res.NatVIPIP.Address))
			if err != nil {
				return nil, fmt.Errorf("error parsing vip ip: %w", err)
			}
			nat.Spec.NatVIPIP = &vipIP
		} else if res.NatInfoType == 1 {
			vipIP, err = netip.ParseAddr(string(natInfoEntry.GetAddress()))
			if err != nil {
				return nil, fmt.Errorf("error parsing vip ip: %w", err)
			}
			nat.Spec.NatVIPIP = &vipIP
		}
		nat.Kind = api.NatKind
		nat.Spec.MinPort = natInfoEntry.MinPort
		nat.Spec.MaxPort = natInfoEntry.MaxPort
		nats[i] = nat
	}
	return &api.NatList{
		TypeMeta: api.TypeMeta{Kind: api.NatListKind},
		Items:    nats,
	}, nil
}

func (c *client) DeleteNeighborNat(ctx context.Context, neigbhorNat api.NeighborNat) error {
	res, err := c.DPDKonmetalClient.DeleteNeighborNAT(ctx, &dpdkproto.DeleteNeighborNATRequest{
		NatVIPIP: &dpdkproto.NATIP{
			IpVersion: api.NetIPAddrToProtoIPVersion(*neigbhorNat.NatVIPIP),
			Address:   []byte(neigbhorNat.NatVIPIP.String()),
		},
		Vni:     neigbhorNat.Spec.Vni,
		MinPort: neigbhorNat.Spec.MinPort,
		MaxPort: neigbhorNat.Spec.MaxPort,
	})
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}

func (c *client) ListFirewallRules(ctx context.Context, interfaceID string) (*api.FirewallRuleList, error) {
	res, err := c.DPDKonmetalClient.ListFirewallRules(ctx, &dpdkproto.ListFirewallRulesRequest{
		InterfaceID: []byte(interfaceID),
	})
	if err != nil {
		return &api.FirewallRuleList{}, err
	}

	fwRules := make([]api.FirewallRule, len(res.GetRules()))
	for i, dpdkFwRule := range res.GetRules() {
		fwRule, err := api.ProtoFwRuleToFwRule(dpdkFwRule, interfaceID)
		if err != nil {
			return &api.FirewallRuleList{}, err
		}
		fwRules[i] = *fwRule
	}

	return &api.FirewallRuleList{
		TypeMeta: api.TypeMeta{Kind: api.FirewallRuleListKind},
		Items:    fwRules,
	}, nil
}

func (c *client) AddFirewallRule(ctx context.Context, fwRule *api.FirewallRule) (*api.FirewallRule, error) {
	var action, direction, ipv uint8

	switch strings.ToLower(fwRule.Spec.FirewallAction) {
	case "accept", "1":
		action = 1
		fwRule.Spec.FirewallAction = "Accept"
	case "drop", "0":
		action = 0
		fwRule.Spec.FirewallAction = "Drop"
	default:
		return &api.FirewallRule{}, fmt.Errorf("firewall action can be only: Drop = 0/Accept = 1")
	}

	switch strings.ToLower(fwRule.Spec.TrafficDirection) {
	case "ingress", "0":
		direction = 0
		fwRule.Spec.TrafficDirection = "Ingress"
	case "egress", "1":
		direction = 1
		fwRule.Spec.TrafficDirection = "Egress"
	default:
		return &api.FirewallRule{}, fmt.Errorf("traffic direction can be only: Ingress = 0/Egress = 1")
	}

	switch strings.ToLower(fwRule.Spec.IpVersion) {
	case "ipv4", "0":
		ipv = 0
		fwRule.Spec.IpVersion = "IPv4"
	case "ipv6", "1":
		ipv = 1
		fwRule.Spec.IpVersion = "IPv6"
	default:
		return &api.FirewallRule{}, fmt.Errorf("ip version can be only: IPv4 = 0/IPv6 = 1")
	}

	res, err := c.DPDKonmetalClient.AddFirewallRule(ctx, &dpdkproto.AddFirewallRuleRequest{
		InterfaceID: []byte(fwRule.FirewallRuleMeta.InterfaceID),
		Rule: &dpdkproto.FirewallRule{
			RuleID:    []byte(fwRule.FirewallRuleMeta.RuleID),
			Direction: dpdkproto.TrafficDirection(direction),
			Action:    dpdkproto.FirewallAction(action),
			Priority:  fwRule.Spec.Priority,
			IpVersion: dpdkproto.IPVersion(ipv),
			SourcePrefix: &dpdkproto.Prefix{
				IpVersion:    dpdkproto.IPVersion(ipv),
				Address:      []byte(fwRule.Spec.SourcePrefix.Addr().String()),
				PrefixLength: uint32(fwRule.Spec.SourcePrefix.Bits()),
			},
			DestinationPrefix: &dpdkproto.Prefix{
				IpVersion:    dpdkproto.IPVersion(ipv),
				Address:      []byte(fwRule.Spec.DestinationPrefix.Addr().String()),
				PrefixLength: uint32(fwRule.Spec.DestinationPrefix.Bits()),
			},
			ProtocolFilter: fwRule.Spec.ProtocolFilter,
		},
	})
	if err != nil {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res.Status)}, err
	}
	if res.Status.Error != 0 {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res.Status)}, apierrors.ErrServerError
	}

	return &api.FirewallRule{
		TypeMeta: api.TypeMeta{Kind: api.FirewallRuleKind},
		FirewallRuleMeta: api.FirewallRuleMeta{
			RuleID:      string(res.RuleID),
			InterfaceID: fwRule.InterfaceID,
		},
		Spec:   fwRule.Spec,
		Status: api.ProtoStatusToStatus(res.Status),
	}, nil
}

func (c *client) GetFirewallRule(ctx context.Context, ruleID string, interfaceID string) (*api.FirewallRule, error) {
	res, err := c.DPDKonmetalClient.GetFirewallRule(ctx, &dpdkproto.GetFirewallRuleRequest{
		InterfaceID: []byte(interfaceID),
		RuleID:      []byte(ruleID),
	})
	if err != nil {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res.Status)}, err
	}
	if errorCode := res.GetStatus().GetError(); errorCode != 0 {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res.Status)}, apierrors.ErrServerError
	}

	return api.ProtoFwRuleToFwRule(res.Rule, interfaceID)
}

func (c *client) DeleteFirewallRule(ctx context.Context, interfaceID string, ruleID string) (*api.FirewallRule, error) {
	res, err := c.DPDKonmetalClient.DeleteFirewallRule(ctx, &dpdkproto.DeleteFirewallRuleRequest{
		InterfaceID: []byte(interfaceID),
		RuleID:      []byte(ruleID),
	})
	if err != nil {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res)}, err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return &api.FirewallRule{Status: api.ProtoStatusToStatus(res)}, apierrors.ErrServerError
	}
	return &api.FirewallRule{Status: api.ProtoStatusToStatus(res)}, nil
}

func (c *client) Initialized(ctx context.Context) (string, error) {
	res, err := c.DPDKonmetalClient.Initialized(ctx, &dpdkproto.Empty{})
	if err != nil {
		return "", err
	}
	return res.Uuid, nil
}

func (c *client) Init(ctx context.Context, initConfig dpdkproto.InitConfig) error {
	res, err := c.DPDKonmetalClient.Init(ctx, &initConfig)
	if err != nil {
		return err
	}
	if errorCode := res.GetError(); errorCode != 0 {
		return apierrors.NewStatusError(errorCode, res.GetMessage())
	}
	return nil
}
