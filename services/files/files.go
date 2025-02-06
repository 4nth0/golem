package files

import (
	"context"
	"net/http"

	"github.com/4nth0/golem/log"
)

type FilesServerConfig struct {
	Folder string `yaml:"folder"`
}

func LaunchService(ctx context.Context, port string, config FilesServerConfig) {

	fs := http.FileServer(http.Dir(config.Folder))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
		log.Error("Unable to start server listening.", "err", err)
	}
}
