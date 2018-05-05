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

func logInfo(s string, v ...interface{}) {
	color.New(color.FgBlue).Println("[maji] " + fmt.Sprintf(s, v...))
}

func logFatal(err error) {
	color.New(color.FgRed).Fprintln(os.Stderr, "[maji] error: "+err.Error())
	os.Exit(1)
}

func main() {
	opt, err := getOptions(os.Args)
	if err != nil {
		logFatal(err)
	}

	if err := run(opt); err != nil {
		logFatal(err)
	}
}

func run(opt *options) error {
	w, err := NewWatcher(opt.Dirs, opt.Exclude)
	if err != nil {
		return err
	}

	logInfo("watching %s", opt.Dirs)

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	p := NewProcess(opt.Command)

	go func() {
		defer w.Close()
		defer p.Stop()
		for {
			select {
			case sig := <-trap:
				logInfo("%s. quit: %s", sig, p)
				return
			case event := <-w.Event:
				logInfo("%s", event)
				logInfo("restart %s", p)
				p.Stop()
				if err := p.Start(); err != nil {
					logInfo("%s", err)
				}
			case err := <-w.Error:
				if err == watcher.ErrWatchedFileDeleted {
					logInfo("%s", err)
					continue
				}
				return
			case <-w.Closed:
				return
			}
		}
	}()

	logInfo("start %s", p)
	if err = p.Start(); err != nil {
		logInfo("%s", err)
	}

	return w.Start(time.Millisecond * 100)
}
