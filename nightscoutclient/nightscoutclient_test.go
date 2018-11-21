package nightscoutclient_test

import (
	"errors"
	"github.com/digiexchris/go-nightscout-indicator/configuration"
	"github.com/digiexchris/go-nightscout-indicator/nightscoutclient"
	"github.com/digiexchris/go-nightscout-indicator/nightscoutclient/nightscoutclientfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
)

var _ = Describe("Get", func() {

	BeforeEach(func() {
		configuration.App = configuration.Config{
			NightscoutHost: httptest.DefaultRemoteAddr,
			ApiSecret:      "1234567890",
		}
	})

	Context("With a good response", func() {
		fakeHttpClient := &nightscoutclientfakes.FakeHttpClient{}
		w := httptest.NewRecorder()
		w.Body.WriteString("[{\"_id\":\"5bf594b187233bfdbeb3f2e3\",\"device\":\"xDrip-DexcomG5 G5 Native\",\"date\":1542821036906,\"dateString\":\"2018-11-21T11:23:56.906-0600\",\"sgv\":122,\"delta\":-0.5,\"direction\":\"Flat\",\"type\":\"sgv\",\"filtered\":177376,\"unfiltered\":177152,\"rssi\":100,\"noise\":1,\"sysTime\":\"2018-11-21T11:23:56.906-0600\"}]")

		fakeHttpClient.DoReturns(w.Result(), nil)

		nightscountClient := nightscoutclient.Client{
			HttpClient: fakeHttpClient,
		}

		reading := nightscountClient.Get("1234", "5678")

		It("Does not error", func() {
			Expect(reading.Error).To(BeNil())
		})

		It("Contains correct values", func() {
			Expect(reading.SGV).To(Equal(float32(122)))
			Expect(reading.Delta).To(Equal(float32(-0.5)))
			Expect(reading.Direction).To(Equal("Flat"))
		})
	})

	Context("With a bad response", func() {
		fakeHttpClient := &nightscoutclientfakes.FakeHttpClient{}
		w := httptest.NewRecorder()
		w.Body.WriteString("[{\"_id\":\"5bf594b187233bfdbeb3f2e3\",\"device\":\"xDrip-DexcomG5 G5 Native\",\"date\":1542821036906,\"dateString\":\"2018-11-21T11:23:56.906-0600\",\"sgv\":122,\"delta\":-0.5,\"direction\":\"Flat\",\"type\":\"sgv\",\"filtered\":177376,\"unfiltered\":177152,\"rssi\":100,\"noise\":1,\"sysTime\":\"2018-11-21T11:23:56.906-0600\"}]")

		fakeHttpClient.DoReturns(w.Result(), errors.New("An error has occurred"))

		nightscountClient := nightscoutclient.Client{
			HttpClient: fakeHttpClient,
		}

		reading := nightscountClient.Get("1234", "5678")

		It("Returns an error", func() {
			Expect(reading.Error).To(Not(BeNil()))
		})

		It("Contains zero values", func() {
			Expect(reading.SGV).To(Equal(float32(0)))
			Expect(reading.Delta).To(Equal(float32(0)))
			Expect(reading.Direction).To(Equal(""))
		})
	})

	Context("With a bad json", func() {
		fakeHttpClient := &nightscoutclientfakes.FakeHttpClient{}
		w := httptest.NewRecorder()
		w.Body.WriteString("oops")

		fakeHttpClient.DoReturns(w.Result(), nil)

		nightscountClient := nightscoutclient.Client{
			HttpClient: fakeHttpClient,
		}

		reading := nightscountClient.Get("1234", "5678")

		It("Returns an error", func() {
			Expect(reading.Error).To(Not(BeNil()))
		})

		It("Contains zero values", func() {
			Expect(reading.SGV).To(Equal(float32(0)))
			Expect(reading.Delta).To(Equal(float32(0)))
			Expect(reading.Direction).To(Equal(""))
		})
	})
})
