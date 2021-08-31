package fs

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type FSClient struct {
	dest string
	file *os.File
}

func NewClient(dest string) *FSClient {
	os.Remove(dest)
	f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return &FSClient{
		dest: dest,
		file: f,
	}
}

func (f FSClient) WriteLine(entry string) error {
	line := NewLine(entry)
	_, err := f.file.WriteString(line.String())
	return err
}

func (f FSClient) Close() {
	f.file.Close()
}

type Line struct {
	CreatedAt time.Time `json:"created_at"`
	Entry     string    `json:"entry"`
}

func NewLine(entry string) Line {
	return Line{
		Entry:     entry,
		CreatedAt: time.Now(),
	}
}

func (l Line) String() string {
	b, _ := json.Marshal(l)
	return string(b) + "\n"
}
