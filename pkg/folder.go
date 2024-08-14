package pkg

import (
	"fmt"
	"time"
)

var MaxDescLength = 20

type Folder struct {
	Name        string
	Description string
	Files       map[string]*File
	CreatedAt   time.Time
	UserName    string
}

func CreateFolder(foldername, desc, username string) *Folder {
	return &Folder{
		Name:        foldername,
		Description: desc,
		Files:       make(map[string]*File, 0),
		CreatedAt:   time.Now(),
		UserName:    username,
	}
}

func (folder *Folder) SetName(foldername string) {
	folder.Name = foldername
}

func (folder *Folder) GetFile(filename string) *File {
	for f := range folder.Files {
		if f == filename {
			return folder.Files[f]
		}
	}
	return nil
}

func (folder *Folder) GetFiles() []*File {
	var files []*File
	for f := range folder.Files {
		files = append(files, folder.Files[f])
	}
	return files
}

func (folder *Folder) AddFile(filename string, file *File) {
	folder.Files[filename] = file
}

func (folder *Folder) ToString() string {
	return fmt.Sprintf("%s %s %s %s",
		folder.Name,
		folder.Description,
		folder.CreatedAt.Format("2006-01-02 15:04:05"),
		folder.UserName,
	)
}
