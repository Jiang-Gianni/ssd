//nolint:all
package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pbar1/pkill-go"
)

var (
	restart    = make(chan interface{})
	testRegexp = regexp.MustCompile("_test.go")
)

//nolint:funlen
func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch entire directory and subdirectories .
	addEntries := func(string) {}
	addEntries = func(name string) {
		entries, err := os.ReadDir(name)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != ".git" && entry.Name() != "cmd" {
				var entryName string
				if name == "./" {
					entryName = "./" + entry.Name()
				} else {
					entryName = name + "/" + entry.Name()
				}

				err = watcher.Add(entryName)
				if err != nil {
					log.Fatal(err)
				}

				log.Println(entryName)
				addEntries(entryName)
			}
		}
	}

	addEntries("./")

	err = watcher.Add("./SCHEMA.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add("./main.go")
	if err != nil {
		log.Fatal(err)
	}

	go RunMainGo()

	log.Println("watching for file changes")

	interval := time.Second
	goTicker := time.NewTicker(interval)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			name := event.Name

			if (name[len(name)-2:] == "go" && !testRegexp.MatchString(name)) || name[len(name)-3:] == "sql" {
				select {
				case <-goTicker.C:
					onGoUpdate(name)
				default:
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}

func onGoUpdate(fileName string) {

	out := regexp.MustCompile("(/)(.*?.go)").ReplaceAllString(fileName, "$1")
	cmd := exec.Command("wsl", "-fix", "./"+out)
	RunCmd(cmd, false)

	restart <- struct{}{}
}

// Kill main if in execution, run main.go and signal to reload chan .
func RunMainGo() {
	for {
		cmd := exec.Command("go", "run", "main.go")
		go RunCmd(cmd, true)
		<-restart
		_, err := pkill.Pkill("main", os.Kill)
		if err != nil {
			log.Println(err)
		}
	}
}

func RunCmd(cmd *exec.Cmd, wantStdout bool) {
	if wantStdout {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
