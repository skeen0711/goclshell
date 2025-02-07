package main

import (
	"context"
	"fmt"
	"os/exec"
)

func main() {
	// Setup signal handling for graceful termination
	ctx, cancel := context.WithCancel(context.Background())
	cmds := make([]*exec.Cmd, 0)
	handleSignals(ctx, cancel, &cmds)

	// Create a history
	history := NewHistory()

	// Command execution loop
	for {
		// Display prompt
		fmt.Print("ccsh> ")

		// Read and trim input
		input := readUserInput()
		if input == "" {
			continue
		}

		// Add to history
		history.AddCommand(input)

		// Parse and execute commands
		err := executeInput(input, ctx, cancel, &cmds, history)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		// Clear commands list for next iteration
		cmds = nil
	}
}
