package pkg

import (
	"fmt"
	"os"
)

type Folder struct {
	Files       []File
	Description string
}

func CreateFolder(desc string) *Folder {
	return &Folder{
		Files:       make([]File, 0),
		Description: desc,
	}
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

	folder := CreateFolder(desc)
	user.AddFolder(foldername, folder)

	fmt.Fprintf(os.Stdout, "Create %s successfully.\n", foldername)
}
