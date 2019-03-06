package kubeconfig

import "github.com/spf13/afero"

func mockFiles() {
	appFS := afero.NewOsFs()
	appFS.MkdirAll("src/", 0755)
	afero.WriteFile(appFS, "src/client.crt", []byte("file b"), 0644)
	afero.WriteFile(appFS, "src/client.key", []byte("file c"), 0644)
}
