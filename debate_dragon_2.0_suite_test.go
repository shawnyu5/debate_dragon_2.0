package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDebateDragon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DebateDragon2.0 Suite")
}
