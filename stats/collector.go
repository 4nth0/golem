package stats

import (
	"context"
	"errors"

	fs "github.com/4nth0/golem/internal/stats/fs"
	pg "github.com/4nth0/golem/internal/stats/postgresql"
	"github.com/4nth0/golem/log"
	"github.com/4nth0/golem/server"
)

type Collector struct {
	writer StatsWriter
	client *StatsClient
}

func NewCollector(driver string, destination string) (*Collector, error) {
	var writer StatsWriter

	switch driver {
	case "fs":
		writer = fs.NewClient(destination)
	case "pg":
		writer = pg.NewClient(destination)
	default:
		return nil, errors.New("UNKNOWN_STATS_DRIVER")
	}

	statsClient := NewClient(writer)

	return &Collector{
		writer: writer,
		client: statsClient,
	}, nil
}

func (c *Collector) Collect(ctx context.Context, requests chan server.InboundRequest) {
	for {
		select {
		case <-ctx.Done():
			c.writer.Close()
			close(requests)
			return
		case request := <-requests:
			err := c.client.PushRequest(request)
			if err != nil {
				log.Error("Unable to collect request", "err", err)
			}
		}
	}
}
