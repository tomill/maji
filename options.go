package main

import (
	"errors"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Dirs    []string `long:"dir" description:"Directories to watch. (default: .)"`
	Exclude []string `long:"exclude" short:"x" description:"Directory/File to ignore."`
	Command []string
}

func GetOptions(osArgs []string) (*options, error) {
	args, command := separate(osArgs)

	opt := &options{}
	app := flags.NewParser(opt, flags.HelpFlag)
	app.Name = "maji"
	app.Usage = `[OPTIONS] [<dir>...] -- <command>`

	dirs, err := app.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	opt.Dirs = append(opt.Dirs, dirs...)
	if len(opt.Dirs) == 0 {
		opt.Dirs = []string{"."}
	}

	opt.Command = command
	if len(opt.Command) == 0 {
		err = errors.New("the required argument <command> was not specified")
		return nil, err
	}

	return opt, nil
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
