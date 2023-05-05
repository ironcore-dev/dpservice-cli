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
	"strconv"
	"strings"

	proto "github.com/onmetal/net-dpservice-go/proto"
)

func ProtoLoadBalancerToLoadBalancer(dpdkLB *proto.GetLoadBalancerResponse, lbID string) (*LoadBalancer, error) {

	var underlayRoute netip.Addr
	if underlayRouteString := string(dpdkLB.GetUnderlayRoute()); underlayRouteString != "" {
		var err error
		underlayRoute, err = netip.ParseAddr(string(dpdkLB.GetUnderlayRoute()))
		if err != nil {
			return nil, fmt.Errorf("error parsing underlay ip: %w", err)
		}
	}
	var lbip netip.Addr
	if lbipString := string(dpdkLB.GetLbVipIP().Address); lbipString != "" {
		var err error
		lbip, err = netip.ParseAddr(string(dpdkLB.GetLbVipIP().Address))
		if err != nil {
			return nil, fmt.Errorf("error parsing lb ip: %w", err)
		}
	}
	var lbports = make([]LBPort, 0, len(dpdkLB.Lbports))
	var p LBPort
	for _, lbport := range dpdkLB.Lbports {
		p.Protocol = uint32(lbport.Protocol)
		p.Port = lbport.Port
		lbports = append(lbports, p)
	}

	return &LoadBalancer{
		TypeMeta: TypeMeta{
			Kind: LoadBalancerKind,
		},
		LoadBalancerMeta: LoadBalancerMeta{
			ID: lbID,
		},
		Spec: LoadBalancerSpec{
			VNI:           dpdkLB.Vni,
			LbVipIP:       lbip,
			Lbports:       lbports,
			UnderlayRoute: underlayRoute,
		},
		Status: LoadBalancerStatus{
			Error:   dpdkLB.Status.Error,
			Message: dpdkLB.Status.Message,
		},
	}, nil
}

func LbipToProtoLbip(lbip netip.Addr) *proto.LBIP {
	return &proto.LBIP{IpVersion: NetIPAddrToProtoIPVersion(lbip), Address: []byte(lbip.String())}
}

func ProtoLbipToLbip(protolbip proto.LBIP) *LBIP {
	var ip netip.Addr
	if lbipString := string(protolbip.Address); lbipString != "" {
		var err error
		ip, err = netip.ParseAddr(string(protolbip.Address))
		if err != nil {
			return nil
		}
	}
	return &LBIP{Address: ip, IpVersion: proto.IPVersion_name[int32(protolbip.IpVersion)]}
}

func StringLbportToLbport(lbport string) (LBPort, error) {
	p := strings.Split(lbport, "/")
	protocolName := strings.ToLower(p[0])
	switch protocolName {
	case "icmp", "tcp", "udp", "sctp":
		protocolName = strings.ToUpper(protocolName)
	case "icmpv6":
		protocolName = "ICMPv6"
	default:
		return LBPort{}, fmt.Errorf("unsupported protocol")
	}
	protocol := proto.Protocol_value[protocolName]
	port, err := strconv.Atoi(p[1])
	if err != nil {
		return LBPort{}, fmt.Errorf("error parsing port number: %w", err)
	}
	return LBPort{Protocol: uint32(protocol), Port: uint32(port)}, nil
}

func ProtoInterfaceToInterface(dpdkIface *proto.Interface) (*Interface, error) {
	var ips []netip.Addr

	if ipv4String := string(dpdkIface.GetPrimaryIPv4Address()); ipv4String != "" {
		ip, err := netip.ParseAddr(ipv4String)
		if err != nil {
			return nil, fmt.Errorf("error parsing primary ipv4: %w", err)
		}

		ips = append(ips, ip)
	}

	if ipv6String := string(dpdkIface.GetPrimaryIPv6Address()); ipv6String != "" {
		ip, err := netip.ParseAddr(ipv6String)
		if err != nil {
			return nil, fmt.Errorf("error parsing primary ipv6: %w", err)
		}

		ips = append(ips, ip)
	}

	var underlayIP netip.Addr
	if underlayIPString := string(dpdkIface.GetUnderlayRoute()); underlayIPString != "" {
		var err error
		underlayIP, err = netip.ParseAddr(string(dpdkIface.GetUnderlayRoute()))
		if err != nil {
			return nil, fmt.Errorf("error parsing underlay ip: %w", err)
		}
	}

	return &Interface{
		TypeMeta: TypeMeta{
			Kind: InterfaceKind,
		},
		InterfaceMeta: InterfaceMeta{
			ID: string(dpdkIface.InterfaceID),
		},
		Spec: InterfaceSpec{
			VNI:    dpdkIface.GetVni(),
			Device: dpdkIface.GetPciDpName(),
			IPs:    ips,
		},
		Status: InterfaceStatus{
			UnderlayIP: underlayIP,
		},
	}, nil
}

