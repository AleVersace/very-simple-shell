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

var TYPE = [...]string{"echo", "type", "exit", "cd", "cat"}

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
		if len(programArgs) == 0 {
			continue
		}
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
				if i+1 < len(input) && input[i+1] == '\'' {
					// go ahead if next char is single quote
					i++
					continue
				}
				inQuotes = false
				newArg := strings.TrimSpace(strings.ReplaceAll(input[quoteIdxStart+1:i], "'", ""))
				args = append(args, newArg)
			}
		case '\n':
		case ' ':
			if !inQuotes {
				if input[i-1] == '\'' { // if previous char was single quote new arg was already registered
					continue
				}
				newArg := strings.TrimSpace(input[prevIdxSpace:i])
				if newArg != "" {
					args = append(args, newArg)
				}
				prevIdxSpace = i
			}
		}
	}

	if len(input) > 0 && input[len(input)-1] != '\'' {
		args = append(args, strings.TrimSpace(input[prevIdxSpace:]))
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
	case "cat":
		cat(args)
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

func cat(args []string) {
	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Print("opening file: %w", err)
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s", line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Print("scanning file: %w", err)
			continue
		}
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
