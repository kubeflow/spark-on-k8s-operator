// Code generated by k8s code-generator DO NOT EDIT.

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	"fmt"

	sparkoperatorv1beta1 "github.com/kubeflow/spark-operator/pkg/client/clientset/versioned/typed/sparkoperator.k8s.io/v1beta1"
	sparkoperatorv1beta2 "github.com/kubeflow/spark-operator/pkg/client/clientset/versioned/typed/sparkoperator.k8s.io/v1beta2"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	SparkoperatorV1beta1() sparkoperatorv1beta1.SparkoperatorV1beta1Interface
	SparkoperatorV1beta2() sparkoperatorv1beta2.SparkoperatorV1beta2Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	sparkoperatorV1beta1 *sparkoperatorv1beta1.SparkoperatorV1beta1Client
	sparkoperatorV1beta2 *sparkoperatorv1beta2.SparkoperatorV1beta2Client
}

// SparkoperatorV1beta1 retrieves the SparkoperatorV1beta1Client
func (c *Clientset) SparkoperatorV1beta1() sparkoperatorv1beta1.SparkoperatorV1beta1Interface {
	return c.sparkoperatorV1beta1
}

// SparkoperatorV1beta2 retrieves the SparkoperatorV1beta2Client
func (c *Clientset) SparkoperatorV1beta2() sparkoperatorv1beta2.SparkoperatorV1beta2Interface {
	return c.sparkoperatorV1beta2
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.sparkoperatorV1beta1, err = sparkoperatorv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.sparkoperatorV1beta2, err = sparkoperatorv1beta2.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.sparkoperatorV1beta1 = sparkoperatorv1beta1.NewForConfigOrDie(c)
	cs.sparkoperatorV1beta2 = sparkoperatorv1beta2.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.sparkoperatorV1beta1 = sparkoperatorv1beta1.New(c)
	cs.sparkoperatorV1beta2 = sparkoperatorv1beta2.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
