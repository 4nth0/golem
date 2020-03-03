package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/gol4ng/logger"
)

type addOpts struct{}

func addCmd(log *logger.Logger) command {
	fs := flag.NewFlagSet("golem add", flag.ExitOnError)

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return add(log, args)
	}}
}

func add(log *logger.Logger, args []string) (err error) {

	if len(args) < 2 {
		return errors.New("Not enough arguments")
	}

	switch args[1] {
	case "http":
		fmt.Println("Create new http service ...")
	default:
		return errors.New("Unknown command")
	}

	return nil
}
