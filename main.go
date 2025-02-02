package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Coding Challenges: Build your own Shell in Go")
	for {
		fmt.Print("goShell> ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" || input == "q" {
			break
		}
		cmd := exec.Command(input)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		_, err := exec.LookPath(input)
		if err != nil {
			fmt.Println("Command not found, try again:")
		}
		cmd.Wait()
	}
}
