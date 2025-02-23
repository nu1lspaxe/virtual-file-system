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
    :exclamation: Binary file cannot be executed in powershell
    ```bash
    # go build -o vfs main.go
    ./vfs
    ```
    

- Get help information (by `-h` or `--help`)

    ```bash
    go run main.go -h
    ```

---


### Test Record
```bash
go test -cover ./...
# ok      system/pkg      0.005s  coverage: 54.7% of statements
```

---

### Verions 

#### v1.0.0
Essential functions are released.

##### v1.0.1
Modify help information display in window/linux system.
Handbook will be shown even cannot find `man` command.

##### v1.0.2
Fix `rename-folder` function (didn't get folder from user)

##### v1.0.3
Append unit tests (coverage: 54.7% of statemetns)