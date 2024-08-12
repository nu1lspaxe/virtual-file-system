package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"system/pkg"
)

func init() {
	path, err := pkg.GetManPath()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	pkg.SetManInfo(path, "./vfs.1")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var greetings = `
Welcome to Virtual File System!
Type 'help' to get details and 'exit' to leave.
`

	fmt.Print(greetings)
	fmt.Print("$ ")

	for scanner.Scan() {
		input := scanner.Text()
		handleCommand(input)
		fmt.Print("$ ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	switch command {
	case "register":
	case "create-folder":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "delete-folder":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "list-folders":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "rename-folder":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "create-file":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "delete-file":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "list-files":
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
	case "help":
		pkg.GetManInfo()
	case "exit":
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, "Error: Unknown command.")
	}
}
