package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/radovskyb/watcher"
)

var logInfo = color.Cyan
var logWarn = color.Red

func main() {
	opt, err := GetOptions(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(opt Options) (err error) {
	w, err := NewWatcher(opt.Dirs, opt.Exclude)
	if err != nil {
		return
	}
	logInfo("Watching: %s", opt.Dirs)

	p := NewProcess(opt.Command)
	err = p.Start()
	if err != nil {
		logWarn("%s", err)
	}
	logInfo("Started: %s", p)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logInfo("Event: %s", event)
				logInfo("Restarting: %s", p)
				p.Stop()
				p.Start()
			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					logWarn("%s", err)
					continue
				}
				return
			case <-w.Closed:
				return
			case <-quit:
				logInfo("Quitting: %s", p)
				p.Stop()
				w.Close()
				return
			}
		}
	}()

	err = w.Start(time.Millisecond * 100)
	return
}
