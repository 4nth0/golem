package command

import (
	"flag"
	"fmt"
)

func HelpCmd(version string) Command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	return Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			err := fs.Parse(args)
			if err != nil {
				return err
			}
			help(version)
			return nil
		},
	}
}

func help(version string) {
	message := "Golem version " + version
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
