package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetManPath get the path of `man` command.
func GetManPath() (string, error) {
	cmd := exec.Command("which", "main")
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

func SetManInfo(manPath, txtPath string) {
	os.Setenv("PATH", fmt.Sprintf("%v:", manPath)+os.Getenv("PATH"))

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		cmd := exec.Command("man", txtPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running man: %v\n", err)
			os.Exit(1)
		}
		return
	}
}
