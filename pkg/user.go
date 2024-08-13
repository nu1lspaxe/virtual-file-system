package pkg

import (
	"fmt"
	"os"
)

type User struct {
	Name    string
	Folders map[string]*Folder
}

func CreateUser(username string) *User {
	return &User{
		Name:    username,
		Folders: make(map[string]*Folder),
	}
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetFolder(foldername string) *Folder {
	for f := range u.Folders {
		if f == foldername {
			return u.Folders[f]
		}
	}
	return nil
}

func (u *User) GetFolders() map[string]*Folder {
	return u.Folders
}

func (u *User) AddFolder(foldername string, folder *Folder) {
	u.Folders[foldername] = folder
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
