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

package dynamic

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	structured "github.com/onmetal/dpservice-cli/dpdk/client"
)

type ObjectKey interface {
	fmt.Stringer
	Name() string
}

type InterfaceKey struct {
	ID string
}

func (k InterfaceKey) String() string {
	return k.ID
}

func (k InterfaceKey) Name() string {
	return k.ID
}

type PrefixKey struct {
	InterfaceID string
	Prefix      netip.Prefix
}

func (k PrefixKey) String() string {
	return fmt.Sprintf("%s/%s", k.InterfaceID, k.Prefix)
}

func (k PrefixKey) Name() string {
	return k.Prefix.String()
}

type VirtualIPKey struct {
	InterfaceID string
}

func (k VirtualIPKey) String() string {
	return k.InterfaceID
}

func (k VirtualIPKey) Name() string {
	return k.InterfaceID
}

type RouteKey struct {
	VNI        uint32
	Prefix     netip.Prefix
	NextHopVNI uint32
	NextHopIP  netip.Addr
}

func (k RouteKey) String() string {
	return fmt.Sprintf("%d:%s-%d:%s", k.VNI, k.Prefix, k.NextHopVNI, k.NextHopIP)
}

func (k RouteKey) Name() string {
	return fmt.Sprintf("%s-%d:%s", k.Prefix, k.NextHopVNI, k.NextHopIP)
}

type emptyKey struct{}

func (emptyKey) String() string {
	return ""
}

func (emptyKey) Name() string {
	return ""
}

var EmptyKey ObjectKey = emptyKey{}

func ObjectKeyFromObject(obj any) ObjectKey {
	switch obj := obj.(type) {
	case *api.Interface:
		return InterfaceKey{ID: obj.ID}
	case *api.Prefix:
		return PrefixKey{
			InterfaceID: obj.InterfaceID,
			Prefix:      obj.Prefix,
		}
	case *api.Route:
		return RouteKey{
			VNI:        obj.VNI,
			Prefix:     obj.Prefix,
			NextHopVNI: obj.NextHop.VNI,
			NextHopIP:  *obj.NextHop.IP,
		}
	case *api.VirtualIP:
		return VirtualIPKey{
			InterfaceID: obj.InterfaceID,
		}
	default:
		return EmptyKey
	}
}

type Client interface {
	Create(ctx context.Context, obj any) error
	Delete(ctx context.Context, obj any) error
}

type client struct {
	structured structured.Client
}

func (c *client) Create(ctx context.Context, obj any) error {
	switch obj := obj.(type) {
	case *api.Interface:
		res, err := c.structured.CreateInterface(ctx, obj)
		if err != nil {
			return err
		}

		*obj = *res
		return nil
	case *api.Prefix:
		res, err := c.structured.AddPrefix(ctx, obj)
		if err != nil {
			return err
		}

		*obj = *res
		return nil
	case *api.Route:
		res, err := c.structured.AddRoute(ctx, obj)
		if err != nil {
			return err
		}

		*obj = *res
		return nil
	case *api.VirtualIP:
		res, err := c.structured.AddVirtualIP(ctx, obj)
		if err != nil {
			return err
		}

		*obj = *res
		return nil
	default:
		return fmt.Errorf("unsupported object %T", obj)
	}
}

func (c *client) Delete(ctx context.Context, obj any) error {
	switch obj := obj.(type) {
	case *api.Interface:
		_, err := c.structured.DeleteInterface(ctx, obj.ID)
		return err
	case *api.Prefix:
		_, err := c.structured.DeletePrefix(ctx, obj.InterfaceID, obj.Prefix)
		return err
	case *api.Route:
		_, err := c.structured.DeleteRoute(ctx, obj.VNI, obj.Prefix)
		return err
	case *api.VirtualIP:
		_, err := c.structured.DeleteVirtualIP(ctx, obj.InterfaceID)
		return err
	default:
		return fmt.Errorf("unsupported object %T", obj)
	}
}

func NewFromStructured(structured structured.Client) Client {
	return &client{structured: structured}
}
