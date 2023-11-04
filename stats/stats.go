package stats

import (
	"github.com/4nth0/golem/server"
)

type StatsClient struct {
	Writer StatsWriter
}

type StatsWriter interface {
	WriteLine(server.InboundRequest) error
	Close()
}

func NewClient(writer StatsWriter) *StatsClient {
	return &StatsClient{
		Writer: writer,
	}
}

func (s StatsClient) PushRequest(req server.InboundRequest) error {
	// line, _ := json.Marshal(req)

	return s.Writer.WriteLine(req)
}
