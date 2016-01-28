package metrics_reporter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/cf-tcp-router/metrics_reporter/haproxy_client"
	"github.com/cloudfoundry-incubator/cf-tcp-router/models"
)

func Convert(proxyStats haproxy_client.HaproxyStats) *HaProxyMetricsReport {

	if len(proxyStats) == 0 {
		return nil
	}

	frontendMetricsReport := MetricsReport{ProxyMetrics: make(map[models.RoutingKey]ProxyStats)}
	backendMetricsReport := MetricsReport{ProxyMetrics: make(map[models.RoutingKey]ProxyStats)}

	var frontendCount, backendCount uint64

	for _, proxyStat := range proxyStats {
		switch proxyStat.Type {
		case 0:
			frontendCount++
			frontendMetricsReport.BytesIn += proxyStat.BytesIn
			frontendMetricsReport.BytesOut += proxyStat.BytesOut
			frontendMetricsReport.MaxSessionPerSec += proxyStat.MaxSessionPerSec
			populateProxyStats(proxyStat, frontendMetricsReport.ProxyMetrics)

		case 1:
			backendCount++
			backendMetricsReport.BytesIn += proxyStat.BytesIn
			backendMetricsReport.BytesOut += proxyStat.BytesOut
			backendMetricsReport.MaxSessionPerSec += proxyStat.MaxSessionPerSec
			backendMetricsReport.TotalCurrentQueuedRequests += proxyStat.CurrentQueued
			backendMetricsReport.TotalBackendConnectionErrors += proxyStat.ErrorConnecting
			backendMetricsReport.AverageQueueTimeMs += proxyStat.AverageQueueTimeMs
			backendMetricsReport.AverageConnectTimeMs += proxyStat.AverageConnectTimeMs
			backendMetricsReport.AverageSessionTimeMs += proxyStat.AverageSessionTimeMs
			populateProxyStats(proxyStat, backendMetricsReport.ProxyMetrics)

		}
	}
	backendMetricsReport.AverageQueueTimeMs = backendMetricsReport.AverageQueueTimeMs / backendCount
	backendMetricsReport.AverageConnectTimeMs = backendMetricsReport.AverageConnectTimeMs / backendCount
	backendMetricsReport.AverageSessionTimeMs = backendMetricsReport.AverageSessionTimeMs / backendCount

	return &HaProxyMetricsReport{
		FrontendMetrics: &frontendMetricsReport,
		BackendMetrics:  &backendMetricsReport,
	}
}

func populateProxyStats(proxyStat haproxy_client.HaproxyStat, proxyStatsMap map[models.RoutingKey]ProxyStats) {
	key, err := proxyKey(proxyStat.ProxyName)
	if err == nil {
		v := ProxyStats{}
		v.ConnectionTime = proxyStat.AverageConnectTimeMs
		v.CurrentSessions = proxyStat.CurrentSessions
		v.MaxSessionPerSec = proxyStat.MaxSessionPerSec
		proxyStatsMap[key] = v
	}
}

// proxyname i.e.  listen_cfg_9001, listen_cfg_9002
func proxyKey(proxy string) (models.RoutingKey, error) {
	routingKey := models.RoutingKey{}

	proxyNameParts := strings.Split(proxy, "_")
	if len(proxyNameParts) != 3 {
		return routingKey, errors.New("not a valid proxy name")
	}

	port, err := strconv.ParseUint(proxyNameParts[2], 10, 16)
	if err != nil {
		return routingKey, err
	}
	routingKey.Port = uint16(port)
	return routingKey, nil
}
