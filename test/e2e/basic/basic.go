package basic

import (
	"fmt"
	"time"

	"github.com/fatedier/frp/test/e2e/framework"
	"github.com/fatedier/frp/test/e2e/framework/consts"

	. "github.com/onsi/ginkgo"
)

var connTimeout = 5 * time.Second

var _ = Describe("[Feature: Basic]", func() {
	f := framework.NewDefaultFramework()

	Describe("TCP", func() {
		It("Expose a TCP echo server", func() {
			localPortName := framework.TCPEchoServerPort
			serverConf := consts.DefaultServerConfig
			clientConf := consts.DefaultClientConfig

			getProxyConf := func(proxyName string, portName string, extra string) string {
				return fmt.Sprintf(`
				[%s]
				type = tcp
				local_port = {{ .%s }}
				remote_port = {{ .%s }}
				`+extra, proxyName, localPortName, portName)
			}

			tests := []struct {
				proxyName   string
				portName    string
				extraConfig string
			}{
				{
					proxyName: "tcp",
					portName:  framework.GenPortName("TCP"),
				},
				{
					proxyName:   "tcp-with-encryption",
					portName:    framework.GenPortName("TCPWithEncryption"),
					extraConfig: "use_encryption = true",
				},
				{
					proxyName:   "tcp-with-compression",
					portName:    framework.GenPortName("TCPWithCompression"),
					extraConfig: "use_compression = true",
				},
				{
					proxyName: "tcp-with-encryption-and-compression",
					portName:  framework.GenPortName("TCPWithEncryptionAndCompression"),
					extraConfig: `
					use_encryption = true
					use_compression = true
					`,
				},
			}

			// build all client config
			for _, test := range tests {
				clientConf += getProxyConf(test.proxyName, test.portName, test.extraConfig) + "\n"
			}
			// run frps and frpc
			f.RunProcesses([]string{serverConf}, []string{clientConf})

			for _, test := range tests {
				framework.ExpectTCPRequest(f.UsedPorts[test.portName],
					[]byte(consts.TestString), []byte(consts.TestString),
					connTimeout, test.proxyName)
			}
		})
	})

	Describe("UDP", func() {
		It("Expose a UDP echo server", func() {
			localPortName := framework.UDPEchoServerPort
			serverConf := consts.DefaultServerConfig
			clientConf := consts.DefaultClientConfig

			getProxyConf := func(proxyName string, portName string, extra string) string {
				return fmt.Sprintf(`
				[%s]
				type = udp
				local_port = {{ .%s }}
				remote_port = {{ .%s }}
				`+extra, proxyName, localPortName, portName)
			}

			tests := []struct {
				proxyName   string
				portName    string
				extraConfig string
			}{
				{
					proxyName: "udp",
					portName:  framework.GenPortName("UDP"),
				},
				{
					proxyName:   "udp-with-encryption",
					portName:    framework.GenPortName("UDPWithEncryption"),
					extraConfig: "use_encryption = true",
				},
				{
					proxyName:   "udp-with-compression",
					portName:    framework.GenPortName("UDPWithCompression"),
					extraConfig: "use_compression = true",
				},
				{
					proxyName: "udp-with-encryption-and-compression",
					portName:  framework.GenPortName("UDPWithEncryptionAndCompression"),
					extraConfig: `
					use_encryption = true
					use_compression = true
					`,
				},
			}

			// build all client config
			for _, test := range tests {
				clientConf += getProxyConf(test.proxyName, test.portName, test.extraConfig) + "\n"
			}
			// run frps and frpc
			f.RunProcesses([]string{serverConf}, []string{clientConf})

			for _, test := range tests {
				framework.ExpectUDPRequest(f.UsedPorts[test.portName],
					[]byte(consts.TestString), []byte(consts.TestString),
					connTimeout, test.proxyName)
			}

		})
	})

	Describe("STCP", func() {
		It("Expose TCP echo server with STCP", func() {
			localPortName := framework.TCPEchoServerPort
			serverConf := consts.DefaultServerConfig
			clientServerConf := consts.DefaultClientConfig
			clientVisitorConf := consts.DefaultClientConfig

			correctSK := "abc"
			wrongSK := "123"

			getProxyServerConf := func(proxyName string, extra string) string {
				return fmt.Sprintf(`
				[%s]
				type = stcp
				role = server
				sk = %s
				local_port = {{ .%s }}
				`+extra, proxyName, correctSK, localPortName)
			}
			getProxyVisitorConf := func(proxyName string, portName, visitorSK, extra string) string {
				return fmt.Sprintf(`
				[%s]
				type = stcp
				role = visitor
				server_name = %s
				sk = %s
				bind_port = {{ .%s }}
				`+extra, proxyName, proxyName, visitorSK, portName)
			}

			tests := []struct {
				proxyName    string
				bindPortName string
				visitorSK    string
				extraConfig  string
				expectError  bool
			}{
				{
					proxyName:    "stcp-normal",
					bindPortName: framework.GenPortName("STCP"),
					visitorSK:    correctSK,
				},
				{
					proxyName:    "stcp-with-encryption",
					bindPortName: framework.GenPortName("STCPWithEncryption"),
					visitorSK:    correctSK,
					extraConfig:  "use_encryption = true",
				},
				{
					proxyName:    "stcp-with-compression",
					bindPortName: framework.GenPortName("STCPWithCompression"),
					visitorSK:    correctSK,
					extraConfig:  "use_compression = true",
				},
				{
					proxyName:    "stcp-with-encryption-and-compression",
					bindPortName: framework.GenPortName("STCPWithEncryptionAndCompression"),
					visitorSK:    correctSK,
					extraConfig: `
					use_encryption = true
					use_compression = true
					`,
				},
				{
					proxyName:    "stcp-with-error-sk",
					bindPortName: framework.GenPortName("STCPWithErrorSK"),
					visitorSK:    wrongSK,
					expectError:  true,
				},
			}

			// build all client config
			for _, test := range tests {
				clientServerConf += getProxyServerConf(test.proxyName, test.extraConfig) + "\n"
			}
			for _, test := range tests {
				clientVisitorConf += getProxyVisitorConf(test.proxyName, test.bindPortName, test.visitorSK, test.extraConfig) + "\n"
			}
			// run frps and frpc
			f.RunProcesses([]string{serverConf}, []string{clientServerConf, clientVisitorConf})

			for _, test := range tests {
				expectResp := []byte(consts.TestString)
				if test.expectError {
					framework.ExpectTCPRequestError(f.UsedPorts[test.bindPortName],
						[]byte(consts.TestString),
						connTimeout, test.proxyName)
					continue
				}

				framework.ExpectTCPRequest(f.UsedPorts[test.bindPortName],
					[]byte(consts.TestString), expectResp,
					connTimeout, test.proxyName)
			}
		})
	})
})
