package pkg

import (
	"fmt"
	"os"
	"regexp"
	"sort"
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

		username, foldername := parts[1], parts[2]
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

		username, foldername := parts[1], parts[2]

		s.DeleteFolder(username, foldername)

	case "list-folders":
		if len(parts) < 2 || len(parts) > 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username := parts[1]
		sortBy, order, msg := ParseArgs(parts)
		if msg != "" {
			fmt.Fprintln(os.Stderr, msg)
			return
		}

		s.ListFolders(username, sortBy, order)

	case "rename-folder":
		if len(parts) != 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}
		username, oldName, newName := parts[1], parts[2], parts[3]

		s.RenameFolder(username, oldName, newName)

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
		fmt.Fprintln(os.Stderr, ErrUnknownCmd.ToString())
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

func (s *System) ListFolders(username, sortBy, order string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(os.Stderr, ErrNotExists.ToString(username))
		return
	}
	if len(user.Folders) == 0 {
		fmt.Fprintln(os.Stderr, WarnNoFolders.ToString(username))
		return
	}

	folders := user.GetFolders()

	switch sortBy {
	case "name":
		sort.Slice(folders, func(i, j int) bool {
			if order == "asc" {
				return folders[i].Name < folders[j].Name
			}
			return folders[i].Name > folders[j].Name
		})

	case "created":
		sort.Slice(folders, func(i, j int) bool {
			if order == "asc" {
				return folders[i].CreatedAt.Before(folders[j].CreatedAt)
			}
			return folders[i].CreatedAt.After(folders[j].CreatedAt)
		})
	}

	fmt.Fprintf(os.Stdout, "Name\t\tDescription\t\tCreatedAt\n")
	for _, folder := range folders {
		fmt.Fprintf(os.Stdout, "%s\t\t%s\t\t\t%s\n",
			folder.Name,
			folder.Description,
			folder.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
}

func (s *System) RenameFolder(username, oldName, newName string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(os.Stderr, ErrNotExists.ToString(username))
		return
	}
	if len(user.Folders) == 0 {
		fmt.Fprintln(os.Stderr, WarnNoFolders.ToString(username))
		return
	}
	if !s.CharsValidator.MatchString(newName) {
		fmt.Fprintln(os.Stderr, ErrInvalidChars.ToString(newName))
		return
	}

	folder := user.GetFolder(oldName)
	user.Folders[newName] = folder
	delete(user.Folders, oldName)

	fmt.Fprintf(os.Stdout, "Rename %s to %s successfully.\n", oldName, newName)
}
