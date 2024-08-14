package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	username := "testuser"
	user := CreateUser(username)

	if user.Name != username {
		t.Errorf("Expected username %s, but got %s\n", username, user.Name)
	}

	if len(user.Folders) != 0 {
		t.Errorf("Expected empty folders, but got %d\n", len(user.Folders))
	}
}

func TestUserFolders(t *testing.T) {
	username := "testuser"
	user := CreateUser(username)

	foldername := "testfolder"
	folder := CreateFolder(foldername, "", username)

	user.AddFolder(foldername, folder)

	userFolder := user.GetFolder(foldername)
	assert.Equal(t, userFolder.Name, foldername)
	assert.Equal(t, userFolder, folder)

	folders := user.GetFolders()
	assert.Equal(t, len(folders), 1)

	fake := user.GetFolder("fakename")
	assert.Nil(t, fake)
}

func TestFolder(t *testing.T) {
	foldername := "testfolder"
	desc := "desc for test"
	username := "testuser"

	folder := CreateFolder(foldername, desc, username)
	assert.Equal(t, folder.Name, foldername)
	assert.Equal(t, len(folder.Files), 0)

	folderTo := "newfoldername"
	folder.SetName(folderTo)
	assert.Equal(t, folder.Name, folderTo)
}

func TestFolderFiles(t *testing.T) {
	foldername := "testfolder"
	desc := "desc for test"
	username := "testuser"

	folder := CreateFolder(foldername, desc, username)

	filename1 := "file1"
	filename2 := "file2"
	file1 := CreateFile(filename1, "", folder.Name, username)
	file2 := CreateFile(filename1, "desc of file2", folder.Name, username)
	folder.AddFile(filename1, file1)
	folder.AddFile(filename2, file2)

	files := folder.GetFiles()
	assert.Equal(t, len(files), 2)

	folderFile1 := folder.GetFile(filename1)
	assert.Equal(t, folderFile1, file1)
}

func TestUserFolderFile(t *testing.T) {
	filename := "testfile"
	foldername := "testfolder"
	desc := "desc for test"
	username := "testuser"

	user := CreateUser(username)

	folder := CreateFolder(foldername, desc, username)
	user.AddFolder(foldername, folder)

	file := CreateFile(filename, desc, foldername, username)
	folder.AddFile(filename, file)

	userFolder := user.GetFolder(foldername)
	userFile := userFolder.GetFile(filename)

	assert.Equal(t, userFile.Name, filename)
}

func TestSetupSystem(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	sys2 := SetupSystem()

	assert.Equal(t, sys, sys2)

}

func GetTestBufs() (outBuf, errBuf *bytes.Buffer) {
	return &bytes.Buffer{}, &bytes.Buffer{}
}

func ResetBufs(outBuf, errBuf *bytes.Buffer) {
	outBuf.Reset()
	errBuf.Reset()
}

func TestRegister(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Register(outBuf, errBuf, "user1")
	assert.Equal(t, "Add user1 successfully.\n", outBuf.String())

	if _, exists := sys.UserTable["user1"]; !exists {
		t.Errorf("User `user1` doesn't register in system\n")
	}
	ResetBufs(outBuf, errBuf)

	sys.Register(outBuf, errBuf, "user1")
	assert.Equal(t, ErrAlreadyExists.ToString("user1")+"\n", errBuf.String())
	ResetBufs(outBuf, errBuf)

	sys.Register(outBuf, errBuf, "u$er")
	assert.Equal(t, ErrInvalidChars.ToString("u$er")+"\n", errBuf.String())
	ResetBufs(outBuf, errBuf)
}

func TestCreateFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Execute("register user1")

	tests := []struct {
		user        string
		folder      string
		desc        string
		expectedOut string
		expectedErr string
	}{
		{"user1", "folder1", "", "Create folder1 successfully.\n", ""},
		{"user2", "folder1", "", "", ErrNotExists.ToString("user2") + "\n"},
		{"user1", "fo[]er1", "", "", ErrInvalidChars.ToString("fo[]er1") + "\n"},
	}

	for _, tt := range tests {
		sys.CreateFolder(outBuf, errBuf, tt.user, tt.folder, tt.desc)

		assert.Equal(t, tt.expectedOut, outBuf.String())
		assert.Equal(t, tt.expectedErr, errBuf.String())
		ResetBufs(outBuf, errBuf)
	}
}

func TestDeleteFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")

	tests := []struct {
		username    string
		foldername  string
		expectedOut string
		expectedErr string
	}{
		{"user1", "folder1", "Delete folder1 successfully.\n", ""},
		{"user2", "folder1", "", ErrNotExists.ToString("user2") + "\n"},
		{"user1", "folder2", "", ErrNotExists.ToString("folder2") + "\n"},
	}

	for _, tt := range tests {
		sys.DeleteFolder(outBuf, errBuf, tt.username, tt.foldername)
		assert.Equal(t, tt.expectedOut, outBuf.String())
		assert.Equal(t, tt.expectedErr, errBuf.String())
		ResetBufs(outBuf, errBuf)
	}
}

func TestRenameFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")

	tests := []struct {
		username    string
		folderFrom  string
		folderTo    string
		expectedOut string
		expectedErr string
	}{
		{"user1", "folder1", "folder2", "Rename folder1 to folder2 successfully.\n", ""},
		{"user2", "folder2", "folder3", "", ErrNotExists.ToString("user2") + "\n"},
		{"user1", "folder3", "folder4", "", WarnNoFolders.ToString("folder3") + "\n"},
		{"user1", "folder2", "folder+", "", ErrInvalidChars.ToString("folder+") + "\n"},
	}

	for _, tt := range tests {
		sys.RenameFolder(outBuf, errBuf, tt.username, tt.folderFrom, tt.folderTo)
		assert.Equal(t, tt.expectedOut, outBuf.String())
		assert.Equal(t, tt.expectedErr, errBuf.String())
		ResetBufs(outBuf, errBuf)
	}

}

func TestCreateFile(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")

	tests := []struct {
		username    string
		foldername  string
		filename    string
		expectedOut string
		expectedErr string
	}{
		{"user1", "folder1", "file1", "Create file1 in user1/folder1 successfully.\n", ""},
		{"user2", "folder1", "file1", "", ErrNotExists.ToString("user2") + "\n"},
		{"user1", "folder2", "file1", "", ErrNotExists.ToString("folder2") + "\n"},
		{"user1", "folder1", "f[]e1", "", ErrInvalidChars.ToString("f[]e1") + "\n"},
		{"user1", "folder1", "file1", "", ErrAlreadyExists.ToString("file1") + "\n"},
	}

	for _, tt := range tests {
		sys.CreateFile(outBuf, errBuf, tt.username, tt.foldername, tt.filename, "")
		assert.Equal(t, tt.expectedOut, outBuf.String())
		assert.Equal(t, tt.expectedErr, errBuf.String())
		ResetBufs(outBuf, errBuf)
	}

}

func TestDeleteFile(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	outBuf, errBuf := GetTestBufs()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")
	sys.Execute("create-file user1 folder1 file1")

	tests := []struct {
		username    string
		foldername  string
		filename    string
		expectedOut string
		expectedErr string
	}{
		{"user1", "folder1", "file1", "Delete file1 in user1/folder1 successfully.\n", ""},
		{"user2", "folder1", "file1", "", ErrNotExists.ToString("user2") + "\n"},
		{"user1", "folder2", "file1", "", ErrNotExists.ToString("folder2") + "\n"},
		{"user1", "folder1", "file1", "", ErrNotExists.ToString("file1") + "\n"},
	}

	for _, tt := range tests {
		sys.DeleteFile(outBuf, errBuf, tt.username, tt.foldername, tt.filename)
		assert.Equal(t, tt.expectedOut, outBuf.String())
		assert.Equal(t, tt.expectedErr, errBuf.String())
		ResetBufs(outBuf, errBuf)
	}

}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args           []string
		expectedSortBy string
		expectedOrder  string
		expectedMsg    string
	}{
		{[]string{}, "name", "asc", ""},

		{[]string{"--sort-name"}, "name", "asc", ""},
		{[]string{"--sort-created"}, "created", "asc", ""},

		{[]string{"asc"}, "name", "asc", ""},
		{[]string{"desc"}, "name", "desc", ""},

		{[]string{"--sort-created", "desc"}, "created", "desc", ""},
		{[]string{"--sort-name", "asc"}, "name", "asc", ""},

		{[]string{"--invalid-flag"}, "", "", ErrInvalidFlag.ToString()},
	}

	for _, tt := range tests {
		sortBy, order, msg := ParseArgs(tt.args)
		assert.Equal(t, tt.expectedSortBy, sortBy)
		assert.Equal(t, tt.expectedOrder, order)
		assert.Equal(t, tt.expectedMsg, msg)
	}
}

func TestCharsValidator(t *testing.T) {
	tests := []struct {
		input         string
		expectedValid bool
	}{
		{"valid_input", true},
		{"another_valid123", true},
		{"valid_input_with_underscore", true},
		{"invalid@char", false},
		{"spaces not allowed", false},
	}

	for _, tt := range tests {
		sys := SetupSystem()
		isValid := sys.CharsValidator.MatchString(tt.input)

		if isValid != tt.expectedValid {
			t.Errorf("input %q: expected validity %v, got %v", tt.input, tt.expectedValid, isValid)
		}
	}
}
