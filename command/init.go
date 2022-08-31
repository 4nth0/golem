package command

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/4nth0/golem/config"
	httpService "github.com/4nth0/golem/services/http"
	log "github.com/sirupsen/logrus"
)

var successMessage string = `

The configuration has successfully initialized.

Configuration file location: %s
Configured port: %s

`

// This command create a new configuration file
// golem init

func InitCmd() Command {
	fs := flag.NewFlagSet("golem init", flag.ExitOnError)

	return Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			err := fs.Parse(args)
			if err != nil {
				return err
			}
			return InitGolem()
		},
	}
}

func InitGolem() (err error) {

	createFolders()

	port := askForDefaultPort()

	cfg := config.InitConfig(ConfigPath)
	cfg.SetPort(port)

	cfg.Services = []config.Service{
		{
			Name: "Ping",
			HTTPConfig: httpService.ServerConfig{
				Routes: map[string]httpService.HTTPHandler{
					"/ping": {
						Body: "pong!!",
					},
					"/multiple-bodies": {
						Bodies: []string{
							"“People don’t just disappear, Dean. Other people just stop looking for them.” — Sam Winchester",
							"“The internet is more than just naked people. You do know that?” — Sam Winchester",
							"“I’ll interrogate the cat.” — Castiel",
							"“If you’re gonna make an omelet, sometimes you have to break some spines.” — Crowley",
						},
					},
					"/user/:id": {
						Body: "Hi User N°${params.id}!",
					},
				},
			},
		},
	}

	err = cfg.Save()
	if err != nil {
		log.WithFields(
			log.Fields{
				"err": err,
			}).Error("Unable to save configuration.")
	}

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
		err = os.Mkdir(BasePath, os.ModePerm)
		if err != nil {
			log.Error("Unable to create base folder", err)
		}
	}

	if _, err := os.Stat(TemplatePath); os.IsNotExist(err) {
		err = os.Mkdir(TemplatePath, os.ModePerm)
		if err != nil {
			log.Error("Unable to create templates folder", err)
		}
	}

	if _, err := os.Stat(DatabasePath); os.IsNotExist(err) {
		err = os.Mkdir(DatabasePath, os.ModePerm)
		if err != nil {
			log.Error("Unable to create db folder", err)
		}
	}
}
