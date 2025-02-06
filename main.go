package main

import (
	"bufio"
	"fmt"
	"io"
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

		commands := strings.Split(input, "|")
		var cmds []*exec.Cmd
		var output io.ReadCloser

		for _, command := range commands {
			command = strings.TrimSpace(command)
			parts := strings.Fields(command)
			var command = parts[0]
			var args = parts[1:]

			switch command {
			case "exit", "q":
				os.Exit(0)

			case "cd":
				changeDirectory(args)

			case "pwd":
				pr, pw := io.Pipe()
				go pwd(pw)
				if len(commands) == 1 {
					io.Copy(os.Stdout, pr)
				} else {
					output = pr
				}

			default:
				// Create system command
				cmd := exec.Command(command, args...)
				cmd.Stderr = os.Stderr
				cmds = append(cmds, cmd)

				if output != nil {
					cmd.Stdin = output
				}
				output, _ = cmd.StdoutPipe()
			}
		}

		if len(cmds) > 0 {
			cmds[len(cmds)-1].Stdout = os.Stdout
		}

		for _, cmd := range cmds {
			cmd.Start()
		}

		for _, cmd := range cmds {
			if err := cmd.Wait(); err != nil {
				if cmd.ProcessState.ExitCode() == -1 {
					fmt.Printf("command not found: %s\n", cmd.Path)
				}
			}
		}
	}
}
