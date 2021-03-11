package kubeconfigtest

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/giantswarm/kubeconfig/v4"
)

type Config struct {
	RestConfig             rest.Config
	RestConfigFromAppError error
}

type KubeConfig struct {
	restConfig             rest.Config
	restConfigFromAppError error
}

func New(config Config) kubeconfig.Interface {
	k := &KubeConfig{
		restConfig:             config.RestConfig,
		restConfigFromAppError: config.RestConfigFromAppError,
	}

	return k
}

func (k *KubeConfig) NewRESTConfigForApp(ctx context.Context, secretName, secretNamespace, secretKey string) (*rest.Config, error) {
	if k.restConfigFromAppError != nil {
		return nil, k.restConfigFromAppError
	}

	return &k.restConfig, nil
}
