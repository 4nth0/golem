package files

import (
	"context"
	"net/http"
)

type FilesServerConfig struct {
	Folder string `yaml:"folder"`
}

func LaunchService(ctx context.Context, port string, config FilesServerConfig) {

	fs := http.FileServer(http.Dir(config.Folder))
	http.Handle("/", fs)

	http.ListenAndServe(":"+port, nil)
}
