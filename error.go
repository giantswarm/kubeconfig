package kubeconfig

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

var missingKubeConfigError = &microerror.Error{
	Kind: "missingKubeConfigError",
}

// IsMissingKubeConfigError asserts missingKubeConfigError.
func IsMissingKubeConfigError(err error) bool {
	return microerror.Cause(err) == missingKubeConfigError
}
