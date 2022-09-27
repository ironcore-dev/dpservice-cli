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

	"github.com/onmetal/dpservice-go-library/dpdk/api"
	apierrors "github.com/onmetal/dpservice-go-library/dpdk/api/errors"
	"github.com/onmetal/dpservice-go-library/netiputil"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
)

type Client interface {
	GetInterface(ctx context.Context, id string) (*api.Interface, error)
	ListInterfaces(ctx context.Context) (*api.InterfaceList, error)
	CreateInterface(ctx context.Context, iface *api.Interface) (*api.Interface, error)
	DeleteInterface(ctx context.Context, id string) error

	GetVirtualIP(ctx context.Context, interfaceID string) (*api.VirtualIP, error)
	CreateVirtualIP(ctx context.Context, virtualIP *api.VirtualIP) (*api.VirtualIP, error)
	DeleteVirtualIP(ctx context.Context, interfaceID string) error

	ListPrefixes(ctx context.Context, interfaceID string) (*api.PrefixList, error)
	CreatePrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error)
	DeletePrefix(ctx context.Context, interfaceID string, prefix netip.Prefix) error

	ListRoutes(ctx context.Context, vni uint32) (*api.RouteList, error)
	CreateRoute(ctx context.Context, route *api.Route) (*api.Route, error)
	DeleteRoute(ctx context.Context, vni uint32, prefix netip.Prefix, nextHopVNI uint32, nextHopIP netip.Addr) error
}

type client struct {
	dpdkproto.DPDKonmetalClient
}

func NewClient(protoClient dpdkproto.DPDKonmetalClient) Client {
	return &client{protoClient}
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

	var ifaces []api.Interface
	for _, dpdkIface := range res.GetInterfaces() {
		iface, err := api.ProtoInterfaceToInterface(dpdkIface)
		if err != nil {
			return nil, err
		}

		ifaces = append(ifaces, *iface)
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

	underlayIP, err := netip.ParseAddr(string(res.GetResponse().GetUnderlayRoute()))
	if err != nil {
		return nil, fmt.Errorf("error parsing underlay route: %w", err)
	}

	return &api.Interface{
		TypeMeta:      api.TypeMeta{Kind: api.InterfaceKind},
		InterfaceMeta: iface.InterfaceMeta,
		Spec:          iface.Spec, // TODO: Enable dynamic device allocation
		Status: api.InterfaceStatus{
			UnderlayIP: underlayIP,
		},
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

func (c *client) CreateVirtualIP(ctx context.Context, virtualIP *api.VirtualIP) (*api.VirtualIP, error) {
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

	return &api.VirtualIP{
		TypeMeta: api.TypeMeta{Kind: api.VirtualIPKind},
		Spec:     virtualIP.Spec,
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

	var prefixes []api.Prefix
	for _, dpdkPrefix := range res.GetPrefixes() {
		prefix, err := api.ProtoPrefixToPrefix(interfaceID, dpdkPrefix)
		if err != nil {
			return nil, err
		}

		prefixes = append(prefixes, *prefix)
	}

	return &api.PrefixList{
		TypeMeta: api.TypeMeta{Kind: api.PrefixListKind},
		Items:    prefixes,
	}, nil
}

func (c *client) CreatePrefix(ctx context.Context, prefix *api.Prefix) (*api.Prefix, error) {
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
	return &api.Prefix{
		TypeMeta: api.TypeMeta{Kind: api.PrefixKind},
		Spec:     prefix.Spec,
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

func (c *client) CreateRoute(ctx context.Context, route *api.Route) (*api.Route, error) {
	res, err := c.DPDKonmetalClient.AddRoute(ctx, &dpdkproto.VNIRouteMsg{
		Vni: &dpdkproto.VNIMsg{Vni: route.VNI},
		Route: &dpdkproto.Route{
			IpVersion: api.NetIPAddrToProtoIPVersion(route.NextHop.IP),
			Weight:    100,
			Prefix: &dpdkproto.Prefix{
				IpVersion:    api.NetIPAddrToProtoIPVersion(route.Prefix.Addr()),
				Address:      []byte(route.Prefix.String()),
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
		TypeMeta: api.TypeMeta{Kind: api.RouteKind},
		Spec:     route.Spec,
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
				Address:      []byte(prefix.String()),
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

	var routes []api.Route
	for _, dpdkRoute := range res.GetRoutes() {
		route, err := api.ProtoRouteToRoute(vni, dpdkRoute)
		if err != nil {
			return nil, err
		}

		routes = append(routes, *route)
	}

	return &api.RouteList{
		TypeMeta: api.TypeMeta{Kind: api.RouteListKind},
		Items:    routes,
	}, nil
}
