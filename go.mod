module github.com/giantswarm/kubeconfig/v4

go 1.15

require (
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/google/go-cmp v0.5.6
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.13
	k8s.io/apimachinery v0.20.13
	k8s.io/client-go v0.20.13
)

replace (
	// Apply fix for CVE-2020-15114 not yet released in github.com/spf13/viper.
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.4
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	// Use v1.4.2 of gorilla/websocket to fix nancy alert.
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
)
