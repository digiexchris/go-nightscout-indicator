package unitconverter_test

import (
	"github.com/digiexchris/go-nightscout-indicator/unitconverter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUnitconverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unitconverter Suite")
}

/**
todo if this gets any more complicated of a package than this,
move these tests to their own test files that describe things better
*/

var _ = Describe("Convert units", func() {
	Context("When MMOL", func() {
		It("Displays correct values", func() {
			output := unitconverter.FormatTitle(unitconverter.MMOL, 100, 12, "Flat")
			Expect(output).To(Equal("5.6 (0.667 →)"))
		})
	})

	Context("When mg/dl", func() {
		It("Displays correct values", func() {
			output := unitconverter.FormatTitle(unitconverter.MGDL, 100, 12, "Flat")
			Expect(output).To(Equal("100 (12 →)"))
		})
	})
})
