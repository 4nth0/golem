package files

import (
	"net/http"
)

type FilesServerConfig struct {
	Folder string `yaml:"folder"`
}

func LaunchService(port string, config FilesServerConfig) {

	fs := http.FileServer(http.Dir(config.Folder))
	http.Handle("/", fs)

	http.ListenAndServe(":"+port, nil)
}
