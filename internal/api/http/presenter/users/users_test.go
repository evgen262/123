package users

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	viewUsers "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/users"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
)

func Test_employeePresenter_PersonToShortView(t *testing.T) {
	type args struct {
		shortUser *entityUser.ShortUser
	}

	testString := "testString"
	testID := uuid.New()

	tests := []struct {
		name string
		args args
		want *viewUsers.ShortUser
	}{
		{
			name: "nil",
			args: args{
				shortUser: nil,
			},
			want: nil,
		},
		{
			name: "correct",
			args: args{
				shortUser: &entityUser.ShortUser{
					LastName:   testString,
					FirstName:  testString,
					MiddleName: testString,
					ImageID:    testID.String(),
					Gender:     entity.GenderFemale,
					PortalData: entityUser.PortalData{
						PersonID:   testID.String(),
						EmployeeID: testID.String(),
					},
				},
			},
			want: &viewUsers.ShortUser{
				LastName:   testString,
				FirstName:  testString,
				MiddleName: testString,
				ImageID:    testID.String(),
				Gender:     "female",
				PortalData: viewUsers.PortalData{
					PersonID:   testID.String(),
					EmployeeID: testID.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := NewUsersPresenter()
			got := ep.ShortUserToView(tt.args.shortUser)

			assert.Equal(t, tt.want, got)
		})
	}
}
