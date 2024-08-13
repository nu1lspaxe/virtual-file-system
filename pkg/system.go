package pkg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type System struct {
	UserTable      map[string]*User
	CharsValidator *regexp.Regexp
}

var (
	VFSystem *System
	once     sync.Once
)

func SetupSystem() *System {
	once.Do(func() {
		VFSystem = &System{
			UserTable:      make(map[string]*User),
			CharsValidator: regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
		}
	})
	return VFSystem
}

func (s *System) Execute(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	switch command {
	case "register":
		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username := parts[1]
		s.Register(username)

	case "create-folder":
		if len(parts) < 3 || len(parts) > 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username := parts[1]
		foldername := parts[2]
		desc := ""
		if len(parts) == 4 {
			desc = parts[3]
		}

		s.CreateFolder(username, foldername, desc)

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
		GetManInfo()
	case "exit":
		fmt.Println("See you.")
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, ErrUnknown.ToString())
	}
}
