package fs

import (
	"log"
	"os"
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

func (f FSClient) WriteLine(line string) error {
	_, err := f.file.WriteString(line + "\n")
	return err
}

func (f FSClient) Close() {
	f.file.Close()
}
