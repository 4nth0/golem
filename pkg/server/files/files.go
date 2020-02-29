package files

import (
	"net/http"

	"github.com/AnthonyCapirchio/golem/pkg/stats"
)

type FilesServerConfig struct {
	Folder string `yaml:"folder"`
}

func LaunchService(ok chan<- bool, stats chan<- stats.StatLine, port string, config FilesServerConfig) {

	fs := http.FileServer(http.Dir(config.Folder))
	http.Handle("/", fs)

	http.ListenAndServe(":"+port, nil)
}
