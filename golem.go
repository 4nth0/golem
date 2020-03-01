package main

import (
	"flag"
	"log"
	"os"
)

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}

var Version string

func main() {
	commands := map[string]command{
		"run":  runCmd(),
		"json": jsonCmd(),
		"help": helpCmd(),
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
