package file

import (
	"fmt"

	"github.com/google/uuid"
)

type File struct {
	Id          uuid.UUID
	Payload     []byte
	Name        string
	Metadata    Metadata
	Permissions *Permissions
}

func (file *File) GetFileName() string {
	if file == nil {
		return ""
	}

	return fmt.Sprintf("%s.%s", file.Name, file.Metadata.Extension)
}

func (file *File) GetPermissions() *Permissions {
	if file == nil || file.Permissions == nil {
		return nil
	}

	return file.Permissions
}
