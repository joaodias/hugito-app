package hugo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHugo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hugo Suite")
}
