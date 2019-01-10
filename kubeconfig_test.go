package kubeconfig

import (
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestKubeConfig_getRESTConfigFromSecret(t *testing.T) {
	testCases := []struct {
		name           string
		presentSecrets []*corev1.Secret
		errorMatcher   func(error) bool
	}{
		{
			name: "case 1: no matching secret",
			presentSecrets: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels:    map[string]string{},
						Name:      "kubeconfig-secret-gs-1",
						Namespace: metav1.NamespaceNone,
					},
					Data: map[string][]byte{
						"test": []byte("test"),
					},
				},
			},
			errorMatcher: errors.IsNotFound,
		},
		{
			name: "case 2: no kubeconfig found",
			presentSecrets: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels:    map[string]string{},
						Name:      "kubeconfig-secret-gs",
						Namespace: metav1.NamespaceNone,
					},
					Data: map[string][]byte{
						"test": []byte("test"),
					},
				},
			},
			errorMatcher: IsMissingKubeConfigError,
		},
		{
			name: "case 3: secret found and no error",
			presentSecrets: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels:    map[string]string{},
						Name:      "kubeconfig-secret-gs",
						Namespace: metav1.NamespaceNone,
					},
					Data: map[string][]byte{
						"kubeConfig": []byte("test"),
					},
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			objs := make([]runtime.Object, 0, len(tc.presentSecrets))
			for _, cc := range tc.presentSecrets {
				objs = append(objs, cc)
			}

			k := KubeConfig{
				logger:    microloggertest.New(),
				k8sClient: fake.NewSimpleClientset(objs...),
			}
			_, err := k.getKubeConfigFromSecret(nil, "kubeconfig-secret-gs", "")

			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case tc.errorMatcher != nil && !tc.errorMatcher(microerror.Cause(err)):
				t.Fatalf("error == %#v, want matching", err)
			}
		})

	}
}
