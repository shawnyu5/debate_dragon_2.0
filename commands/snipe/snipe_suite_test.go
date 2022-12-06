package snipe_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSnipe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snipe Suite")
}
