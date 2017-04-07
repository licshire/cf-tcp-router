package config_test

import (
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/cf-tcp-router/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		expectedCfg config.Config
		configYaml  string
		configFile  *os.File
	)
	BeforeEach(func() {
		expectedCfg = config.Config{
			OAuth: config.OAuthConfig{
				TokenEndpoint:     "uaa.service.cf.internal",
				ClientName:        "someclient",
				ClientSecret:      "somesecret",
				Port:              8443,
				SkipSSLValidation: true,
				CACerts:           "some-ca-cert",
			},
			RoutingAPI: config.RoutingAPIConfig{
				URI:          "http://routing-api.service.cf.internal",
				Port:         3000,
				AuthDisabled: false,
			},
			HaProxyPidFile:  "/path/to/pid/file",
			RouterGroupName: "some-router-group",
		}
	})

	JustBeforeEach(func() {
		var err error
		configFile, err = ioutil.TempFile("", "test_config")
		Expect(err).NotTo(HaveOccurred())

		count, err := configFile.Write([]byte(configYaml))
		Expect(err).NotTo(HaveOccurred())
		Expect(count).NotTo(Equal(0))

		err = configFile.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(configFile.Name()) // clean up
	})

	Context("when a valid config", func() {
		BeforeEach(func() {
			configYaml = `
oauth:
  token_endpoint: "uaa.service.cf.internal"
  client_name: "someclient"
  client_secret: "somesecret"
  port: 8443
  skip_ssl_validation: true
  ca_certs: "some-ca-cert"

routing_api:
  uri: http://routing-api.service.cf.internal
  port: 3000
  auth_disabled: false

haproxy_pid_file: /path/to/pid/file
router_group: "some-router-group"
`
		})

		It("loads the config", func() {
			cfg, err := config.New(configFile.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(*cfg).To(Equal(expectedCfg))
		})
	})

	Context("when given an invalid config", func() {
		Context("non existing config", func() {
			It("return error", func() {
				_, err := config.New("non_existing_config.yml")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("malformed YAML config", func() {
			BeforeEach(func() {
				configYaml = `

token_endpoint_invalid_porperty: "uaa.service.cf.internal"
  client_name: "someclient"
  client_secret: "somesecret"
  port: 8080

  routing_api:
    uri: http://routing-api.service.cf.internal
    port: 3000

haproxy_pid_file: /path/to/pid/file
router_group: "some-router-group"
`
			})

			It("returns error", func() {
				_, err := config.New(configFile.Name())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when haproxy pid file is missing", func() {
		BeforeEach(func() {
			configYaml = `
oauth:
  token_endpoint: "uaa.service.cf.internal"
  client_name: "someclient"
  client_secret: "somesecret"
  port: 8443
  skip_ssl_validation: true
  ca_certs: "some-ca-cert"

routing_api:
  uri: http://routing-api.service.cf.internal
  port: 3000
  auth_disabled: false

router_group: "some-router-group"
`
		})

		It("returns error", func() {
			_, err := config.New(configFile.Name())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when router group name is missing", func() {
		BeforeEach(func() {
			configYaml = `
oauth:
  token_endpoint: "uaa.service.cf.internal"
  client_name: "someclient"
  client_secret: "somesecret"
  port: 8443
  skip_ssl_validation: true
  ca_certs: "some-ca-cert"

routing_api:
  uri: http://routing-api.service.cf.internal
  port: 3000
  auth_disabled: false

haproxy_pid_file: /path/to/pid/file
`
		})
		It("returns error", func() {
			_, err := config.New(configFile.Name())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when oauth section is  missing", func() {
		BeforeEach(func() {
			configYaml = `
routing_api:
  uri: http://routing-api.service.cf.internal
  port: 3000

haproxy_pid_file: /path/to/pid/file
router_group: "some-router-group"
`
		})

		It("loads only routing api section", func() {
			expectedCfg = config.Config{
				RoutingAPI: config.RoutingAPIConfig{
					URI:  "http://routing-api.service.cf.internal",
					Port: 3000,
				},
				HaProxyPidFile:  "/path/to/pid/file",
				RouterGroupName: "some-router-group",
			}
			cfg, err := config.New(configFile.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(*cfg).To(Equal(expectedCfg))
		})
	})

	Context("when oauth section has some missing fields", func() {
		BeforeEach(func() {
			configYaml = `
oauth:
  token_endpoint: "uaa.service.cf.internal"
  client_name:
  client_secret:
  port: 8443
  skip_ssl_validation: true

routing_api:
  uri: http://routing-api.service.cf.internal
  port: 3000

haproxy_pid_file: /path/to/pid/file
router_group: some-router-group
`
		})

		It("loads config and defaults missing fields", func() {
			expectedCfg := config.Config{
				OAuth: config.OAuthConfig{
					TokenEndpoint:     "uaa.service.cf.internal",
					ClientName:        "",
					ClientSecret:      "",
					Port:              8443,
					SkipSSLValidation: true,
				},
				RoutingAPI: config.RoutingAPIConfig{
					URI:  "http://routing-api.service.cf.internal",
					Port: 3000,
				},
				HaProxyPidFile:  "/path/to/pid/file",
				RouterGroupName: "some-router-group",
			}
			cfg, err := config.New(configFile.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(*cfg).To(Equal(expectedCfg))
		})
	})
})
