package kubeconfig

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
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
			errorMatcher: IsNotFoundError,
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
			errorMatcher: IsNotFoundError,
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

func TestKubeConfig_Unmarshal(t *testing.T) {
	testCases := []struct {
		name                 string
		input                []byte
		matchKubeConfigValue KubeConfigValue
		errorMatcher         func(error) bool
	}{
		{
			name: "case 1: unmarshal kubeconfig",
			input: []byte(`
apiVersion: v1
clusters:
- cluster:
    certificate-authority: /workdir/.minikube/ca.crt
    server: https://10.142.5.51:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: /workdir/.minikube/client.crt
    client-key-data: /workdir/.minikube/client.key
`),
			matchKubeConfigValue: KubeConfigValue{
				APIVersion: "v1",
				Kind:       "Config",
				Clusters: []KubeconfigNamedCluster{
					{
						Name: "minikube",
						Cluster: KubeconfigCluster{
							Server: "https://10.142.5.51:8443",
						},
					},
				},
				Users: []KubeconfigUser{
					{
						Name: "minikube",
						User: KubeconfigUserKeyPair{
							ClientCertificateData: "/workdir/.minikube/client.crt",
							ClientKeyData:         "/workdir/.minikube/client.key",
						},
					},
				},
				Contexts: []KubeconfigNamedContext{
					{
						Name: "minikube",
						Context: KubeconfigContext{
							Cluster: "minikube",
							User:    "minikube",
						},
					},
				},
				CurrentContext: "minikube",
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := Unmarshal(tc.input)

			switch {
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case tc.errorMatcher != nil && !tc.errorMatcher(microerror.Cause(err)):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(tc.matchKubeConfigValue, *result) {
				t.Fatalf("want matching kubeconfig value \n %s", cmp.Diff(*result, tc.matchKubeConfigValue))
			}
		})
	}
}

func TestKubeConfig_Marshal(t *testing.T) {
	matchKubeConfigValue := KubeConfigValue{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []KubeconfigNamedCluster{
			{
				Name: "minikube",
				Cluster: KubeconfigCluster{
					Server:                   "https://10.142.5.51:8443",
					CertificateAuthorityData: "Y2FkYXRhdGVzdA==",
				},
			},
		},
		Users: []KubeconfigUser{
			{
				Name: "minikube",
				User: KubeconfigUserKeyPair{
					ClientCertificateData: "Y2NkYXRhdGVzdA==",
					ClientKeyData:         "a2V5ZGF0YXRlc3Q=",
				},
			},
		},
		Contexts: []KubeconfigNamedContext{
			{
				Name: "minikube",
				Context: KubeconfigContext{
					Cluster: "minikube",
					User:    "minikube",
				},
			},
		},
		CurrentContext: "minikube",
	}
	output, err := Marshal(&matchKubeConfigValue)
	if err != nil {
		t.Fatalf("expect nil got %#v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(output)
	if err != nil {
		t.Fatalf("expect nil got %#v", microerror.Mask(err))
	}

	if config.Host != "https://10.142.5.51:8443" {
		t.Fatalf("expect config.Host same as %#v got %#v", "https://10.142.5.51:8443", config.Host)
	}

	if !bytes.Equal(config.CertData, []byte("ccdatatest")) {
		t.Fatalf("expect config.CertData same as %#v got %#v", "ccdatatest", string(config.CertData))
	}

	if !bytes.Equal(config.CAData, []byte("cadatatest")) {
		t.Fatalf("expect config.CAData same as %#v got %#v", "cadatatest", string(config.CAData))
	}

	if !bytes.Equal(config.KeyData, []byte("keydatatest")) {
		t.Fatalf("expect config.CAData same as %#v got %#v", "keydatatest", string(config.CAData))
	}
}
