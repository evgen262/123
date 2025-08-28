package files

import (
	"testing"

	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

func Test_fileMapper_FilePbToEntity(t *testing.T) {
	type args struct {
		f *filev1.File
	}
	tests := []struct {
		name string
		args args
		want func(a args) *entityFile.File
	}{
		{
			name: "success",
			args: args{
				f: &filev1.File{
					Id: "b6d1395a-4c64-4636-a1a0-c5460e3e92c6",
					Meta: &filev1.Meta{
						ContentType: "image/png",
					},
					Permissions: &filev1.Permissions{},
				},
			},
			want: func(a args) *entityFile.File {
				return &entityFile.File{
					Id: uuid.MustParse("b6d1395a-4c64-4636-a1a0-c5460e3e92c6"),
					Metadata: entityFile.Metadata{
						ContentType: "image/png",
					},
					Permissions: &entityFile.Permissions{},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want(tt.args)

			fm := NewFileMapper()
			got := fm.FilePbToEntity(tt.args.f)

			assert.Equal(t, want, got)
		})
	}
}
