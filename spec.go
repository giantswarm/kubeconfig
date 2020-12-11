package kubeconfig

import (
	"context"

	"k8s.io/client-go/rest"
)

type Interface interface {
	// NewRESTConfigForApp returns a Kubernetes REST Config for the cluster configured
	// in the secrets objects.
	NewRESTConfigForApp(ctx context.Context, secretName, secretNamespace string) (*rest.Config, error)
}
