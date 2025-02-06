package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// `pwdCommand` outputs the current working directory
func Pwd(w *io.PipeWriter) {
	defer w.Close()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, dir)
}

// `changeDirectory` changes the current directory based on input arguments
func changeDirectory(args []string) {
	var path string
	var err error

	if len(args) > 0 {
		path = args[0]
	} else {
		path, _ = os.UserHomeDir()
	}
	err = os.Chdir(path)
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
