package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/4nth0/golem/internal/command"
)

func addCmd() command.Command {
	fs := flag.NewFlagSet("golem add", flag.ExitOnError)

	return command.Command{
		FlagSet: fs,
		Handler: func(args []string) error {
			fs.Parse(args)
			return Add(args)
		},
	}
}

func Add(args []string) (err error) {

	if len(args) < 2 {
		return errors.New("NOT ENOUGH ARGMUENTS")
	}

	switch args[1] {
	case "http":
		fmt.Println("Create new http service ...")
	default:
		return errors.New("UNKONWN COMMAND")
	}

	return nil
}
