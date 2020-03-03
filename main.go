package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}

var Version string
var BasePath string = "./.golem"
var TemplatePath string = BasePath + "/templates"
var DatabasePath string = BasePath + "/db"
var ConfigPath string = "./golem.yaml"
var DefaultPort string = "3000"

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	commands := map[string]command{
		"init": initCmd(),
		"help": helpCmd(),
		"run":  runCmd(),
		"json": jsonCmd(),
		"add":  addCmd(),
	}

	fs := flag.NewFlagSet("golem", flag.ExitOnError)
	fs.Parse(os.Args[1:])
	args := fs.Args()

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		help()
		log.Print(err)
	}
}
