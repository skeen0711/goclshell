package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	for {
		fmt.Print("ccsh> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		var command = parts[0]
		var args = parts[1:]

		switch command {
		case "exit":
			os.Exit(0)

		case "cd":
			var path string
			var err error

			if len(args) > 0 {
				path = args[0]
			} else {
				path, _ = os.UserHomeDir()
			}
			err = os.Chdir(path)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

		case "pwd":
			dir, _ := os.Getwd()
			fmt.Println(dir)

		default:
			cmdList := strings.Split(input, "|")
			for _, cmd := range cmdList {
				parts := strings.Fields(cmd)
				pCommand := parts[0]
				args = parts[1:]

				pCommand = strings.TrimSpace(pCommand)
				pArgs := strings.TrimSpace(strings.Join(args, " "))
				fCom := exec.Command(pCommand, pArgs)
				fCom.Stdout = os.Stdout
				fCom.Stderr = os.Stderr

				err := fCom.Run()
				if err != nil {
					fmt.Println(err)
				}
			}
			cmd := exec.Command(command, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
