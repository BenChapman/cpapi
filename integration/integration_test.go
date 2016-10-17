package integration_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	It("should reply to pings on the API", func() {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/ping", appPort))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
