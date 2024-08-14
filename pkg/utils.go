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
		return "", fmt.Errorf("could not find 'man'")
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
		fmt.Fprintln(os.Stdout, GetHelpInfo())
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

func GetHelpInfo() string {
	return `
Virtual File System(1)                              User Commands                              Virtual File System(1)


NAME
       Virtual File System - a CLI file system.

SYNOPSIS
       vfs [-h | --help] [command] [options]

DESCRIPTION
       This is a pure CLI file system written in Go. The system is used to deal with three types of management: User,
       Folder, and File.

COMMANDS
       register [username]
              Register a new user.

       create-folder [username] [foldername] [description]
              Create a folder for the specified user.

       delete-folder [username] [foldername]
              Delete the specified folder for the user.

       list-folders [username] [--sort-name|--sort-created] [asc|desc]
              List all the folders for the user.

       rename-folder [username] [foldername] [new-folder-name]
              Rename the folder.

       create-file [username] [foldername] [filename] [description]?
              Create file from a folder for the user.

       delete-file [username] [foldername] [filename]
              Delete file from a folder for the user.

       list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]
              List all files under the folder for the user.

OPTIONS
       -h, --help
              Show help options.

Virtual File System 1.0                              August 2024                               Virtual File System(1)
`
}
