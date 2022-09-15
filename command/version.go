package command

import (
	"flag"
	"fmt"
)

func VersionCmd(version string) Command {
	fs := flag.NewFlagSet("golem json", flag.ExitOnError)

	return Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			err := fs.Parse(args)
			if err != nil {
				return err
			}
			fmt.Println("Golem version " + version)
			return nil
		},
	}
}
