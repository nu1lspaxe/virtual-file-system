package pkg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type System struct {
	UserTable      map[string]User
	CharsValidator *regexp.Regexp
}

var (
	VFSystem *System
	once     sync.Once
)

func SetupSystem() *System {
	once.Do(func() {
		VFSystem = &System{
			UserTable:      make(map[string]User),
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
		fmt.Fprintln(os.Stderr, "Error: Not implement yet.")
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
		GetManInfo()
	case "exit":
		fmt.Println("See you.")
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, ErrUnknown.ToString())
	}
}
