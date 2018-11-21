package nightscoutclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNightscoutClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nightscout Client Suite")
}
