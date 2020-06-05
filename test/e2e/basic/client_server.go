package basic

import (
	"github.com/fatedier/frp/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Feature: Client-Server]", func() {
	f := framework.NewDefaultFramework()

	Describe("Protocol", func() {
		It("KCP", func() {
			// TODO
			_ = f
		})
		It("Websocket", func() {
			// TODO
		})
	})
})
