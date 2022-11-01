package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	cmd    string
	cmdDir string
	watch  string
	filter string
)

func init() {
	flag.StringVar(&cmd, "c", "", "command to run e.g. \"go run .\"")
	flag.StringVar(&cmdDir, "d", "./", "run directory, the directory to run the command in")
	flag.StringVar(&watch, "w", "./", "watch directory, directory to recursively watch")
	flag.StringVar(&filter, "f", "", "comma separated extensions to watch e.g. \"go,cc,c,h\" (default=\"\")")
}

func main() {
	flag.Parse()

	if cmd == "" {
		log.Printf("must provide command with -c\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	split := strings.Split(cmd, " ")
	filters := make([]string, 0)
	if filter != "" {
		filters = strings.Split(filter, ",")
	}

	log.Printf("gowatch: command to run: %v\n", cmd)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("gowatch: failed to init file watcher: %v\n", err)
		flag.PrintDefaults()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	d := make(chan struct{})

	go func(done chan struct{}) {
		var pg ProcessGroup
		var run *exec.Cmd
		var last time.Time
		m := &sync.Mutex{}

		createCmd := func() *exec.Cmd {
			run := exec.Command(split[0])
			run.Args = append(run.Args, split[1:]...)
			run.Dir = cmdDir
			run.Stderr = os.Stderr
			run.Stdin = os.Stdin
			run.Stdout = os.Stdout
			return run
		}

		createPg := func() ProcessGroup {
			pg, err := NewProcessGroup()
			if err != nil {
				log.Fatalf("gowatch: %v\n", err)
			}
			return pg
		}

		rebuild := func() {
			if time.Since(last) < 100*time.Millisecond {
				return
			}
			if !m.TryLock() {
				return
			}
			if err := pg.Kill(); err != nil {
				log.Printf("gowatch: %v\n", err)
			}

			pg = createPg()
			run = createCmd()
			pg.SetPgidToCmd(run)

			last = time.Now()
			go func() {
				if err := run.Start(); err != nil {
					log.Printf("gowatch: %v\n", err)
				} else {
					pg.AddProcess(run.Process)
				}
				m.Unlock()
			}()
		}

		rebuild()

	loop:
		for {
			select {
			case sig := <-sigs:
				log.Printf("gowatch recv: %v\n", sig)
				break loop
			case event, ok := <-watcher.Events:
				if ok {
					if event.Has(fsnotify.Write) {
						for _, ext := range filters {
							if strings.HasSuffix(event.Name, ext) {
								log.Printf("gowatch: file changed: %v\n", event.Name)
								rebuild()
								break
							}
						}
						if len(filters) == 0 {
							log.Printf("gowatch: file changed: %v\n", event.Name)
							rebuild()
						}
					}
				}
			}
		}
		pg.Kill()
		done <- struct{}{}
	}(d)

	err = filepath.Walk(watch, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			log.Printf("gowatch: watching: %v\n", path)
			if err := watcher.Add(path); err != nil {
				log.Printf("gowatch: failed to watch: %v\n", path)
			}
		}
		return err
	})
	if err != nil {
		log.Printf("gowatch: %v\n", err)
	}

	<-d
	log.Printf("gowatch: exiting...\n")
}
