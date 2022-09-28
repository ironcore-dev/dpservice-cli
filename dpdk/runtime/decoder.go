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

package runtime

import (
	"encoding/json"
	"fmt"
	"io"

	yaml2 "github.com/ghodss/yaml"
	dpsvcio "github.com/onmetal/dpservice-go-library/io"
	"gopkg.in/yaml.v2"
)

type Decoder interface {
	Decode(v any) error
}

type KindDecoder struct {
	scheme  *Scheme
	decoder PeekDecoder
}

func NewKindDecoder(scheme *Scheme, decoder PeekDecoder) *KindDecoder {
	return &KindDecoder{
		scheme:  scheme,
		decoder: decoder,
	}
}

type PeekDecoder interface {
	Decoder
	Undecode() error
}

type peekDecoder struct {
	decoded bool
	decoder Decoder
	reader  *dpsvcio.CheckpointReader
}

func (d *peekDecoder) Decode(v any) error {
	d.decoded = true
	d.reader.Checkpoint()
	return d.decoder.Decode(v)
}

func (d *peekDecoder) Undecode() error {
	if !d.decoded {
		return fmt.Errorf("must call Decode before Undecode")
	}
	d.decoded = false
	if _, err := d.reader.Unread(); err != nil {
		return err
	}
	return nil
}

func NewPeekDecoder(rd io.Reader, newDecoder func(rd io.Reader) Decoder) PeekDecoder {
	reader := dpsvcio.NewCheckpointReader(rd)
	return &peekDecoder{
		decoder: newDecoder(reader),
		reader:  reader,
	}
}

func (d *KindDecoder) Next() (any, error) {
	typeMeta := &struct {
		Kind string `json:"kind"`
	}{}
	if err := d.decoder.Decode(typeMeta); err != nil {
		return nil, fmt.Errorf("error decoding type meta: %w", err)
	}

	if err := d.decoder.Undecode(); err != nil {
		return nil, fmt.Errorf("error reversing decoder: %w", err)
	}

	res, err := d.scheme.New(typeMeta.Kind)
	if err != nil {
		return nil, fmt.Errorf("error creating new %s: %w", typeMeta.Kind, err)
	}

	if err = d.decoder.Decode(res); err != nil {
		return nil, fmt.Errorf("error decoding %s: %w", typeMeta.Kind, err)
	}
	return res, nil
}

func NewExtDecoderFactory(ext string) (func(reader io.Reader) Decoder, error) {
	switch ext {
	case "json", ".json":
		return func(rd io.Reader) Decoder {
			return json.NewDecoder(rd)
		}, nil
	case "yaml", ".yaml", "yml", ".yml":
		return func(rd io.Reader) Decoder {
			return NewYAMLToJSONDecoder(rd)
		}, nil
	default:
		return nil, fmt.Errorf("unsupported extension %q", ext)
	}
}

type YAMLToJSONDecoder struct {
	decoder *yaml.Decoder
}

func NewYAMLToJSONDecoder(rd io.Reader) *YAMLToJSONDecoder {
	return &YAMLToJSONDecoder{decoder: yaml.NewDecoder(rd)}
}

func (d *YAMLToJSONDecoder) Decode(v any) error {
	var yamlValue any
	if err := d.decoder.Decode(&yamlValue); err != nil {
		return err
	}

	yamlData, err := yaml2.Marshal(yamlValue)
	if err != nil {
		return err
	}

	jsonData, err := yaml2.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, v)
}
