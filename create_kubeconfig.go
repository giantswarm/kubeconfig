package kubeconfig

import (
	"context"
	"encoding/base64"
	"fmt"
	yaml "gopkg.in/yaml.v2"

	"github.com/giantswarm/microerror"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func marshal(config *KubeConfigValue) ([]byte, error) {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return bytes, nil
}

func unmarshal(bytes []byte) (*KubeConfigValue, error) {
	var kubeConfig KubeConfigValue
	err := yaml.Unmarshal(bytes, &kubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return &kubeConfig, nil
}

func (k *KubeConfig) NewKubeConfigForRESTConfig(ctx context.Context, config *rest.Config, clusterName string) ([]byte, error) {
	if clusterName == "" {
		return nil, microerror.Maskf(executionError, "clusterName must not be empty")
	} else if config == nil {
		return nil, microerror.Maskf(executionError, "config must not be empty")
	}

	kubeConfig := KubeConfigValue{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []KubeconfigNamedCluster{
			{
				Name: clusterName,
				Cluster: KubeconfigCluster{
					Server:                   config.Host,
					CertificateAuthorityData: base64.StdEncoding.EncodeToString(config.TLSClientConfig.CAData),
				},
			},
		},
		Contexts: []KubeconfigNamedContext{
			{
				Name: fmt.Sprintf("%s-context", clusterName),
				Context: KubeconfigContext{
					Cluster: clusterName,
					User:    fmt.Sprintf("%s-user", clusterName),
				},
			},
		},
		Users: []KubeconfigUser{
			{
				Name: fmt.Sprintf("%s-user", clusterName),
				User: KubeconfigUserKeyPair{
					ClientCertificateData: base64.StdEncoding.EncodeToString(config.CertData),
					ClientKeyData:         base64.StdEncoding.EncodeToString(config.KeyData),
				},
			},
		},
		CurrentContext: fmt.Sprintf("%s-context", clusterName),
	}

	bytes, err := yaml.Marshal(kubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return bytes, nil
}

func (k *KubeConfig) NewRESTConfigForKubeConfig(ctx context.Context, kubeConfig []byte) (*rest.Config, error) {
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return restConfig, nil
}
