package kubeconfig

import (
	"context"

	"k8s.io/client-go/rest"
)

type Interface interface {
	// NewRESTConfigForApp returns a Kubernetes REST Config for the cluster configured
	// in the kubeconfig section of the app CR.
	NewRESTConfigForApp(ctx context.Context, secretName, secretNamespace string) (*rest.Config, error)
}
