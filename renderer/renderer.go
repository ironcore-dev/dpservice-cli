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

	"github.com/ghodss/yaml"
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
