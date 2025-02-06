package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	var cmds []*exec.Cmd
	ctx, cancel := context.WithCancel(context.Background())

	// Setup history
	history := NewHistory()

	// Goroutine to listen for signals
	go func(ctx context.Context) {
		select {
		case sig := <-signalChannel:
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				notifyCmds(cmds, sig)
			}
		case <-ctx.Done():
			close(signalChannel)
			return
		}
	}(ctx)

	for {
		// Command prompt
		fmt.Print("ccsh> ")

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Continue if input is empty
		if input == "" {
			continue
		}

		// Add command to history
		history.AddCommand(input)

		// Parse commands (split on pipes)
		commands := strings.Split(input, "|")
		var output io.ReadCloser

		// Loop through each command
		for _, commandLine := range commands {
			commandLine = strings.TrimSpace(commandLine)

			parts := strings.Fields(commandLine)
			if len(parts) == 0 {
				continue
			}

			command := parts[0]
			args := parts[1:]

			switch command {
			case "exit", "q":
				cancel()
				os.Exit(0)

			case "cd":
				changeDirectory(args)

			case "pwd":
				pr, pw := io.Pipe()
				go Pwd(pw)
				if len(commands) == 1 {
					io.Copy(os.Stdout, pr) // Output directly to stdout
				} else {
					output = pr // Pipe output for subsequent commands
				}

			case "history":
				for _, cmd := range history.commands {
					fmt.Println(cmd)
				}

			default:
				// Execute external commands
				cmd := exec.Command(command, args...)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr

				cmds = append(cmds, cmd)

				// Handle piping
				if output != nil {
					cmd.Stdin = output
				}
				output, _ = cmd.StdoutPipe()
			}
		}

		// Ensure the last command outputs to Stdout
		if len(cmds) > 0 {
			cmds[len(cmds)-1].Stdout = os.Stdout
		}

		// Start all commands
		for _, cmd := range cmds {
			cmd.Start()
		}

		// Wait for all commands to finish
		for _, cmd := range cmds {
			err := cmd.Wait()
			if err != nil {
				if err.Error() != "signal: interrupt" {
					if cmd.ProcessState.ExitCode() == -1 {
						fmt.Printf("command not found: %s\n", cmd.Path)
					}
				}
			}
		}
		cmds = nil // Clear command list for the next loop
	}
}
