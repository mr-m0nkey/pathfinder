package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	inputDirectory := os.Args[1]
	inputCommand := os.Args[2:]

	
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = filepath.Walk(currentWorkingDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Name() == inputDirectory {
		
				command := strings.Join(inputCommand, " ")
				//TODO add support for other operating systems asides windows
				cmd := exec.Command("cmd", "/C", command)
				cmd.Dir = path
				cmd.Stderr = os.Stdout
				commandOutput, err := cmd.Output()
				if err != nil {
					return err
				}
				log.Println(fmt.Sprintf("Running '%s' on path: '%s'", command, path))
				log.Println(string(commandOutput))
				return io.EOF
			}

			return nil
		})

	if err != nil && err != io.EOF {
		log.Println(err)
	}

}
