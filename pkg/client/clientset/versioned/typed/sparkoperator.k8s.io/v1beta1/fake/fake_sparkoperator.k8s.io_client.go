// Code generated by k8s code-generator DO NOT EDIT.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/kubeflow/spark-operator/pkg/client/clientset/versioned/typed/sparkoperator.k8s.io/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSparkoperatorV1beta1 struct {
	*testing.Fake
}

func (c *FakeSparkoperatorV1beta1) ScheduledSparkApplications(namespace string) v1beta1.ScheduledSparkApplicationInterface {
	return &FakeScheduledSparkApplications{c, namespace}
}

func (c *FakeSparkoperatorV1beta1) SparkApplications(namespace string) v1beta1.SparkApplicationInterface {
	return &FakeSparkApplications{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSparkoperatorV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
