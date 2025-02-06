package main

import (
	"fmt"
	"io"
	"os"
)

// pwd function: Outputs the current working directory
func pwd(w *io.PipeWriter) {
	defer w.Close()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, dir)
}

// changeDirectory function: Changes the current working directory
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
