package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	ValidChars = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	manFile    string
)

// GetManPath get the path of `man` command.
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

// SetManInfo is used to set man path and then run man [txtPath] to show deatils.
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
		fmt.Fprintf(os.Stderr, "Error running man: %v\n", err)
		os.Exit(1)
	}
}
