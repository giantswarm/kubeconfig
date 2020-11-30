module github.com/giantswarm/kubeconfig/v3

go 1.15

require (
	github.com/giantswarm/apiextensions/v3 v3.8.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.4
	github.com/google/go-cmp v0.5.3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.18.9
	k8s.io/apimachinery v0.18.9
	k8s.io/client-go v0.18.9
)

// Apply fix for CVE-2020-15114 not yet released in github.com/spf13/viper.
replace github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.3
