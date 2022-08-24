package command

import "flag"

type Command struct {
	FlagSet *flag.FlagSet
	Handler func(args []string) error
}
