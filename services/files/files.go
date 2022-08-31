package files

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type FilesServerConfig struct {
	Folder string `yaml:"folder"`
}

func LaunchService(ctx context.Context, port string, config FilesServerConfig) {

	fs := http.FileServer(http.Dir(config.Folder))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
		log.WithFields(
			log.Fields{
				"err": err,
			}).Error("Unable to start server listening.")
	}
}
