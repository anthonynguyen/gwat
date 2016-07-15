package main

import (
	"os/exec"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		cmd := exec.Command("go", "run", "sample/sample.go")
		for {
			select {
			case event := <- watcher.Events:
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
					out, err := cmd.Output()
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Print(string(out))
					}
				}
			case err := <- watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./sample/sample.go")
	if err != nil {
		panic(err)
	}
	<- done
}
