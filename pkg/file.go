package pkg

import (
	"fmt"
	"time"
)

type File struct {
	Name        string
	Description string
	CreatedAt   time.Time
	FolderName  string
	UserName    string
}

func CreateFile(filename, desc, foldername, username string) *File {
	return &File{
		Name:        filename,
		Description: desc,
		CreatedAt:   time.Now(),
		FolderName:  foldername,
		UserName:    username,
	}
}

func (file *File) ToString() string {
	return fmt.Sprintf("%s %s %s %s %s",
		file.Name,
		file.Description,
		file.CreatedAt.Format("2006-01-02 15:04:05"),
		file.FolderName,
		file.UserName,
	)
}
