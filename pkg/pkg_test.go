package pkg

import (
	"testing"
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

	if userFolder != folder {
		t.Errorf("Expected be the same folder\n")
	}

	if userFolder.Name != foldername {
		t.Errorf("Expected foldername %s, but got %s\n", foldername, userFolder.Name)
	}

	folders := user.GetFolders()
	if len(folders) != 1 {
		t.Errorf("Expected one folder from the user\n")
	}

	if fake := user.GetFolder("fakename"); fake != nil {
		t.Errorf("Expected nil but got %v\n", fake)
	}
}

func TestFolder(t *testing.T) {
	foldername := "testfolder"
	desc := "desc for test"
	username := "testuser"

	folder := CreateFolder(foldername, desc, username)

	if folder.Name != foldername {
		t.Errorf("Expected foldername %s, but got %s\n", foldername, folder.Name)
	}

	if len(folder.Files) != 0 {
		t.Errorf("Expected empty files, but got %d\n", len(folder.Files))
	}

	newName := "newfoldername"
	folder.SetName(newName)

	if folder.Name != newName {
		t.Errorf("Expected foldername %s, but got %s\n", foldername, folder.Name)
	}
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
	if len(files) != 2 {
		t.Errorf("Expected two files, but got %d\n", len(files))
	}

	folderFile1 := folder.GetFile(filename1)

	if folderFile1 != file1 {
		t.Errorf("Expected be the same folder\n")
	}

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

	if userFile.Name != filename {
		t.Errorf("Expected foldername %s, but got %s\n", filename, userFile.Name)
	}
}

func TestSetupSystem(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()
	sys2 := SetupSystem()

	if sys != sys2 {
		t.Errorf("Expected only one system exists\n")
	}

}

func TestRegister(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")

	if _, exists := sys.UserTable["user1"]; !exists {
		t.Errorf("User `user1` doesn't register in system\n")
	}

	sys.Execute("register user2")

	if len(sys.UserTable) != 2 {
		t.Errorf("Amount of regitered user want 2, but got %d\n", len(sys.UserTable))
	}

	sys.Execute("register user1")

	if len(sys.UserTable) != 2 {
		t.Errorf("user1 cannot be registered twice\n")
	}

	sys.Execute("register u$er")
	if len(sys.UserTable) != 2 {
		t.Errorf("Names contain invalid chars cannot be used\n")
	}
}

func TestCreateFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")

	if len(sys.UserTable["user1"].Folders) != 1 {
		t.Errorf("Folder `folder1` doesn't be create for user1\n")
	}
}

func TestDeleteFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")
	sys.Execute("delete-folder user1 folder1")

	if len(sys.UserTable["user1"].Folders) != 0 {
		t.Errorf("Folder `folder1` doesn't delete from user1\n")
	}
}

func TestRenameFolder(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")

	sys.Execute("rename-folder user1 folder1 folder2")

	folder1 := sys.UserTable["user1"].GetFolder("folder1")
	folder2 := sys.UserTable["user1"].GetFolder("folder2")

	if folder1 != nil {
		t.Errorf("Folder `folder1` doesn't rename successfully\n")
	}

	if folder2.Name != "folder2" {
		t.Errorf("Expected foldername folder2, but got %s\n", folder2.Name)
	}

	sys.Execute("rename-folder user1 folder2 folder+")
	if folder2.Name != "folder2" {
		t.Errorf("Names contain invalid chars cannot be used\n")
	}
}

func TestCreateFile(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")
	sys.Execute("create-file user1 folder1 file1")
	sys.Execute("create-file user1 folder1 file2")
	sys.Execute("create-file user2 folder1 file2")

	folder1 := sys.UserTable["user1"].Folders["folder1"]
	if len(folder1.Files) != 2 {
		t.Errorf("Expected contain 2 files, but got %d\n", len(folder1.Files))
	}

	sys.Execute("create-file user1 folder1 fil@")
	if len(folder1.Files) != 2 {
		t.Errorf("Names contain invalid chars cannot be used\n")
	}
}

func TestDeleteFile(t *testing.T) {
	sys := SetupSystem()
	defer sys.Reset()

	sys.Execute("register user1")
	sys.Execute("create-folder user1 folder1")
	sys.Execute("create-file user1 folder1 file1")
	sys.Execute("delete-file user1 folder1 file1")

	folder1 := sys.UserTable["user1"].Folders["folder1"]
	if len(folder1.Files) != 0 {
		t.Errorf("Expected contain 0 files, but got %d\n", len(folder1.Files))
	}
}
