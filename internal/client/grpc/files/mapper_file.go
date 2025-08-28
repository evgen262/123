package files

import (
	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	"github.com/google/uuid"

	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

type fileMapper struct{}

func NewFileMapper() *fileMapper {
	return &fileMapper{}
}

func (fm fileMapper) FilePbToEntity(file *filev1.File) *entityFile.File {
	if file == nil {
		return nil
	}

	id, err := uuid.Parse(file.GetId())
	if err != nil {
		return nil
	}

	pbMetadata := file.GetMeta()
	if pbMetadata == nil {
		return nil
	}

	var permissions *entityFile.Permissions
	pbPermissions := file.GetPermissions()
	if pbPermissions != nil {
		permissions = &entityFile.Permissions{
			UserIds:    pbPermissions.GetUserIds(),
			OrgUnitIds: pbPermissions.GetOrgUnitIds(),
		}
	}

	return &entityFile.File{
		Id:      id,
		Payload: file.GetPayload(),
		Name:    file.GetName(),
		Metadata: entityFile.Metadata{
			ContentType: pbMetadata.GetContentType(),
			Size:        pbMetadata.GetSize(),
			OwnerId:     pbMetadata.GetOwnerId(),
			Extension:   file.GetFileExtension(),
		},
		Permissions: permissions,
	}
}
