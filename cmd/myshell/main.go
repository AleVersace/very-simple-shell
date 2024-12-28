package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		execCommandInPath(command, args)
	}
}

func execCommandInPath(command string, args []string) {
	pathsEnv := os.Getenv("PATH")
	paths := strings.Split(pathsEnv, ":")
	for _, path := range paths {
		dir, err := os.Open(path)
		if err != nil {
			commandNotFound(command)
		}
		defer dir.Close()
		files, err := dir.Readdir(-1)
		if err != nil {
			commandNotFound(command)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == command {
				execCommand(command, args)
				return
			}
		}
	}
}

func execCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(stdout))
}

func commandNotFound(command string) {
	fmt.Printf("%s: command not found\n", command)
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

	// check if built-in
	for _, command := range TYPE {
		if command == args[0] {
			fmt.Printf("%s is a shell builtin\n", command)
			return
		}
	}

	// check if in PATH
	typeCommandInPath(args[0], args[1:])

	fmt.Printf("%s: not found\n", args[0])
}

func typeCommandInPath(command string, args []string) {
	pathsEnv := os.Getenv("PATH")
	paths := strings.Split(pathsEnv, ":")
	for _, path := range paths {
		dir, err := os.Open(path)
		if err != nil {
			commandNotFound(command)
		}
		defer dir.Close()
		files, err := dir.Readdir(-1)
		if err != nil {
			commandNotFound(command)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == command {
				fmt.Printf("%s is %s\n", command, path)
				return
			}
		}
	}
}
