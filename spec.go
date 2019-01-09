package kubeconfig

import (
	"context"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	// NewG8sClient returns a new generated clientset for a tenant cluster.
	NewG8sClient(ctx context.Context, secretName, secretNamespace string) (versioned.Interface, error)
	// NewK8sClient returns a new Kubernetes clientset for a tenant cluster.
	NewK8sClient(ctx context.Context, secretName, secretNamespace string) (kubernetes.Interface, error)
}
