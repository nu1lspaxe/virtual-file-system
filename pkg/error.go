package pkg

import "fmt"

type RespondType int

const (
	Succeed RespondType = -1

	ErrAlreadyExists RespondType = iota
	ErrInvalidChars
	ErrNotExists
	ErrArgsLength
	ErrInvalidFlag
	ErrUnknownCmd

	WarnNoFolders
	WarnEmptyFolder
)

func (r RespondType) ToString(item ...string) string {
	switch r {
	case Succeed:
		return ""
	case ErrAlreadyExists:
		return fmt.Sprintf("Error: The %v has already existed.", item)
	case ErrInvalidChars:
		return fmt.Sprintf("Error: The %v contain invalid chars.", item)
	case ErrNotExists:
		return fmt.Sprintf("Error: The %v doesn't exist.", item)
	case ErrArgsLength:
		return "Error: Invalid command syntax. Check `help` to get info!"
	case ErrInvalidFlag:
		return "Error: Invalid flags. They can be [--sort-name|--sort-created] [asc|desc]."
	case ErrUnknownCmd:
		return "Unrecognized command."
	case WarnNoFolders:
		return fmt.Sprintf("Warning: The %v doesn't have any folders", item)
	case WarnEmptyFolder:
		return "Warning: The folder is empty."
	default:
		return "Undefined"
	}
}
