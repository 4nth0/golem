package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/4nth0/golem/command"

	log "github.com/sirupsen/logrus"
)

var Version string
var BasePath string = "./.golem"
var TemplatePath string = BasePath + "/templates"
var DatabasePath string = BasePath + "/db"
var ConfigPath string = "./golem.yaml"
var DefaultPort string = "3000"

func configureLog() {
	log.SetOutput(os.Stdout)

	logLevel := "info"

	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = value
	}

	logrusLevel, errLogLevel := log.ParseLevel(logLevel)

	if errLogLevel != nil {
		log.Fatalf("ENV LOG_LEVEL provided is not a viable option, can be either: panic, fatal, error, warn, info, debug, trace")
	}
	log.Printf("Set log level to: %s", logLevel)
	log.SetLevel(logrusLevel)
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	configureLog()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()

	commands := map[string]command.Command{
		"init": command.InitCmd(),
		"help": command.HelpCmd(),
		"run":  command.RunCmd(ctx, ConfigPath),
		"json": command.JsonCmd(ctx),
		"add":  command.AddCmd(),
	}

	fs := flag.NewFlagSet("golem", flag.ExitOnError)
	fs.Parse(os.Args[1:])
	args := fs.Args()

	if len(args) == 0 {
		command.HelpCmd()
		log.Print("No argument provided")
		return
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.Handler(args[1:]); err != nil {
		fmt.Println(err)
		command.HelpCmd()
	}
}
