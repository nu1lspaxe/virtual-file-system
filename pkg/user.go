package pkg

type User struct {
	Name    string
	Folders []Folder
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetFolders() []Folder {
	return u.Folders
}

func (u *User) AddFolder(folder Folder) {
	u.Folders = append(u.Folders, folder)
}

func (s *System) Register(username string) RespondType {
	if !s.CharsValidator.MatchString(username) {
		return ErrInvalidChars
	}
	if s.ExistsUser(username) {
		return ErrAlreadyExists
	}

	s.UserTable[username] = User{
		Name:    username,
		Folders: []Folder{},
	}
	return Succeed
}

func (s *System) ExistsUser(username string) bool {
	for name, _ := range s.UserTable {
		if name == username {
			return true
		}
	}
	return false
}
