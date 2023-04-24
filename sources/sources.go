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

package sources

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	runtime2 "github.com/onmetal/dpservice-cli/dpdk/runtime"
)

type Iterator struct {
	sources []string
	idx     int
}

func NewIterator(sources []string) *Iterator {
	return &Iterator{
		sources: sources,
	}
}

func (r *Iterator) Next() (Source, error) {
	if r.idx >= len(r.sources) {
		return nil, io.EOF
	}

	src, err := NewSource(r.sources[r.idx])
	r.idx++
	return src, err
}

func NewSource(source string) (Source, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("error parsing source: %w", err)
	}

	switch u.Scheme {
	case "", "file":
		stat, err := os.Stat(u.Path)
		if err != nil {
			return nil, fmt.Errorf("error stat-ing file source %s: %w", u.Path, err)
		}

		if stat.IsDir() {
			entries, err := os.ReadDir(u.Path)
			if err != nil {
				return nil, fmt.Errorf("error reading dir %s: %w", u.Path, err)
			}
			return &DirSource{
				path:    u.Path,
				entries: entries,
			}, nil
		}
		return &FileIterator{
			path: u.Path,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported source scheme %s", u.Scheme)
	}
}

type Source interface {
	Next() (ReadCloserExt, error)
}

type ReadCloserExt interface {
	io.ReadCloser
	Ext() string
}

type FileIterator struct {
	path string
	read bool
}

func (f *FileIterator) Next() (ReadCloserExt, error) {
	if f.read {
		return nil, io.EOF
	}

	f.read = true
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}

	return &FileSource{File: file, path: f.path}, nil
}

type FileSource struct {
	*os.File
	path string
}

func (f *FileSource) Ext() string {
	return filepath.Ext(f.path)
}

type DirSource struct {
	path    string
	entries []os.DirEntry

	idx int
}

func (s *DirSource) Next() (ReadCloserExt, error) {
	for ; s.idx < len(s.entries); s.idx++ {
		if s.entries[s.idx].IsDir() {
			continue
		}

		entry := s.entries[s.idx]
		s.idx++

		p := filepath.Join(s.path, entry.Name())
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %w", p, err)
		}
		return &FileSource{File: f, path: p}, nil
	}
	return nil, io.EOF
}

func IterateObjects(iterator *Iterator, scheme *runtime2.Scheme, f func(obj any) error) error {
	for {
		src, err := iterator.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		for {
			rce, err := src.Next()
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}

			newDecoder, err := runtime2.NewExtDecoderFactory(rce.Ext())
			if err != nil {
				return err
			}

			decoder := runtime2.NewKindDecoder(api.DefaultScheme, runtime2.NewPeekDecoder(rce, newDecoder))
			for {
				obj, err := decoder.Next()
				if err != nil {
					if err != io.EOF {
						return err
					}
					break
				}

				if err := f(obj); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CollectObjects(iterator *Iterator, scheme *runtime2.Scheme) ([]any, error) {
	var objs []any
	if err := IterateObjects(iterator, scheme, func(obj any) error {
		objs = append(objs, obj)
		return nil
	}); err != nil {
		return nil, err
	}
	return objs, nil
}
