# virtual-file-system

## Functionality

**Virtual File System** has mainly three functionalities: 
- User Management
- Folder Management
- File Management

The design pattern of the system follows Singleton pattern. 

The system will be created ans only once as the program start running.

<br>

### User Management
- Allow users to register a unique, case insensitive username.
- Users can have an arbitrary number of folders and files.

#### Commands

```bash
register [username]
```

### Folder Management
- Users can create, delete, and rename folders.
- Folder names must be unique within the user's scope and are case insensitive.
- Folders have an optional description field.

#### Commands

```bash
create-folder [username] [foldername] [description]?

delete-folder [username] [foldername]

list-folders [username] [--sort-name|--sort-created] [asc|desc]

rename-folder [username] [foldername] [new-folder-name]
```

### File Management
- Users can create, delete, and list all files within a specified folder.
- File names must be unique within the same folder and are case insensitive.
- Files have an optional description field.

#### Commands

```bash
create-file [username] [foldername] [filename] [description]?

delete-file [username] [foldername] [filename]

list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]
```

:exclamation: Name of the User | Folder | File are only acceptable with character (a-zA-Z), integer (0-9) and underscore (_)

---

## CLI to start

- Start the program
  - Run main file directly
    ```bash
    go run main.go
    ```
  - Run binary file
    ```bash
    # go build -o vfs main.go
    ./vfs # only works in linux-based system
    ```

- Get help information (by `-h` or `--help`)

    ```bash
    go run main.go -h
    ```