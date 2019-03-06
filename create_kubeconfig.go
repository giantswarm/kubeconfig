package kubeconfig

import (
	"github.com/giantswarm/microerror"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubeconfig is a struct used to create a kubectl configuration YAML file
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
	ClientCertificateData string `yaml:"client-certificate"`
	ClientKeyData         string `yaml:"client-key"`
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

func Marshal(config *KubeConfigValue) ([]byte, error) {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return bytes, nil
}

func Unmarshal(bytes []byte) (*KubeConfigValue, error) {
	var kubeConfig KubeConfigValue
	err := yaml.Unmarshal(bytes, &kubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return &kubeConfig, nil
}

func (k *KubeConfig) NewKubeConfigForRESTConfig(config *rest.Config) ([]byte, error) {
	kubeConfig := KubeConfigValue{
		Clusters: []KubeconfigNamedCluster{
			{
				Cluster: KubeconfigCluster{
					Server: config.Host,
				},
			},
		},
	}

	bytes, err := yaml.Marshal(kubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return bytes, nil
}

func (k *KubeConfig) NewRESTConfigForKubeConfig(config *KubeConfigValue) (*rest.Config, error) {
	bytes, err := Marshal(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(bytes)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return restConfig, nil
}
