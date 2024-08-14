package main

import (
	"bufio"
	"fmt"
	"os"
	"system/pkg"
)

func init() {
	path, err := pkg.GetManPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cannot load help info because %v\n", err)
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

	vfs := pkg.SetupSystem()

	for scanner.Scan() {
		input := scanner.Text()
		vfs.Execute(input)
		fmt.Print("$ ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
