package kubeconfig

import (
	"context"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Config represents the configuration used to create a new kubeconfig library instance.
type Config struct {
	RestConfig rest.Config
	Logger     micrologger.Logger
}

// TenantCluster provides functionality for connecting to tenant clusters.
type KubeConfig struct {
	logger     micrologger.Logger
	restConfig *rest.Config
}

// New creates a new tenant cluster service.
func New(config Config) (*KubeConfig, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if &config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Config must not be empty", config)
	}

	g := &KubeConfig{
		logger:     config.Logger,
		restConfig: rest.CopyConfig(&config.RestConfig),
	}

	return g, nil
}

// NewG8sClientFromSecret returns a generated clientset based on the secret kubeconfig information.
func (k KubeConfig) NewG8sClientFromSecret(ctx context.Context, secretName, secretNamespace string) (versioned.Interface, error) {
	restConfig, err := k.getRESTConfigFromSecret(ctx, secretName, secretNamespace)
	if err != nil {
		return nil, err
	}

	client, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "rest.RESTClientFor")
	}
	return client, nil
}

// NewK8sClientFromSecret returns a Kubernetes clientset based on the secret kubeconfig information.
func (k KubeConfig) NewK8sClientFromSecret(ctx context.Context, secretName, secretNamespace string) (kubernetes.Interface, error) {
	restConfig, err := k.getRESTConfigFromSecret(ctx, secretName, secretNamespace)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "kubernetes.RESTClientFor")
	}
	return client, nil
}

// getRESTConfigFromSecret returns Kubernetes REST config based on the specified secret kubeconfig information.
func (k KubeConfig) getRESTConfigFromSecret(ctx context.Context, secretName, secretNamespace string) (*rest.Config, error) {
	clientset, err := kubernetes.NewForConfig(k.restConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "kubernetes.NewForConfig")
	}

	secret, err := clientset.CoreV1().Secrets(secretNamespace).Get(secretName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, microerror.Maskf(err, "secret namespace: %v, name: %v not found", secretNamespace, secretName)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		return nil, microerror.Maskf(err, "error getting secret %v", statusError.ErrStatus.Message)
	} else if err != nil {
		return nil, microerror.Maskf(err, "unknown error")
	}
	if bytes, ok := secret.Data["kubeConfig"]; ok {
		restConfig, err := clientcmd.RESTConfigFromKubeConfig(bytes)
		if err != nil {
			return nil, microerror.Maskf(err, "clientcmd.RESTConfigFromKubeConfig")
		}
		return restConfig, nil
	} else {
		return nil, microerror.New("secret object do not contain kubeConfig key")
	}
}
