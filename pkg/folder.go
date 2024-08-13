package pkg

import "time"

type Folder struct {
	Name        string
	Description string
	Files       []File
	CreatedAt   time.Time
}

func CreateFolder(foldername, desc string) *Folder {
	return &Folder{
		Name:        foldername,
		Description: desc,
		Files:       make([]File, 0),
		CreatedAt:   time.Now(),
	}
}
