// Copyright (c) 2017 Intel Corporation
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

package functional

import (
	"fmt"

	. "github.com/clearcontainers/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func withWorkload(workload string, expectedExitCode int) TableEntry {
	w := []string{"sh", "-c", workload}
	return Entry(fmt.Sprintf("with '%v' as workload", w), w, expectedExitCode)
}

func withoutOption(option string, fail bool) TableEntry {
	return Entry(fmt.Sprintf("without '%s' option", option), option, fail)
}

var _ = Describe("run", func() {
	var (
		container *Container
		err       error
	)

	BeforeEach(func() {
		container, err = NewContainer([]string{}, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(container).NotTo(BeNil())
	})

	AfterEach(func() {
		Expect(container.Teardown()).To(Succeed())
	})

	DescribeTable("container",
		func(workload []string, expectedExitCode int) {
			Expect(container.SetWorkload(workload)).To(Succeed())

			_, _, exitCode := container.Run()

			Expect(expectedExitCode).To(Equal(exitCode))
		},
		withWorkload("true", 0),
		withWorkload("false", 1),
		withWorkload("exit 0", 0),
		withWorkload("exit 1", 1),
		withWorkload("exit 15", 15),
		withWorkload("exit 123", 123),
	)
})

var _ = Describe("run", func() {
	var (
		container *Container
		err       error
	)

	BeforeEach(func() {
		container, err = NewContainer([]string{"true"}, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(container).NotTo(BeNil())
	})

	AfterEach(func() {
		Expect(container.Teardown()).To(Succeed())
	})

	DescribeTable("container",
		func(option string, fail bool) {
			Expect(container.RemoveOption(option)).To(Succeed())

			_, stderr, exitCode := container.Run()

			if fail {
				Expect(exitCode).ToNot(Equal(0))
				Expect(stderr).NotTo(BeEmpty())
			} else {
				Expect(exitCode).To(Equal(0))
				Expect(stderr).To(BeEmpty())
			}
		},
		withoutOption("--bundle", shouldFail),
		withoutOption("-b", shouldFail),
		withoutOption("--pid-file", shouldNotFail),
		withoutOption("--console", shouldNotFail),
	)
})