func NetIPAddrToProtoIPVersion(addr netip.Addr) proto.IPVersion {
	switch {
	case addr.Is4():
		return proto.IPVersion_IPv4
	case addr.Is6():
		return proto.IPVersion_IPv6
	default:
		return 0
	}
}

func NetIPAddrToProtoIPConfig(addr netip.Addr) *proto.IPConfig {
	if !addr.IsValid() {
		return nil
	}

	return &proto.IPConfig{
		IpVersion:      NetIPAddrToProtoIPVersion(addr),
		PrimaryAddress: []byte(addr.String()),
	}
}

func ProtoVirtualIPToVirtualIP(interfaceID string, dpdkVIP *proto.InterfaceVIPIP) (*VirtualIP, error) {
	ip, err := netip.ParseAddr(string(dpdkVIP.GetAddress()))
	if err != nil {
		return nil, fmt.Errorf("error parsing virtual ip address: %w", err)
	}

	return &VirtualIP{
		TypeMeta: TypeMeta{
			Kind: VirtualIPKind,
		},
		VirtualIPMeta: VirtualIPMeta{
			InterfaceID: interfaceID,
			IP:          ip,
		},
		Spec: VirtualIPSpec{},
	}, nil
}

func ProtoPrefixToPrefix(interfaceID string, dpdkPrefix *proto.Prefix) (*Prefix, error) {
	addr, err := netip.ParseAddr(string(dpdkPrefix.GetAddress()))
	if err != nil {
		return nil, fmt.Errorf("error parsing dpdk prefix address: %w", err)
	}

	prefix, err := addr.Prefix(int(dpdkPrefix.PrefixLength))
	if err != nil {
		return nil, fmt.Errorf("invalid dpdk prefix length %d for address %s", dpdkPrefix.PrefixLength, addr)
	}

	return &Prefix{
		TypeMeta: TypeMeta{
			Kind: PrefixKind,
		},
		PrefixMeta: PrefixMeta{
			InterfaceID: interfaceID,
			Prefix:      prefix,
		},
		Spec: PrefixSpec{UnderlayRoute: dpdkPrefix.UnderlayRoute},
	}, nil
}

func ProtoRouteToRoute(vni uint32, dpdkRoute *proto.Route) (*Route, error) {
	prefixAddr, err := netip.ParseAddr(string(dpdkRoute.GetPrefix().GetAddress()))
	if err != nil {
		return nil, fmt.Errorf("error parsing prefix address: %w", err)
	}

	prefix := netip.PrefixFrom(prefixAddr, int(dpdkRoute.GetPrefix().GetPrefixLength()))

	nextHopIP, err := netip.ParseAddr(string(dpdkRoute.GetNexthopAddress()))
	if err != nil {
		return nil, fmt.Errorf("error parsing netxt hop address: %w", err)
	}

	return &Route{
		TypeMeta: TypeMeta{
			RouteKind,
		},
		RouteMeta: RouteMeta{
			VNI:    vni,
			Prefix: prefix,
			NextHop: RouteNextHop{
				VNI: dpdkRoute.GetNexthopVNI(),
				IP:  nextHopIP,
			},
		},
		Spec: RouteSpec{},
	}, nil
}

func ProtoLBPrefixToProtoPrefix(lbprefix proto.LBPrefix) *proto.Prefix {
	return &proto.Prefix{
		IpVersion:     lbprefix.IpVersion,
		Address:       lbprefix.Address,
		PrefixLength:  lbprefix.PrefixLength,
		UnderlayRoute: lbprefix.UnderlayRoute,
	}
}

func ProtoNatToNat(dpdkNat *proto.GetNATResponse, interfaceID string) (*Nat, error) {
	var underlayRoute netip.Addr
	if underlayRouteString := string(dpdkNat.GetUnderlayRoute()); underlayRouteString != "" {
		var err error
		underlayRoute, err = netip.ParseAddr(string(dpdkNat.GetUnderlayRoute()))
		if err != nil {
			return nil, fmt.Errorf("error parsing underlay ip: %w", err)
		}
	}
	var natvipip netip.Addr
	if natvipipString := string(dpdkNat.GetNatVIPIP().Address); natvipipString != "" {
		var err error
		natvipip, err = netip.ParseAddr(string(dpdkNat.GetNatVIPIP().Address))
		if err != nil {
			return nil, fmt.Errorf("error parsing nat ip: %w", err)
		}
	}

	return &Nat{
		TypeMeta: TypeMeta{
			Kind: NatKind,
		},
		NatMeta: NatMeta{
			InterfaceID: interfaceID,
		},
		Spec: NatSpec{
			NatVIPIP:      natvipip,
			MinPort:       dpdkNat.MinPort,
			MaxPort:       dpdkNat.MaxPort,
			UnderlayRoute: underlayRoute,
		},
		Status: NatStatus{
			Error:   dpdkNat.Status.Error,
			Message: dpdkNat.Status.Message,
		},
	}, nil
}
