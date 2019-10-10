/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sparkapplication

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/apis/sparkoperator.k8s.io/v1beta2"
	"github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/config"
)

const (
	metricsPropertiesKey       = "metrics.properties"
	prometheusConfigKey        = "prometheus.yaml"
	prometheusScrapeAnnotation = "prometheus.io/scrape"
	prometheusPortAnnotation   = "prometheus.io/port"
	prometheusPathAnnotation   = "prometheus.io/path"
)

func configPrometheusMonitoring(ctx context.Context, app *v1beta2.SparkApplication, client client.Client) error {
	port := config.DefaultPrometheusJavaAgentPort
	if app.Spec.Monitoring.Prometheus.Port != nil {
		port = *app.Spec.Monitoring.Prometheus.Port
	}

	var javaOption string
	if app.HasPrometheusConfigFile() {
		configFile := *app.Spec.Monitoring.Prometheus.ConfigFile
		logger.V(1).Info("Overriding the default Prometheus configuration with config file in the Spark image.", "configFile", configFile)
		javaOption = fmt.Sprintf("-javaagent:%s=%d:%s", app.Spec.Monitoring.Prometheus.JmxExporterJar,
			port, configFile)
	} else {
		logger.V(1).Info("Using the default Prometheus configuration.")
		configMapName := config.GetPrometheusConfigMapName(app)
		configMap := buildPrometheusConfigMap(app, configMapName)
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			namespacedName := types.NamespacedName{
				Namespace: app.Namespace,
				Name:      configMapName,
			}
			cm := &corev1.ConfigMap{}
			err := client.Get(ctx, namespacedName, cm)
			if apiErrors.IsNotFound(err) {
				return client.Create(ctx, configMap)
			}
			if err != nil {
				return err
			}

			cm.Data = configMap.Data
			return client.Update(ctx, cm)
		})

		if retryErr != nil {
			return fmt.Errorf("failed to apply %s in namespace %s: %v", configMapName, app.Namespace, retryErr)
		}

		javaOption = fmt.Sprintf(
			"-javaagent:%s=%d:%s/%s",
			app.Spec.Monitoring.Prometheus.JmxExporterJar,
			port,
			config.PrometheusConfigMapMountPath,
			prometheusConfigKey)
	}

	/* work around for push gateway issue: https://github.com/prometheus/pushgateway/issues/97 */
	metricNamespace := fmt.Sprintf("%s.%s", app.Namespace, app.Name)
	metricConf := fmt.Sprintf("%s/%s", config.PrometheusConfigMapMountPath, metricsPropertiesKey)
	if app.Spec.SparkConf == nil {
		app.Spec.SparkConf = make(map[string]string)
	}
	app.Spec.SparkConf["spark.metrics.namespace"] = metricNamespace
	app.Spec.SparkConf["spark.metrics.conf"] = metricConf

	if app.Spec.Monitoring.ExposeDriverMetrics {
		if app.Spec.Driver.Annotations == nil {
			app.Spec.Driver.Annotations = make(map[string]string)
		}
		app.Spec.Driver.Annotations[prometheusScrapeAnnotation] = "true"
		app.Spec.Driver.Annotations[prometheusPortAnnotation] = fmt.Sprintf("%d", port)
		app.Spec.Driver.Annotations[prometheusPathAnnotation] = "/metrics"

		if app.Spec.Driver.JavaOptions == nil {
			app.Spec.Driver.JavaOptions = &javaOption
		} else {
			*app.Spec.Driver.JavaOptions = *app.Spec.Driver.JavaOptions + " " + javaOption
		}
	}
	if app.Spec.Monitoring.ExposeExecutorMetrics {
		if app.Spec.Executor.Annotations == nil {
			app.Spec.Executor.Annotations = make(map[string]string)
		}
		app.Spec.Executor.Annotations[prometheusScrapeAnnotation] = "true"
		app.Spec.Executor.Annotations[prometheusPortAnnotation] = fmt.Sprintf("%d", port)
		app.Spec.Executor.Annotations[prometheusPathAnnotation] = "/metrics"

		if app.Spec.Executor.JavaOptions == nil {
			app.Spec.Executor.JavaOptions = &javaOption
		} else {
			*app.Spec.Executor.JavaOptions = *app.Spec.Executor.JavaOptions + " " + javaOption
		}
	}

	return nil
}

func buildPrometheusConfigMap(app *v1beta2.SparkApplication, prometheusConfigMapName string) *corev1.ConfigMap {
	metricsProperties := config.DefaultMetricsProperties
	if app.Spec.Monitoring.MetricsProperties != nil {
		metricsProperties = *app.Spec.Monitoring.MetricsProperties
	}
	prometheusConfig := config.DefaultPrometheusConfiguration
	if app.Spec.Monitoring.Prometheus.Configuration != nil {
		prometheusConfig = *app.Spec.Monitoring.Prometheus.Configuration
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            prometheusConfigMapName,
			Namespace:       app.Namespace,
			OwnerReferences: []metav1.OwnerReference{*getOwnerReference(app)},
		},
		Data: map[string]string{
			metricsPropertiesKey: metricsProperties,
			prometheusConfigKey:  prometheusConfig,
		},
	}
}
