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
		fmt.Println(err)
		os.Exit(1)
	}

	if err := run(opt); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(opt Options) (err error) {
	w, err := NewWatcher(opt.Dirs, opt.Exclude)
	if err != nil {
		return
	} else {
		logInfo("Watching: %s", opt.Dirs)
	}

	p := NewProcess(opt.Command)
	err = p.Start()
	if err != nil {
		logWarn("%s", err)
	} else {
		logInfo("Started: %s", p)
	}

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
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		logInfo("Quitting: %s", p)
		p.Stop()
		w.Close()
		os.Exit(0)
	}()

	err = w.Start(time.Millisecond * 100)
	return
}
