package pkg

import (
	"fmt"
	"io"
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

// SetupSystem create a singleton instance (system)
func SetupSystem() *System {
	once.Do(func() {
		VFSystem = &System{
			UserTable:      make(map[string]*User, 0),
			CharsValidator: regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
		}
	})
	return VFSystem
}

func (s *System) Reset() {
	s = nil
	once = sync.Once{}
}

// Execute to call APIs by command
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
		s.Register(os.Stdout, os.Stderr, username)

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

		s.CreateFolder(os.Stdout, os.Stderr, username, foldername, desc)

	case "delete-folder":
		if len(parts) != 3 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username, foldername := parts[1], parts[2]

		s.DeleteFolder(os.Stdout, os.Stderr, username, foldername)

	case "list-folders":
		if len(parts) < 2 || len(parts) > 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username := parts[1]
		sortBy, order, msg := ParseArgs(parts[2:])
		if msg != "" {
			fmt.Fprintln(os.Stderr, msg)
			return
		}

		s.ListFolders(os.Stdout, os.Stderr, username, sortBy, order)

	case "rename-folder":
		if len(parts) != 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}
		username, oldName, newName := parts[1], parts[2], parts[3]

		s.RenameFolder(os.Stdout, os.Stderr, username, oldName, newName)

	case "create-file":
		if len(parts) < 4 || len(parts) > 5 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}
		username, foldername, filename := parts[1], parts[2], parts[3]
		desc := ""
		if len(parts) == 5 {
			desc = parts[4]
		}

		s.CreateFile(os.Stdout, os.Stderr, username, foldername, filename, desc)

	case "delete-file":
		if len(parts) != 4 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}
		username, foldername, filename := parts[1], parts[2], parts[3]

		s.DeleteFile(os.Stdout, os.Stderr, username, foldername, filename)

	case "list-files":
		if len(parts) < 3 || len(parts) > 5 {
			fmt.Fprintln(os.Stderr, ErrArgsLength.ToString())
			return
		}

		username, foldername := parts[1], parts[2]
		sortBy, order, msg := ParseArgs(parts[3:])
		if msg != "" {
			fmt.Fprintln(os.Stderr, msg)
			return
		}

		s.ListFiles(os.Stdout, os.Stderr, username, foldername, sortBy, order)

	case "help":
		GetManInfo()

	case "exit":
		fmt.Println("See you.")
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, ErrUnknownCmd.ToString())
	}
}

// Register a new user
func (s *System) Register(w io.Writer, ew io.Writer, username string) {
	if !s.CharsValidator.MatchString(username) {
		fmt.Fprintln(ew, ErrInvalidChars.ToString(username))
		return
	}
	if user := s.GetUser(username); user != nil {
		fmt.Fprintln(ew, ErrAlreadyExists.ToString(username))
		return
	}

	s.UserTable[username] = CreateUser(username)
	fmt.Fprintf(w, "Add %s successfully.\n", username)
}

// GetUser to find and return user if exists
func (s *System) GetUser(username string) *User {
	for name := range s.UserTable {
		if name == username {
			return s.UserTable[username]
		}
	}
	return nil
}

// CreateFolder to create a folder for a user, description is optional
func (s *System) CreateFolder(w io.Writer, ew io.Writer, username, foldername, desc string) {

	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	if !s.CharsValidator.MatchString(foldername) {
		fmt.Fprintln(ew, ErrInvalidChars.ToString(foldername))
		return
	}
	if folder := user.GetFolder(foldername); folder != nil {
		fmt.Fprintln(ew, ErrAlreadyExists.ToString(foldername))
		return
	}

	folder := CreateFolder(foldername, desc, username)
	user.AddFolder(foldername, folder)

	fmt.Fprintf(w, "Create %s successfully.\n", foldername)
}

// DeleteFolder to delete a folder from a user if exists
func (s *System) DeleteFolder(w io.Writer, ew io.Writer, username, foldername string) {

	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	if folder := user.GetFolder(foldername); folder == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(foldername))
		return
	}

	delete(user.Folders, foldername)

	fmt.Fprintf(w, "Delete %v successfully.\n", foldername)
}

// ListFolders to list all the folders of a user if exist
func (s *System) ListFolders(w io.Writer, ew io.Writer, username, sortBy, order string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	if len(user.Folders) == 0 {
		fmt.Fprintln(ew, WarnNoFolders.ToString(username))
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

	for _, folder := range folders {
		fmt.Fprintln(w, folder.ToString())
	}
}

// RenameFolder to rename a folder of a user
func (s *System) RenameFolder(w io.Writer, ew io.Writer, username, oldName, newName string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	folder := user.GetFolder(oldName)
	if folder == nil {
		fmt.Fprintln(ew, WarnNoFolders.ToString(oldName))
		return
	}
	if !s.CharsValidator.MatchString(newName) {
		fmt.Fprintln(ew, ErrInvalidChars.ToString(newName))
		return
	}

	folder.SetName(newName)
	user.Folders[newName] = folder
	delete(user.Folders, oldName)

	fmt.Fprintf(w, "Rename %s to %s successfully.\n", oldName, newName)
}

// CreateFile to create a file under a folder of a user
func (s *System) CreateFile(w io.Writer, ew io.Writer, username, foldername, filename, desc string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	folder := user.GetFolder(foldername)
	if folder == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(foldername))
		return
	}
	if !s.CharsValidator.MatchString(filename) {
		fmt.Fprintln(ew, ErrInvalidChars.ToString(filename))
		return
	}
	file := folder.GetFile(filename)
	if file != nil {
		fmt.Fprintln(ew, ErrAlreadyExists.ToString(filename))
		return
	}

	folder.AddFile(filename, CreateFile(
		filename, desc, foldername, username,
	))

	fmt.Fprintf(w, "Create %s in %s/%s successfully.\n", filename, username, foldername)
}

// DeleteFile to delete file under a folder from a user if exist
func (s *System) DeleteFile(w io.Writer, ew io.Writer, username, foldername, filename string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	folder := user.GetFolder(foldername)
	if folder == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(foldername))
		return
	}
	file := folder.GetFile(filename)
	if file == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(filename))
		return
	}

	delete(folder.Files, filename)

	fmt.Fprintf(w, "Delete %s in %s/%s successfully.\n", filename, username, foldername)
}

// ListFiles to list all files from a folder of a user
func (s *System) ListFiles(w io.Writer, ew io.Writer, username, foldername, sortBy, order string) {
	user := s.GetUser(username)
	if user == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(username))
		return
	}
	folder := user.GetFolder(foldername)
	if folder == nil {
		fmt.Fprintln(ew, ErrNotExists.ToString(foldername))
		return
	}
	if len(folder.Files) == 0 {
		fmt.Fprintln(ew, WarnEmptyFolder.ToString())
		return
	}

	files := folder.GetFiles()

	switch sortBy {
	case "name":
		sort.Slice(files, func(i, j int) bool {
			if order == "asc" {
				return files[i].Name < files[j].Name
			}
			return files[i].Name > files[j].Name
		})

	case "created":
		sort.Slice(files, func(i, j int) bool {
			if order == "asc" {
				return files[i].CreatedAt.Before(files[j].CreatedAt)
			}
			return files[i].CreatedAt.After(files[j].CreatedAt)
		})
	}

	for _, file := range files {
		fmt.Fprintln(w, file.ToString())
	}
}
