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

package io_test

import (
	"bytes"

	. "github.com/onmetal/dpservice-go-library/io"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {
	Context("CheckpointReader", func() {
		var (
			data []byte
			rd   *CheckpointReader
		)
		BeforeEach(func() {
			data = []byte("foo.bar.baz.qux")
			rd = NewCheckpointReader(bytes.NewReader(data))
		})

		It("should allow unreading", func() {
			fooBytes := make([]byte, 3)
			n, err := rd.Read(fooBytes)
			Expect(n).To(Equal(3))
			Expect(err).NotTo(HaveOccurred())

			Expect(fooBytes).To(Equal([]byte("foo")))

			n, err = rd.Unread()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(3))

			fooBytes = make([]byte, 3)
			n, err = rd.Read(fooBytes)
			Expect(n).To(Equal(3))
			Expect(err).NotTo(HaveOccurred())
			Expect(fooBytes).To(Equal([]byte("foo")))
		})

		It("should allow unreading uneven amounts", func() {
			fooDotBarBytes := make([]byte, 7)
			n, err := rd.Read(fooDotBarBytes)
			Expect(n).To(Equal(7))
			Expect(err).NotTo(HaveOccurred())

			Expect(fooDotBarBytes).To(Equal([]byte("foo.bar")))

			n, err = rd.Unread()
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(7))

			fooBytes := make([]byte, 3)
			n, err = rd.Read(fooBytes)
			Expect(n).To(Equal(3))
			Expect(err).NotTo(HaveOccurred())

			dotBarBytes := make([]byte, 4)
			n, err = rd.Read(dotBarBytes)
			Expect(n).To(Equal(4))
			Expect(err).NotTo(HaveOccurred())
			Expect(dotBarBytes).To(Equal([]byte(".bar")))
		})
	})
})
