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
	"fmt"
	"reflect"
)

type Scheme struct {
	typeByKind map[string]reflect.Type
	kindByType map[reflect.Type]string
}

func NewScheme() *Scheme {
	return &Scheme{
		typeByKind: make(map[string]reflect.Type),
		kindByType: make(map[reflect.Type]string),
	}
}

func (s *Scheme) Add(objs ...any) error {
	for _, obj := range objs {
		t := reflect.TypeOf(obj)
		if t.Kind() != reflect.Ptr {
			return fmt.Errorf("object %T must be a pointer to a struct", obj)
		}

		t = t.Elem()
		if err := s.AddWithKind(t.Name(), obj); err != nil {
			return fmt.Errorf("[name %s] %w", t.Name(), err)
		}
	}
	return nil
}

func (s *Scheme) AddWithKind(name string, obj any) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("object %T must be a pointer to a struct", obj)
	}

	s.typeByKind[name] = t.Elem()
	s.kindByType[t.Elem()] = name
	return nil
}

func (s *Scheme) KindFor(obj any) (string, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return "", fmt.Errorf("object %T must be a pointer to a struct", obj)
	}

	kind, ok := s.kindByType[t.Elem()]
	if !ok {
		return "", fmt.Errorf("no kind registered for type %T", obj)
	}
	return kind, nil
}

func (s *Scheme) New(name string) (any, error) {
	typ, ok := s.typeByKind[name]
	if !ok {
		return nil, fmt.Errorf("no type %q registered", name)
	}
	return reflect.New(typ).Interface(), nil
}
