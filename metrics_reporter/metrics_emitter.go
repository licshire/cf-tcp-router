package metrics_reporter

var (
	totalCurrentQueuedRequests   = PrefixedValue("TotalCurrentQueuedRequests")
	totalBackendConnectionErrors = PrefixedValue("TotalBackendConnectionErrors")
	averageQueueTimeMs           = PrefixedDurationMs("AverageQueueTimeMs")
	averageConnectTimeMs         = PrefixedDurationMs("AverageConnectTimeMs")
	totalBytesIn                 = PrefixedValue("TotalBytesIn")
	totalBytesOut                = PrefixedValue("TotalBytesOut")
	totalMaxSessionPerSec        = PrefixedValue("MaxSessionPerSec")

	connectionTime2   = PrefixedDurationMsByPort("ConnectionTime")
	currentSessions2  = PrefixedValueByPort("CurrentSessions")
	maxSessionsPerSec = PrefixedValueByPort("MaxSessionPerSecond")
)

const (
	backendPrefix  = "backend"
	frontendPrefix = "frontend"
)

//go:generate counterfeiter -o fakes/fake_metrics_emitter.go . MetricsEmitter
type MetricsEmitter interface {
	Emit(*HaProxyMetricsReport)
}

type metricsEmitter struct{}

func NewMetricsEmitter() MetricsEmitter {
	return &metricsEmitter{}
}

func (e *metricsEmitter) Emit(r *HaProxyMetricsReport) {
	if r != nil {
		backendMetrics := r.BackendMetrics
		e.emitCommonMetrics(backendPrefix, backendMetrics)

		frontendMetrics := r.FrontendMetrics
		e.emitCommonMetrics(frontendPrefix, frontendMetrics)

		e.emitBackendMetrics(backendMetrics)

		for k, v := range backendMetrics.ProxyMetrics {
			currentSessions2.Send(backendPrefix, k.String(), v.CurrentSessions)
			maxSessionsPerSec.Send(backendPrefix, k.String(), v.MaxSessionPerSec)

			connectionTime2.Send(backendPrefix, k.String(), v.ConnectionTime)
		}
		for k, v := range frontendMetrics.ProxyMetrics {
			currentSessions2.Send(frontendPrefix, k.String(), v.CurrentSessions)
			maxSessionsPerSec.Send(frontendPrefix, k.String(), v.MaxSessionPerSec)
		}
	}
}

func (e *metricsEmitter) emitCommonMetrics(predix string, r *MetricsReport) {
	totalBytesIn.Send(predix, r.BytesIn)
	totalBytesOut.Send(predix, r.BytesOut)
	totalMaxSessionPerSec.Send(predix, r.MaxSessionPerSec)

}

func (e *metricsEmitter) emitBackendMetrics(r *MetricsReport) {
	totalCurrentQueuedRequests.Send(backendPrefix, r.TotalCurrentQueuedRequests)
	totalBackendConnectionErrors.Send(backendPrefix, r.TotalBackendConnectionErrors)
	averageQueueTimeMs.Send(backendPrefix, r.AverageQueueTimeMs)
	averageConnectTimeMs.Send(backendPrefix, r.AverageConnectTimeMs)
}
