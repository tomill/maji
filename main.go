package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
	"github.com/radovskyb/watcher"
)

const Name = "maji"

type option struct {
	Dirs    []string `long:"dir"               description:"Directory/File to watch. (default: .)"`
	Exclude []string `long:"exclude" short:"x" description:"Directory/File to ignore."`
	Command []string
}

func main() {
	opt := &option{}
	app := flags.NewParser(opt, flags.HelpFlag)
	app.Name = Name
	app.Usage = `[OPTIONS] [<dir>...] -- <command>`

	args := os.Args[1:]
	for i, v := range os.Args {
		if v == "--" {
			args = os.Args[1:i]
			opt.Command = os.Args[i+1:]
			break
		}
	}

	bare, err := app.ParseArgs(args)
	if err != nil {
		log.Fatalln(err)
	}

	opt.Dirs = append(opt.Dirs, bare...)
	if len(opt.Dirs) == 0 {
		opt.Dirs = []string{"."}
	}

	if err := run(opt); err != nil {
		log.Fatalln(err)
	}
}

func run(opt *option) error {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.IgnoreHiddenFiles(true)

	listing := func() {
		for _, v := range opt.Dirs {
			if err := w.AddRecursive(v); err != nil {
				infof("ignored %s", err)
			}
		}
		if err := w.Ignore(opt.Exclude...); err != nil {
			infof("ignored %s", err)
		}

		if len(w.WatchedFiles()) == 0 {
			log.Fatalln("error no files to watch")
		}
	}

	p := NewProcess(opt.Command)
	infof("start %s", p)
	if err := p.Start(); err != nil {
		infof("error %s", err)
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	go func() {
		defer w.Close()
		defer p.Stop()
		defer infof("bye~~")
		for {
			select {
			case sig := <-trap:
				infof("%s. quit: %s", sig, p)
				return
			case event := <-w.Event:
				infof("event %s", event)
				infof("restart %s", p)
				p.Stop()
				if err := p.Start(); err != nil {
					infof("error %s", err)
				}
			case err := <-w.Error:
				infof("event %s", err)
				if err == watcher.ErrWatchedFileDeleted {
					time.Sleep(time.Millisecond * 500)
					listing()
					continue
				}
				return
			case <-w.Closed:
				return
			}
		}
	}()

	listing()
	return w.Start(time.Millisecond * 100)
}

func infof(s string, v ...interface{}) {
	_, _ = color.New(color.FgBlue).Printf("["+Name+"] "+s+"\n", v...)
}
