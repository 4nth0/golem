package main

import (
	"flag"
	"fmt"

	"github.com/4nth0/golem/internal/command"
)

func helpCmd() command.Command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	return command.Command{fs, func(args []string) error {
		fs.Parse(args)
		help()
		return nil
	}}
}

func help() {
	message := "Golem version " + Version
	message += `

Usage: golem <command> [command flags]

run command:
  -config string
    Config file path (default golem.yaml)


json command:
  -port
    The port used to start server (default 3000)
  -entity
    The entity name.
    Use the -entity flag for each entity (ie. golem json -entity users -entity songs)
  -sync
  	Synchronize in memory store with local file
`

	fmt.Println(message)
}
