package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var TYPE = [...]string{"echo", "type", "exit", "cd"}

func main() {
	for {
		_, _ = fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		command = strings.TrimSpace(command)

		programArgs := strings.Split(command, " ")
		program := programArgs[0]
		programArgs = programArgs[1:]
		parseCommand(program, programArgs)
	}
}

func parseCommand(command string, args []string) {
	switch command {
	case "echo":
		echo(args)
	case "cd":
	case "exit":
		exit(args)
	case "type":
		typeCommand(args)
	default:
		fmt.Printf("%s: command not found\n", command)
	}
}

func exit(args []string) {
	if len(args) == 0 || len(args) > 1 {
		os.Exit(0)
	}
	exitCode, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		os.Exit(0)
	}
	os.Exit(int(exitCode))
}

func echo(args []string) {
	fmt.Printf("%s\n", strings.Join(args, " "))
}

func typeCommand(args []string) {
	if len(args) == 0 {
		return
	}
	for _, command := range TYPE {
		if command == args[0] {
			fmt.Printf("%s is a shell builtin\n", command)
			return
		}
	}
	fmt.Printf("%s: not found\n", args[0])
}
