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
		if len(parts) != 3 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username := parts[1]
		foldername := parts[2]

		s.DeleteFolder(username, foldername)

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

func (s *System) Register(username string) {
	if !s.CharsValidator.MatchString(username) {
		fmt.Fprintln(os.Stderr, ErrInvalidChars.ToString(username))
		return
	}
	if user := s.GetUser(username); user != nil {
		fmt.Fprintln(os.Stderr, ErrAlreadyExists.ToString(username))
		return
	}

	s.UserTable[username] = CreateUser(username)
	fmt.Fprintf(os.Stdout, "Add %s successfully.\n", username)
}

func (s *System) GetUser(username string) *User {
	for name := range s.UserTable {
		if name == username {
			return s.UserTable[username]
		}
	}
	return nil
}

func (s *System) CreateFolder(username, foldername, desc string) {

	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(os.Stderr, ErrNotExists.ToString(username))
		return
	}
	if !s.CharsValidator.MatchString(foldername) {
		fmt.Fprintln(os.Stderr, ErrInvalidChars.ToString(foldername))
		return
	}
	if folder := user.GetFolder(foldername); folder != nil {
		fmt.Fprintln(os.Stderr, ErrAlreadyExists.ToString(foldername))
		return
	}

	folder := CreateFolder(foldername, desc)
	user.AddFolder(foldername, folder)

	fmt.Fprintf(os.Stdout, "Create %s successfully.\n", foldername)
}

func (s *System) DeleteFolder(username, foldername string) {

	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(os.Stderr, ErrNotExists.ToString(username))
		return
	}
	if folder := user.GetFolder(foldername); folder == nil {
		fmt.Fprintln(os.Stderr, ErrNotExists.ToString(foldername))
		return
	}

	fmt.Fprintf(os.Stdout, "Delete %v successfully.", foldername)
}

func (s *System) ListFolders(username string) {

}
