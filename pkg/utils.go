package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	manFile string
)

// GetManPath to get the path of `man` command.
func GetManPath() (string, error) {
	cmd := exec.Command("which", "man")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not find 'man': %v", err)
	}
	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", fmt.Errorf("'man' command not found in PATH")
	}
	return path, nil
}

// SetManInfo to set man path and then run man [txtPath] to show deatils.
func SetManInfo(manPath, txtPath string) {
	os.Setenv("PATH", fmt.Sprintf("%v:", manPath)+os.Getenv("PATH"))
	manFile = txtPath

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		GetManInfo()
		os.Exit(0)
	}
}

func GetManInfo() {
	cmd := exec.Command("man", manFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func ParseArgs(args []string) (sortBy, order, msg string) {
	sortBy = "name"
	order = "asc"

	if len(args) == 0 {
		return sortBy, order, ""
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--sort-name":
			sortBy = "name"
		case "--sort-created":
			sortBy = "created"
		case "asc":
			order = "asc"
		case "desc":
			order = "desc"
		default:
			return "", "", ErrInvalidFlag.ToString()
		}
	}
	return sortBy, order, ""
}
