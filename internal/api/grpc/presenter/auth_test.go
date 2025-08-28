package presenter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

func Test_authPresenter_UserToPb(t *testing.T) {
	type args struct {
		user *entity.User
	}
	tests := []struct {
		name string
		args args
		want *authv1.User
	}{
		{
			name: "ok no employees no info",
			args: args{user: &entity.User{
				CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			}},
			want: &authv1.User{
				CloudId:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				Employees: []*authv1.Employee{},
			},
		},
		{
			name: "ok with employees no info",
			args: args{user: &entity.User{
				CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				Employees: []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "770987654321",
						OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:     "Иванов Иван Иванович",
					},
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "771234567890",
						OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
						FIO:     "Иванов Иван Иванович",
					},
				},
			}},
			want: &authv1.User{
				CloudId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				Employees: []*authv1.Employee{
					{
						Inn:   "770987654321",
						OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Fio:   "Иванов Иван Иванович",
					},
					{
						Inn:   "771234567890",
						OrgId: "71bac977-ea7f-4156-9504-60f7d443ab62",
						Fio:   "Иванов Иван Иванович",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewAuthPresenter()

			got := au.UserToPb(tt.args.user)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_authPresenter_EmployeesToPb(t *testing.T) {
	type args struct {
		employees []entity.EmployeeInfo
	}
	tests := []struct {
		name string
		args args
		want []*authv1.Employee
	}{
		{
			name: "correct",
			args: args{
				employees: []entity.EmployeeInfo{
					{
						CloudID: "3c5cbb16-011a-310e-97e2-565400a26506",
						Inn:     "770123456789",
						OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:     "Иванов Иван Иванович",
					},
					{
						CloudID: "3c5cbb16-011a-310e-97e2-565400a26506",
						Inn:     "771234567890",
						OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
						FIO:     "Иванов Иван Иванович",
					},
				},
			},
			want: []*authv1.Employee{
				{
					Inn:   "770123456789",
					OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					Fio:   "Иванов Иван Иванович",
				},
				{
					Inn:   "771234567890",
					OrgId: "71bac977-ea7f-4156-9504-60f7d443ab62",
					Fio:   "Иванов Иван Иванович",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewAuthPresenter()

			got := au.EmployeesToPb(tt.args.employees)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_authPresenter_UserInfoToPb(t *testing.T) {
	type args struct {
		userInfo *entity.UserInfo
	}
	tests := []struct {
		name string
		args args
		want *authv1.UserInfo
	}{
		{
			name: "ok",
			args: args{userInfo: &entity.UserInfo{
				Sub:        "test@some.domain.mos.ru",
				LastName:   "Testov",
				FirstName:  "Test",
				MiddleName: "Testovich",
				LogonName:  "TestovTT",
				Email:      "test@it.mos.ru",
				Company:    "OOO Test",
				Department: "Otdel ooo test",
				Position:   "developer",
			}},
			want: &authv1.UserInfo{
				Email:           "test@it.mos.ru",
				LogonName:       "TestovTT",
				Sub:             "test@some.domain.mos.ru",
				LastName:        "Testov",
				FirstName:       "Test",
				MiddleName:      "Testovich",
				SudirCompany:    "OOO Test",
				SudirDepartment: "Otdel ooo test",
				SudirPosition:   "developer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			au := NewAuthPresenter()

			got := au.UserInfoToPb(tt.args.userInfo)
			assert.Equal(t, tt.want, got)
		})
	}
}
