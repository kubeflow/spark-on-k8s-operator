package sparkapplication

import (
	"github.com/kubeflow/spark-operator/pkg/util"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSparkAppMetrics(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	// Test with label containing "-". Expect them to be converted to "_".
	metricsConfig := &util.MetricConfig{
		MetricsPrefix:                 "",
		MetricsLabels:                 []string{"app-id", "namespace"},
		MetricsJobStartLatencyBuckets: []float64{30, 60, 90, 120},
	}
	metrics := newSparkAppMetrics(metricsConfig)
	app1 := map[string]string{"app_id": "test1", "namespace": "default"}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			metrics.sparkAppCount.With(app1).Inc()
			metrics.sparkAppSubmitCount.With(app1).Inc()
			metrics.sparkAppRunningCount.Inc(app1)
			metrics.sparkAppSuccessCount.With(app1).Inc()
			metrics.sparkAppFailureCount.With(app1).Inc()
			metrics.sparkAppFailedSubmissionCount.With(app1).Inc()
			metrics.sparkAppSuccessExecutionTime.With(app1).Observe(float64(100 * i))
			metrics.sparkAppFailureExecutionTime.With(app1).Observe(float64(500 * i))
			metrics.sparkAppStartLatency.With(app1).Observe(float64(10 * i))
			metrics.sparkAppStartLatencyHistogram.With(app1).Observe(float64(10 * i))
			metrics.sparkAppExecutorRunningCount.Inc(app1)
			metrics.sparkAppExecutorSuccessCount.With(app1).Inc()
			metrics.sparkAppExecutorFailureCount.With(app1).Inc()
		}
		for i := 0; i < 5; i++ {
			metrics.sparkAppRunningCount.Dec(app1)
			metrics.sparkAppExecutorRunningCount.Dec(app1)
		}
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppCount, app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppSubmitCount, app1))
	assert.Equal(t, float64(5), metrics.sparkAppRunningCount.Value(app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppSuccessCount, app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppFailureCount, app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppFailedSubmissionCount, app1))
	assert.Equal(t, float64(5), metrics.sparkAppExecutorRunningCount.Value(app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppExecutorFailureCount, app1))
	assert.Equal(t, float64(10), fetchCounterValue(metrics.sparkAppExecutorSuccessCount, app1))
}
