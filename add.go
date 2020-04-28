package main

import (
	"errors"
	"flag"
	"fmt"
)

type addOpts struct{}

func addCmd() command {
	fs := flag.NewFlagSet("golem add", flag.ExitOnError)

	return command{fs, func(args []string) error {
		fs.Parse(args)
		return Add(args)
	}}
}

func Add(args []string) (err error) {

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
