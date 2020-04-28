package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/4nth0/golem/internal/config"
	httpService "github.com/4nth0/golem/pkg/server/http"
)

type initOpts struct {
	configFile string
}

var successMessage string = `

The configuration has successfully initialized.

Configuration file location: %s
Configured port: %s

`

// This command create a new configuration file
// golem init

func initCmd() command {
	fs := flag.NewFlagSet("golem init", flag.ExitOnError)

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return InitGolem()
	}}
}

func InitGolem() (err error) {

	createFolders()

	port := askForDefaultPort()

	cfg := config.InitConfig(ConfigPath)
	cfg.SetPort(port)

	cfg.Services = []config.Service{
		config.Service{
			Name: "Ping",
			HTTPConfig: httpService.ServerConfig{
				Routes: map[string]httpService.HTTPHandler{
					"/ping": httpService.HTTPHandler{
						Body: "pong!!",
					},
				},
			},
		},
	}

	cfg.Save()

	fmt.Printf(successMessage, ConfigPath, port)

	return nil
}

func askForDefaultPort() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Which port do you want to use (default %s) ? ", DefaultPort)

	port, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	port = strings.TrimSuffix(port, "\n")

	if port == "" {
		port = DefaultPort
	}

	return port
}

func createFolders() {
	if _, err := os.Stat(BasePath); os.IsNotExist(err) {
		os.Mkdir(BasePath, os.ModePerm)
	}

	if _, err := os.Stat(TemplatePath); os.IsNotExist(err) {
		os.Mkdir(TemplatePath, os.ModePerm)
	}

	if _, err := os.Stat(DatabasePath); os.IsNotExist(err) {
		os.Mkdir(DatabasePath, os.ModePerm)
	}
}
