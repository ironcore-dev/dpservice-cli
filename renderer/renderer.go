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

package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/onmetal/dpservice-cli/dpdk/api"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
)

type Renderer interface {
	Render(v any) error
}

type JSON struct {
	w      io.Writer
	pretty bool
}

func NewJSON(w io.Writer, pretty bool) *JSON {
	return &JSON{w, pretty}
}

func (j *JSON) Render(v any) error {
	enc := json.NewEncoder(j.w)
	if j.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}

type YAML struct {
	w io.Writer
}

func NewYAML(w io.Writer) *YAML {
	return &YAML{w}
}

func (y *YAML) Render(v any) error {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	data, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}

	_, err = io.Copy(y.w, bytes.NewReader(data))
	return err
}

type Name struct {
	w         io.Writer
	operation string
}

func NewName(w io.Writer, operation string) *Name {
	return &Name{
		w:         w,
		operation: operation,
	}
}

func (n *Name) Render(v any) error {
	objs, err := getObjs(v)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if err := n.renderObject(obj); err != nil {
			return err
		}
	}
	return nil
}

func (n *Name) renderObject(obj api.Object) error {
	var parts []string
	if kind := obj.GetKind(); kind != "" {
		parts = append(parts, fmt.Sprintf("%s/%s", strings.ToLower(kind), obj.GetName()))
	} else {
		parts = append(parts, obj.GetName())
	}

	if n.operation != "" {
		parts = append(parts, n.operation)
	}

	_, err := fmt.Fprintf(n.w, "%s\n", strings.Join(parts, " "))
	return err
}

