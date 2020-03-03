package main

import (
	"flag"
	"log"
	"os"

	"github.com/gol4ng/logger"
	"github.com/gol4ng/logger/formatter"
	"github.com/gol4ng/logger/handler"
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
	l := logger.NewLogger(handler.Stream(os.Stdout, formatter.NewDefaultFormatter(formatter.WithContext(true))))

	commands := map[string]command{
		"init": initCmd(l),
		"help": helpCmd(l),
		"run":  runCmd(l),
		"json": jsonCmd(l),
		"add":  addCmd(l),
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
