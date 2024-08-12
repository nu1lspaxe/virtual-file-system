# virtual-file-system

## Functionality


### User Management
- Allow users to register a unique, case insensitive username.
- Users can have an arbitrary number of folders and files.

### Folder Management
- Users can create, delete, and rename folders.
- Folder names must be unique within the user's scope and are case insensitive.
- Folders have an optional description field.

### File Management
- Users can create, delete, and list all files within a specified folder.
- File names must be unique within the same folder and are case insensitive.
- Files have an optional description field.

---

## Commands
- Get help information (under linux system)
```bash
# run main.go
go run main.go -h
# or run binary file
./vfs -h
```