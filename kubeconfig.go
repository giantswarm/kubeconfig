package kubeconfig

import (
	"context"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewG8sClientFromSecret(ctx context.Context, secretName, secretNamespace string) (versioned.Interface, error){
	restConfig, err := getRESTConfigFromSecret(secretName, secretNamespace)
	if err != nil {
		return nil, err
	}

	client, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "rest.RESTClientFor")
	}
	return client, nil
}

func NewK8sClientFromSecret(ctx context.Context, secretName, secretNamespace string) (kubernetes.Interface, error) {
	restConfig, err := getRESTConfigFromSecret(secretName, secretNamespace)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Maskf(err, "rest.RESTClientFor")
	}
	return client, nil
}

func getRESTConfigFromSecret(secretName, secretNamespace string) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, microerror.Maskf(err, "rest.InClusterConfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
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

