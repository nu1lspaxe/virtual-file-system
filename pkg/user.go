package pkg

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

func (u *User) GetFolders() []*Folder {
	var folders []*Folder
	for f := range u.Folders {
		folders = append(folders, u.Folders[f])
	}
	return folders
}

func (u *User) AddFolder(foldername string, folder *Folder) {
	u.Folders[foldername] = folder
}
