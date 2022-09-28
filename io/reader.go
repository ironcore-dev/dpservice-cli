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

package io

import (
	"bytes"
	"io"
)

type CheckpointReader struct {
	offset int
	src    io.Reader
	buf    bytes.Buffer
}

func NewCheckpointReader(src io.Reader) *CheckpointReader {
	return &CheckpointReader{
		src: src,
	}
}

func (r *CheckpointReader) Read(p []byte) (n int, err error) {
	n, err = io.MultiReader(bytes.NewReader(r.buf.Bytes()[r.offset:]), io.TeeReader(r.src, &r.buf)).Read(p)
	r.offset += n
	return n, err
}

func (r *CheckpointReader) Checkpoint() {
	r.buf.Reset()
	r.offset = 0
}

func (r *CheckpointReader) Unread() (n int, err error) {
	n = r.offset
	r.offset = 0
	return n, nil
}
