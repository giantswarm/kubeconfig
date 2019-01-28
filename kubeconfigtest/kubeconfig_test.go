package kubeconfigtest

import "testing"

func Test_New(t *testing.T) {
	// Test that New doesn't panic and kubeconfig.Interface is implemented.
	New(Config{})
}
