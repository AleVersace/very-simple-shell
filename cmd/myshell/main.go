package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

		programArgs := splitArgs(command)
		program := programArgs[0]
		if len(programArgs) > 1 {
			programArgs = programArgs[1:]
			for idx, arg := range programArgs {
				programArgs[idx] = strings.Trim(arg, " ")
			}
		}
		parseCommand(program, programArgs)
	}
}

func splitArgs(input string) []string {
	var args []string
	prevIdxSpace := 0
	quoteIdxStart := 0
	inQuotes := false

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\'':
			if !inQuotes {
				inQuotes = true
				quoteIdxStart = i
			} else {
				inQuotes = false
				args = append(args, input[quoteIdxStart+1:i])
			}
		case '\n':
		case ' ':
			if !inQuotes {
				if quoteIdxStart == i-1 {
					continue
				}
				args = append(args, input[prevIdxSpace:i])
				prevIdxSpace = i
			}
		}
	}

	if len(input) > 0 && input[len(input)-1] != '\'' {
		args = append(args, input[prevIdxSpace:])
	}

	return args
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
	var pathDelimiter string
	if runtime.GOOS == "windows" {
		pathDelimiter = ";"
	} else {
		pathDelimiter = ":"
	}

	pathsEnv := os.Getenv("PATH")
	paths := strings.Split(pathsEnv, pathDelimiter)
	for _, path := range paths {
		dir, err := os.Open(path)
		if err != nil {
			continue
		}
		defer dir.Close()
		files, err := dir.Readdir(-1)
		if err != nil {
			continue
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
	commandNotFound(command)
}

func execCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(string(stdout))
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
	found := typeCommandInPath(args[0], args[1:])
	if !found {
		fmt.Printf("%s: not found\n", args[0])
	}
}

func typeCommandInPath(command string, args []string) bool {
	var pathDelimiter string
	if runtime.GOOS == "windows" {
		pathDelimiter = ";"
	} else {
		pathDelimiter = ":"
	}

	pathsEnv := os.Getenv("PATH")
	paths := strings.Split(pathsEnv, pathDelimiter)
	for _, path := range paths {
		dir, err := os.Open(path)
		if err != nil {
			continue
		}
		defer dir.Close()
		files, err := dir.Readdir(-1)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if file.Name() == command {
				fmt.Printf("%s is %s/%s\n", command, path, command)
				return true
			}
		}
	}
	return false
}
