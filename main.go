package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var data AppData

func init() {
	if _, err := os.Stat("data/db.json"); os.IsNotExist(err) {
		data = AppData{
			SearchHistory: []Search{},
		}
		jsonBytes, err := json.Marshal(&data)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Mkdir("data", 0755)

		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := os.Create("data/db.json")

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, err = f.Write(jsonBytes)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		content, err := ioutil.ReadFile("data/db.json")

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func main() {

	inputDirectory := os.Args[1]
	inputCommand := os.Args[2:]

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer saveData()

	log.Println("Searching...")

	err = filepath.Walk(currentWorkingDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
			}

			if info.Name() == inputDirectory {

				log.Println("Directory found at '%s'. Enter 'Y' to use this directory or any other key to continue searching.", path)

				var userInput string
				fmt.Scanln(&userInput)
				if !strings.EqualFold(userInput, "Y") {
					fmt.Println("Searching...")
					return nil
				} else {
					data.SearchHistory = append(data.SearchHistory, Search{
						UserInput: inputDirectory,
						Result: path,
					})

				}




				command := strings.Join(inputCommand, " ")

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

func saveData() {
	f, err := os.OpenFile("data/db.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        log.Fatal(err)
    }

	jsonBytes, err := json.Marshal(&data)
	
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(jsonBytes)

	if err != nil {
		log.Fatal(err)
	}

    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}

type AppData struct {
	SearchHistory []Search `json:"searchHistory"`
}

type Search struct {
	UserInput string `json:"userInput"`
	Result    string `json:"result"`
}
