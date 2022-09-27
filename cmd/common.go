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

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"time"

	"github.com/onmetal/dpservice-go-library/dpdk/client"
	"github.com/onmetal/dpservice-go-library/renderer"
	"github.com/onmetal/dpservice-go-library/sources"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DPDKClientFactory interface {
	NewClient(ctx context.Context) (client.Client, func() error, error)
}

type DPDKClientOptions struct {
	Address        string
	ConnectTimeout time.Duration
}

func (o *DPDKClientOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "address", "localhost:1337", "net-dpservice address.")
	fs.DurationVar(&o.ConnectTimeout, "connect-timeout", 4*time.Second, "Timeout to connect to the net-dpservice.")
}

func (o *DPDKClientOptions) NewClient(ctx context.Context) (client.Client, func() error, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, o.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to %s: %w", o.Address, err)
	}

	protoClient := dpdkproto.NewDPDKonmetalClient(conn)
	c := client.NewClient(protoClient)

	cleanup := conn.Close
	return c, cleanup, nil
}

func SubcommandRequired(cmd *cobra.Command, args []string) error {
	if err := cmd.Help(); err != nil {
		return err
	}
	return errors.New("subcommand is required")
}

func MultipleOfArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args)%n != 0 {
			return fmt.Errorf("expected a multiple of %d args but got %d args", n, len(args))
		}
		return nil
	}
}

func CommandNames(cmds []*cobra.Command) []string {
	res := make([]string, len(cmds))
	for i, cmd := range cmds {
		res[i] = cmd.Name()
	}
	return res
}

type RendererOptions struct {
	Output string
	Pretty bool
}

func (o *RendererOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Output, "output", "o", o.Output, "Output format.")
	fs.BoolVar(&o.Pretty, "pretty", o.Pretty, "Whether to render pretty output.")
}

func (o *RendererOptions) NewRenderer(w io.Writer) (renderer.Renderer, error) {
	// TODO: Factor out instantiation of registry & make it more modular.
	registry := renderer.NewRegistry()

	if err := registry.Register("json", func(w io.Writer) renderer.Renderer {
		return renderer.NewJSON(w, o.Pretty)
	}); err != nil {
		return nil, err
	}

	if err := registry.Register("yaml", func(w io.Writer) renderer.Renderer {
		return renderer.NewYAML(w)
	}); err != nil {
		return nil, err
	}

	return registry.New(o.Output, w)
}

type RendererFactory interface {
	NewRenderer(w io.Writer) (renderer.Renderer, error)
}

type SourcesOptions struct {
	Filename []string
}

func (o *SourcesOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVarP(&o.Filename, "filename", "f", o.Filename, "Filename, directory, or URL to file to use to create the resource")
}

func (o *SourcesOptions) NewIterator() (*sources.Iterator, error) {
	return sources.NewIterator(o.Filename), nil
}

type SourcesReaderFactory interface {
	NewIterator() (*sources.Iterator, error)
}

type RouteKey struct {
	Prefix     netip.Prefix
	NextHopVNI uint32
	NextHopIP  netip.Addr
}

func ParseRouteKey(prefixStr, nextHopVNIStr, nextHopIPStr string) (RouteKey, error) {
	prefix, err := netip.ParsePrefix(prefixStr)
	if err != nil {
		return RouteKey{}, fmt.Errorf("error parsing prefix: %w", err)
	}

	nextHopVNI, err := strconv.ParseUint(nextHopVNIStr, 10, 32)
	if err != nil {
		return RouteKey{}, fmt.Errorf("error parsing next hop vni: %w", err)
	}

	nextHopIP, err := netip.ParseAddr(nextHopIPStr)
	if err != nil {
		return RouteKey{}, fmt.Errorf("error parsing next hop ip: %w", err)
	}

	return RouteKey{
		Prefix:     prefix,
		NextHopVNI: uint32(nextHopVNI),
		NextHopIP:  nextHopIP,
	}, nil
}

func ParseRouteKeyArgs(args []string) ([]RouteKey, error) {
	if len(args)%3 != 0 {
		return nil, fmt.Errorf("expected args to be a multiple of 3 but got %d", len(args))
	}

	keys := make([]RouteKey, len(args)%3)
	for i := 0; i < len(args); i += 3 {
		key, err := ParseRouteKey(args[i], args[i+1], args[i+2])
		if err != nil {
			return nil, fmt.Errorf("[route key %d] %w", i, err)
		}

		keys[i/3] = key
	}
	return keys, nil
}

func ParsePrefixArgs(args []string) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, len(args))
	for i, arg := range args {
		prefix, err := netip.ParsePrefix(arg)
		if err != nil {
			return nil, fmt.Errorf("[prefix %d] %w", i, err)
		}

		prefixes[i] = prefix
	}
	return prefixes, nil
}
