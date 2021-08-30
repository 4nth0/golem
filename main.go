package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/4nth0/golem/internal/command"
	"github.com/4nth0/golem/run"

	log "github.com/sirupsen/logrus"
)

var Version string
var BasePath string = "./.golem"
var TemplatePath string = BasePath + "/templates"
var DatabasePath string = BasePath + "/db"
var ConfigPath string = "./golem.yaml"
var DefaultPort string = "3000"

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	commands := map[string]command.Command{
		"init": initCmd(),
		"help": helpCmd(),
		"run":  run.RunCmd(ConfigPath),
		"json": jsonCmd(),
		"add":  addCmd(),
	}

	fs := flag.NewFlagSet("golem", flag.ExitOnError)
	fs.Parse(os.Args[1:])
	args := fs.Args()

	if len(args) == 0 {
		help()
		log.Print("No argument provided")
		return
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.Handler(args[1:]); err != nil {
		fmt.Println(err)
		help()
	}
}