func getObjs(v any) ([]api.Object, error) {
	switch v := v.(type) {
	case api.Object:
		return []api.Object{v}, nil
	case api.List:
		return v.GetItems(), nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

type Table struct {
	w              io.Writer
	tableConverter TableConverter
}

func NewTable(w io.Writer, converter TableConverter) *Table {
	return &Table{w, converter}
}

type TableData struct {
	Headers []any
	Columns [][]any
}

type TableConverter interface {
	ConvertToTable(v any) (*TableData, error)
}

type defaultTableConverter struct{}

var DefaultTableConverter = defaultTableConverter{}

func (t defaultTableConverter) ConvertToTable(v any) (*TableData, error) {
	switch obj := v.(type) {
	case *api.LoadBalancer:
		return t.loadBalancerTable(*obj)
	case *api.LoadBalancerTarget:
		return t.loadBalancerTargetTable([]api.LoadBalancerTarget{*obj})
	case *api.LoadBalancerTargetList:
		return t.loadBalancerTargetTable(obj.Items)
	case *api.Interface:
		return t.interfaceTable([]api.Interface{*obj})
	case *api.InterfaceList:
		return t.interfaceTable(obj.Items)
	case *api.Prefix:
		return t.prefixTable([]api.Prefix{*obj})
	case *api.PrefixList:
		return t.prefixTable(obj.Items)
	case *api.Route:
		return t.routeTable([]api.Route{*obj})
	case *api.RouteList:
		return t.routeTable(obj.Items)
	case *api.VirtualIP:
		return t.virtualIPTable([]api.VirtualIP{*obj})
	case *api.Nat:
		return t.natTable([]api.Nat{*obj})
	case *api.NeighborNat:
		return t.neighborNatTable([]api.NeighborNat{*obj})
	case *api.NatList:
		return t.natTable(obj.Items)
	case *api.FirewallRule:
		return t.fwruleTable([]api.FirewallRule{*obj})
	case *api.FirewallRuleList:
		return t.fwruleTable(obj.Items)
	case *api.Init:
		return t.initTable(*obj)
	case *api.Initialized:
		return t.initializedTable(*obj)
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}

func (t defaultTableConverter) loadBalancerTable(lb api.LoadBalancer) (*TableData, error) {
	headers := []any{"ID", "VNI", "LbVipIP", "Lbports", "UnderlayRoute", "Status"}

	columns := make([][]any, 1)

	var ports = make([]string, 0, len(lb.Spec.Lbports))
	for _, port := range lb.Spec.Lbports {
		p := dpdkproto.Protocol_name[int32(port.Protocol)] + "/" + strconv.Itoa(int(port.Port))
		ports = append(ports, p)
	}
	columns[0] = []any{lb.ID, lb.Spec.VNI, lb.Spec.LbVipIP, ports, lb.Spec.UnderlayRoute, lb.Status.String()}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) loadBalancerTargetTable(lbtargets []api.LoadBalancerTarget) (*TableData, error) {
	headers := []any{"LoadBalancerID", "IpVersion", "Address", "Status"}

	columns := make([][]any, len(lbtargets))
	for i, lbtarget := range lbtargets {
		columns[i] = []any{
			lbtarget.LoadBalancerTargetMeta.LoadbalancerID,
			api.NetIPAddrToProtoIPVersion(*lbtarget.Spec.TargetIP),
			lbtarget.Spec.TargetIP,
			lbtarget.Status.String(),
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) interfaceTable(ifaces []api.Interface) (*TableData, error) {
	var headers []any
	if ifaces[0].Spec.VirtualFunction == nil {
		headers = []any{"ID", "VNI", "Device", "IPs", "UnderlayRoute"}
	} else {
		headers = []any{"ID", "VNI", "Device", "IPs", "UnderlayRoute", "VirtualFunction"}
	}
	if len(ifaces) == 1 {
		headers = append(headers, "Status")
	}
	columns := make([][]any, len(ifaces))
	for i, iface := range ifaces {
		if ifaces[0].Spec.VirtualFunction == nil {
			columns[i] = []any{iface.ID, iface.Spec.VNI, iface.Spec.Device, iface.Spec.IPs, iface.Spec.UnderlayRoute}
		} else {
			columns[i] = []any{iface.ID, iface.Spec.VNI, iface.Spec.Device, iface.Spec.IPs, iface.Spec.UnderlayRoute, iface.Spec.VirtualFunction.String()}
		}
		if len(ifaces) == 1 {
			columns[i] = append(columns[i], iface.Status.String())
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) prefixTable(prefixes []api.Prefix) (*TableData, error) {
	headers := []any{"Prefix", "UnderlayRoute"}
	if len(prefixes) == 1 {
		headers = append(headers, "Status")
	}
	columns := make([][]any, len(prefixes))
	for i, prefix := range prefixes {
		columns[i] = []any{prefix.Prefix, prefix.Spec.UnderlayRoute}
		if len(prefixes) == 1 {
			columns[i] = append(columns[i], prefix.Status.String())
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) routeTable(routes []api.Route) (*TableData, error) {
	headers := []any{"Prefix", "VNI", "NextHopVNI", "NextHopIP"}
	if len(routes) == 1 {
		headers = append(headers, "Status")
	}
	columns := make([][]any, len(routes))
	for i, route := range routes {
		columns[i] = []any{route.Prefix, route.VNI, route.NextHop.VNI, route.NextHop.IP}
		if len(routes) == 1 {
			columns[i] = append(columns[i], route.Status.String())
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) virtualIPTable(virtualIPs []api.VirtualIP) (*TableData, error) {
	headers := []any{"InterfaceID", "VirtualIP", "UnderlayRoute", "Status"}

	columns := make([][]any, len(virtualIPs))
	for i, virtualIP := range virtualIPs {
		columns[i] = []any{virtualIP.InterfaceID, virtualIP.IP, virtualIP.Spec.UnderlayRoute, virtualIP.Status.String()}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) natTable(nats []api.Nat) (*TableData, error) {
	var headers []any
	if nats[0].Spec.UnderlayRoute == nil {
		headers = []any{"NatIP", "MinPort", "MaxPort", "Status"}
	} else if nats[0].NatMeta.InterfaceID == "" {
		headers = []any{"NatIP", "MinPort", "MaxPort", "UnderlayRoute", "Status"}
	} else {
		headers = []any{"InterfaceID", "NatIP", "MinPort", "MaxPort", "UnderlayRoute", "Status"}
	}

	columns := make([][]any, len(nats))
	for i, nat := range nats {
		if nats[0].Spec.UnderlayRoute == nil {
			columns[i] = []any{nat.Spec.NatVIPIP, nat.Spec.MinPort, nat.Spec.MaxPort, nat.Status.String()}
		} else if nats[0].NatMeta.InterfaceID == "" {
			columns[i] = []any{nat.Spec.NatVIPIP, nat.Spec.MinPort, nat.Spec.MaxPort, nat.Spec.UnderlayRoute, nat.Status.String()}
		} else {
			columns[i] = []any{nat.NatMeta.InterfaceID, nat.Spec.NatVIPIP, nat.Spec.MinPort, nat.Spec.MaxPort, nat.Spec.UnderlayRoute, nat.Status.String()}
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) neighborNatTable(nats []api.NeighborNat) (*TableData, error) {

	headers := []any{"VNI", "NatIP", "MinPort", "MaxPort", "UnderlayRoute", "Status"}

	columns := make([][]any, len(nats))
	for i, nat := range nats {

		columns[i] = []any{nat.Spec.Vni, nat.NeighborNatMeta.NatVIPIP, nat.Spec.MinPort, nat.Spec.MaxPort, nat.Spec.UnderlayRoute, nat.Status.String()}

	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) fwruleTable(fwrules []api.FirewallRule) (*TableData, error) {
	headers := []any{"InterfaceID", "RuleID", "Direction", "Src", "Dst", "Action", "Protocol", "Priority"}
	if len(fwrules) == 1 {
		headers = append(headers, "Status")
	}
	columns := make([][]any, len(fwrules))
	for i, fwrule := range fwrules {
		columns[i] = []any{
			fwrule.FirewallRuleMeta.InterfaceID,
			fwrule.FirewallRuleMeta.RuleID,
			fwrule.Spec.TrafficDirection,
			fwrule.Spec.SourcePrefix,
			fwrule.Spec.DestinationPrefix,
			fwrule.Spec.FirewallAction,
			fwrule.Spec.ProtocolFilter.String(),
			fwrule.Spec.Priority,
		}
		if len(fwrules) == 1 {
			columns[i] = append(columns[i], fwrule.Status.String())
		}
	}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) initTable(init api.Init) (*TableData, error) {
	headers := []any{"Error", "Message"}
	columns := make([][]any, 1)
	columns[0] = []any{init.Status.Error, init.Status.Message}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

func (t defaultTableConverter) initializedTable(initialized api.Initialized) (*TableData, error) {
	headers := []any{"UUID", "Error", "Message"}
	columns := make([][]any, 1)
	columns[0] = []any{initialized.Spec.UUID, initialized.Status.Error, initialized.Status.Message}

	return &TableData{
		Headers: headers,
		Columns: columns,
	}, nil
}

var (
	lightBoxStyle = table.BoxStyle{
		BottomLeft:       "",
		BottomRight:      "",
		BottomSeparator:  "",
		EmptySeparator:   " ",
		Left:             "",
		LeftSeparator:    "",
		MiddleHorizontal: "",
		MiddleSeparator:  "",
		MiddleVertical:   " ",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "",
		RightSeparator:   "",
		TopLeft:          "",
		TopRight:         "",
		TopSeparator:     "",
		UnfinishedRow:    "",
	}
	tableStyle = table.Style{Box: lightBoxStyle}
)

func (t *Table) Render(v any) error {
	data, err := t.tableConverter.ConvertToTable(v)
	if err != nil {
		return err
	}

	tw := table.NewWriter()
	tw.SetStyle(tableStyle)
	tw.SetOutputMirror(t.w)

	tw.AppendHeader(data.Headers)
	for _, col := range data.Columns {
		tw.AppendRow(col)
	}

	tw.Render()
	return nil
}

type NewFunc func(w io.Writer) Renderer

type Registry struct {
	newFuncByName map[string]NewFunc
}

func NewRegistry() *Registry {
	return &Registry{
		newFuncByName: make(map[string]NewFunc),
	}
}

func (r *Registry) Register(name string, newFunc NewFunc) error {
	if _, ok := r.newFuncByName[name]; ok {
		return fmt.Errorf("renderer %q is already registered", name)
	}

	r.newFuncByName[name] = newFunc
	return nil
}

func (r *Registry) New(name string, w io.Writer) (Renderer, error) {
	newFunc, ok := r.newFuncByName[name]
	if !ok {
		return nil, fmt.Errorf("unknown renderer %q", name)
	}

	return newFunc(w), nil
}
