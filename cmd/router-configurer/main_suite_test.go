package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"code.cloudfoundry.org/cf-tcp-router/testutil"
	"code.cloudfoundry.org/cf-tcp-router/utils"
	"code.cloudfoundry.org/consuladapter/consulrunner"
	"code.cloudfoundry.org/routing-api"
	"code.cloudfoundry.org/routing-api/cmd/routing-api/test_helpers"
	routingtestrunner "code.cloudfoundry.org/routing-api/cmd/routing-api/testrunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	routerConfigurerPath    string
	routingAPIBinPath       string
	routerConfigurerPort    int
	haproxyConfigFile       string
	haproxyConfigBackupFile string
	haproxyBaseConfigFile   string

	consulRunner *consulrunner.ClusterRunner
	dbAllocator  test_helpers.DbAllocator

	etcdPort int

	routingAPIAddress string
	routingAPIArgs    routingtestrunner.Args
	routingAPIPort    uint16
	routingAPIIP      string
	routingApiClient  routing_api.Client

	dbEnv = os.Getenv("DB")
)

func TestRouterConfigurer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RouterConfigurer Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	routerConfigurer, err := gexec.Build("code.cloudfoundry.org/cf-tcp-router/cmd/router-configurer", "-race")
	Expect(err).NotTo(HaveOccurred())
	routingAPIBin, err := gexec.Build("code.cloudfoundry.org/routing-api/cmd/routing-api", "-race")
	Expect(err).NotTo(HaveOccurred())
	payload, err := json.Marshal(map[string]string{
		"router-configurer": routerConfigurer,
		"routing-api":       routingAPIBin,
	})

	Expect(err).NotTo(HaveOccurred())

	return payload
}, func(payload []byte) {
	context := map[string]string{}

	err := json.Unmarshal(payload, &context)
	Expect(err).NotTo(HaveOccurred())

	routerConfigurerPort = 7000 + GinkgoParallelNode()
	routerConfigurerPath = context["router-configurer"]
	routingAPIBinPath = context["routing-api"]

	setupConsul()
})

var _ = BeforeEach(func() {
	randomFileName := testutil.RandomFileName("haproxy_", ".cfg")
	randomBackupFileName := fmt.Sprintf("%s.bak", randomFileName)
	randomBaseFileName := testutil.RandomFileName("haproxy_base_", ".cfg")
	haproxyConfigFile = path.Join(os.TempDir(), randomFileName)
	haproxyConfigBackupFile = path.Join(os.TempDir(), randomBackupFileName)
	haproxyBaseConfigFile = path.Join(os.TempDir(), randomBaseFileName)

	err := utils.WriteToFile(
		[]byte(
			`global maxconn 4096
defaults
  log global
  timeout connect 300000
  timeout client 300000
  timeout server 300000
  maxconn 2000`),
		haproxyBaseConfigFile)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(utils.FileExists(haproxyBaseConfigFile)).To(BeTrue())

	err = utils.CopyFile(haproxyBaseConfigFile, haproxyConfigFile)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(utils.FileExists(haproxyConfigFile)).To(BeTrue())

	switch dbEnv {
	case "etcd":
		port := 4001 + GinkgoParallelNode()
		dbAllocator = test_helpers.NewEtcdAllocator(port)
	default:
		dbAllocator = test_helpers.NewMySQLAllocator()
	}
	dbId, err := dbAllocator.Create()
	Expect(err).NotTo(HaveOccurred())

	// test_helpers.
	// 	etcdUrl = fmt.Sprintf("http://127.0.0.1:%d", etcdPort)
	// etcdRunner = etcdstorerunner.NewETCDClusterRunner(etcdPort, 1, nil)
	// etcdRunner.Start()

	// etcdAdapter = etcdRunner.Adapter(nil)

	routingAPIPort = uint16(6900 + GinkgoParallelNode())
	routingAPIIP = "127.0.0.1"
	routingAPIAddress = fmt.Sprintf("http://%s:%d", routingAPIIP, routingAPIPort)

	routingAPIArgs = routingtestrunner.Args{
		Port:       routingAPIPort,
		IP:         routingAPIIP,
		ConfigPath: createConfig(dbId, consulRunner.URL()),
		DevMode:    true,
	}
	routingApiClient = routing_api.NewClient(routingAPIAddress, false)
})

var _ = AfterEach(func() {
	err := os.Remove(haproxyConfigFile)
	Expect(err).ShouldNot(HaveOccurred())

	os.Remove(haproxyConfigBackupFile)

	etcdAdapter.Disconnect()
	etcdRunner.Reset()
	etcdRunner.Stop()
})

var _ = SynchronizedAfterSuite(func() {
	teardownConsul()
}, func() {
	gexec.CleanupBuildArtifacts()
})

func createConfig(dbId, consulUrl string) string {
	var configBytes []byte
	configFilePath := fmt.Sprintf("/tmp/example_%d.yml", GinkgoParallelNode())

	switch dbEnv {
	case "etcd":
		etcdConfigStr := `log_guid: "my_logs"
uaa_verification_key: "-----BEGIN PUBLIC KEY-----

      MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d

      KVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX

      qHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug

      spULZVNRxq7veq/fzwIDAQAB

      -----END PUBLIC KEY-----"

debug_address: "1.2.3.4:1234"
metron_config:
  address: "1.2.3.4"
  port: "4567"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: "default-tcp"
  type: "tcp"
  reservable_ports: "1024-65535"
etcd:
  node_urls: ["%s"]
consul_cluster:
  servers: "%s"`
		configBytes = []byte(fmt.Sprintf(etcdConfigStr, dbId, consulUrl))
	default:
		mysqlConfigStr := `log_guid: "my_logs"
uaa_verification_key: "-----BEGIN PUBLIC KEY-----

      MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d

      KVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX

      qHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug

      spULZVNRxq7veq/fzwIDAQAB

      -----END PUBLIC KEY-----"

debug_address: "1.2.3.4:1234"
metron_config:
  address: "1.2.3.4"
  port: "4567"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: "default-tcp"
  type: "tcp"
  reservable_ports: "1024-65535"
sqldb:
  username: "root"
  password: "password"
  schema: "%s"
  port: 3306
  host: "localhost"
  type: "mysql"
consul_cluster:
  servers: "%s"`
		configBytes = []byte(fmt.Sprintf(mysqlConfigStr, dbId, consulUrl))
	}

	err := utils.WriteToFile(configBytes, configFilePath)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(utils.FileExists(configFilePath)).To(BeTrue())

	return configFilePath
}

func setupConsul() {
	consulRunner = consulrunner.NewClusterRunner(9001+GinkgoParallelNode()*consulrunner.PortOffsetLength, 1, "http")
	consulRunner.Start()
	consulRunner.WaitUntilReady()
}

func teardownConsul() {
	consulRunner.Stop()
}
