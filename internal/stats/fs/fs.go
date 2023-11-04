package fs

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/4nth0/golem/server"
)

type Client struct {
	dest string
	file *os.File
}

func NewClient(dest string) *Client {
	os.Remove(dest)
	f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return &Client{
		dest: dest,
		file: f,
	}
}

func (f Client) WriteLine(entry server.InboundRequest) error {
	line := NewLine(entry)
	_, err := f.file.WriteString(line.String())
	return err
}

func (f Client) Close() {
	f.file.Close()
}

type Line struct {
	CreatedAt time.Time             `json:"created_at"`
	Entry     server.InboundRequest `json:"entry"`
}

func NewLine(entry server.InboundRequest) Line {
	return Line{
		Entry:     entry,
		CreatedAt: time.Now(),
	}
}

func (l Line) String() string {
	b, _ := json.Marshal(l)
	return string(b) + "\n"
}
