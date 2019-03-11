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

// KubeconfigValue is a struct used to create a kubectl configuration YAML file
type KubeConfigValue struct {
	APIVersion     string                   `yaml:"apiVersion"`
	Kind           string                   `yaml:"kind"`
	Clusters       []KubeconfigNamedCluster `yaml:"clusters"`
	Users          []KubeconfigUser         `yaml:"users"`
	Contexts       []KubeconfigNamedContext `yaml:"contexts"`
	CurrentContext string                   `yaml:"current-context"`
	Preferences    struct{}                 `yaml:"preferences"`
}

// KubeconfigUser is a struct used to create a kubectl configuration YAML file
type KubeconfigUser struct {
	Name string                `yaml:"name"`
	User KubeconfigUserKeyPair `yaml:"user"`
}

// KubeconfigUserKeyPair is a struct used to create a kubectl configuration YAML file
type KubeconfigUserKeyPair struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

// KubeconfigNamedCluster is a struct used to create a kubectl configuration YAML file
type KubeconfigNamedCluster struct {
	Name    string            `yaml:"name"`
	Cluster KubeconfigCluster `yaml:"cluster"`
}

// KubeconfigCluster is a struct used to create a kubectl configuration YAML file
type KubeconfigCluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

// KubeconfigNamedContext is a struct used to create a kubectl configuration YAML file
type KubeconfigNamedContext struct {
	Name    string            `yaml:"name"`
	Context KubeconfigContext `yaml:"context"`
}

// KubeconfigContext is a struct used to create a kubectl configuration YAML file
type KubeconfigContext struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

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

func (k *KubeConfig) NewRESTConfigForKubeConfig(ctx context.Context, config *KubeConfigValue) (*rest.Config, error) {
	bytes, err := marshal(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(bytes)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return restConfig, nil
}
