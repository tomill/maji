package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/radovskyb/watcher"
)

func main() {
	opt := struct {
		Dirs    []string `long:"dir"               description:"Directory/File to watch. (default: .)"`
		Exclude []string `long:"exclude" short:"x" description:"Directory/File to ignore."`
	}{}

	app := flags.NewParser(&opt, flags.HelpFlag)
	app.Name = "maji"
	app.Usage = `[OPTIONS] [<dir>...] -- <command>`

	var command []string
	args := os.Args[1:]
	for i, v := range os.Args {
		if v == "--" {
			args = os.Args[1:i]
			command = os.Args[i+1:]
			break
		}
	}

	bare, err := app.ParseArgs(args)
	if err != nil {
		logFatal(err)
	}

	opt.Dirs = append(opt.Dirs, bare...)
	if len(opt.Dirs) == 0 {
		opt.Dirs = []string{"."}
	}

	w, err := NewWatcher(opt.Dirs, opt.Exclude)
	if err != nil {
		logFatal(err)
	}

	p := NewProcess(command)

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	go func() {
		defer w.Close()
		defer p.Stop()
		for {
			select {
			case sig := <-trap:
				logInfof("%s. quit: %s", sig, p)
				return
			case event := <-w.Event:
				logInfof("%s", event)
				logInfof("restart %s", p)
				p.Stop()
				if err := p.Start(); err != nil {
					logInfof("%s", err)
				}
			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					logInfof("%s", err)
					continue
				}
				return
			case <-w.Closed:
				return
			}
		}
	}()

	logInfof("start %s", p)
	if err = p.Start(); err != nil {
		logInfof("%s", err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		logFatal(err)
	}
}

func logInfof(s string, v ...interface{}) {
	_, _ = color.New(color.FgBlue).Println("[maji] " + fmt.Sprintf(s, v...))
}

func logFatal(err error) {
	_, _ = color.New(color.FgRed).Fprintln(os.Stderr, "[maji] error: "+err.Error())
	os.Exit(1)
}
