package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	command = strings.TrimSpace(command)

	parseCommand(command)
}

func parseCommand(command string) {
	switch command {
	case "echo":
	case "cd":
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}
