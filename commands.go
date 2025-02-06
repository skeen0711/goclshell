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

// pwdCommand outputs the current working directory
func pwdCommand(w *io.PipeWriter) {
	defer w.Close()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Fprintln(w, dir)
}

// changeDirectory changes the current directory to the provided path
func changeDirectory(args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path, _ = os.UserHomeDir()
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
	}
}

// `notifyCmds` sends a signal to all running commands
func notifyCmds(cmds []*exec.Cmd, sig os.Signal) {
	for _, cmd := range cmds {
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
	}
}

// handleSignals sets up signal handling for clean termination of commands
func handleSignals(ctx context.Context, cancel context.CancelFunc, cmds *[]*exec.Cmd) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func(ctx context.Context) {
		for {
			select {
			case sig := <-signalChannel:
				switch sig {
				case os.Interrupt, syscall.SIGTERM:
					notifyCmds(*cmds, sig)
				}
			case <-ctx.Done():
				close(signalChannel)
				return
			}
		}
	}(ctx)
}

// readUserInput reads a command line input from the user
func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// executeInput parses and executes a given input string
func executeInput(input string, ctx context.Context, cancel context.CancelFunc, cmds *[]*exec.Cmd, history *History) error {
	commands := strings.Split(input, "|")
	var output io.ReadCloser

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
			go pwdCommand(pw)
			if len(commands) == 1 {
				io.Copy(os.Stdout, pr) // Output directly to stdout
			} else {
				output = pr // Pipe output
			}

		case "history":
			history.ListCommands()

		default:
			err := executeExternalCommand(command, args, cmds, output)
			if err != nil {
				fmt.Printf("Command execution failed: %v\n", err)
			}
		}
	}

	// Handle final output
	if len(*cmds) > 0 {
		(*cmds)[len(*cmds)-1].Stdout = os.Stdout
		for _, cmd := range *cmds {
			err := cmd.Start()
			if err != nil {
				return err
			}
		}
		for _, cmd := range *cmds {
			err := cmd.Wait()
			if err != nil && err.Error() != "signal: interrupt" {
				if cmd.ProcessState.ExitCode() == -1 {
					fmt.Printf("Command not found: %s\n", cmd.Path)
				}
			}
		}
	}

	return nil
}

// executeExternalCommand handles execution of non-built-in commands
func executeExternalCommand(command string, args []string, cmds *[]*exec.Cmd, output io.ReadCloser) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	*cmds = append(*cmds, cmd)

	// Handle piping
	if output != nil {
		cmd.Stdin = output
	}
	var err error
	output, err = cmd.StdoutPipe()
	return err
}
