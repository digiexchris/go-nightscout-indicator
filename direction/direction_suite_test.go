package direction_test

import (
	"github.com/digiexchris/go-nightscout-indicator/direction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDirection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Direction Suite")
}

var _ = Describe("Get Direction String", func() {
	Context("When the value is in the list", func() {
		It("Returns the correct string", func() {
			Expect(direction.GetDirectionForTrend("NONE")).To(Equal("⇼"))
			Expect(direction.GetDirectionForTrend("DoubleUp")).To(Equal("⇈"))
			Expect(direction.GetDirectionForTrend("SingleUp")).To(Equal("↑"))
			Expect(direction.GetDirectionForTrend("FortyFiveUp")).To(Equal("↗"))
			Expect(direction.GetDirectionForTrend("Flat")).To(Equal("→"))
			Expect(direction.GetDirectionForTrend("FortyFiveDown")).To(Equal("↘"))
			Expect(direction.GetDirectionForTrend("SingleDown")).To(Equal("↓"))
			Expect(direction.GetDirectionForTrend("DoubleDown")).To(Equal("⇊"))

		})
	})
})
