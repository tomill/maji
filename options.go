package main

import (
	"errors"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Dirs    []string `long:"dir" description:"Diretories to watch. (default: .)"`
	Exclude []string `long:"exclude" short:"x" description:"Diretories to ignore."`
	Command []string
}

func GetOptions(osArgs []string) (opt Options, err error) {
	args, command := separate(osArgs)

	parser := flags.NewParser(&opt, flags.HelpFlag)
	parser.Name = "maji"
	parser.Usage = `[OPTIONS] [<dir>...] -- <command>`

	dirs, err := parser.ParseArgs(args)
	if err != nil {
		return
	}

	opt.Dirs = append(opt.Dirs, dirs...)
	if len(opt.Dirs) == 0 {
		opt.Dirs = []string{"."}
	}

	opt.Command = command
	if len(opt.Command) == 0 {
		err = errors.New("the required argument <command> was not specified")
		return
	}

	return
}

func separate(osArgs []string) (args []string, command []string) {
	args = osArgs

	for i, v := range osArgs {
		if v == "--" {
			args = osArgs[1:i]     // skip $0
			command = osArgs[i+1:] // skip "--"
			break
		}
	}

	return
}
