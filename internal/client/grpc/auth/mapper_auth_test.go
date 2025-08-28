package auth

import (
	"net"
	"testing"
	"time"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_mapperAuth_UserToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		user *authv1.User
	}

	tests := []struct {
		name string
		args args
		want func(a args) *entityAuth.UserSudir
	}{
		{
			name: "correct",
			args: args{
				user: &authv1.User{
					CloudId: "testCloudID",
					Snils:   "000-000-000-00",
					Info: &authv1.UserInfo{
						LogonName:  "testLogin",
						Email:      "testEmail",
						Fio:        "testFIO",
						LastName:   "testLastName",
						FirstName:  "testFirstName",
						MiddleName: "testMiddleName",
					},
					Portals: []*authv1.Portal{
						{
							Id:        1,
							ShortName: "testPortalName",
							Url:       "testURL",
							LogoUrl:   "testURL",
						},
					},
					Employees: []*authv1.Employee{
						{
							Inn:   "1234567890",
							OrgId: "1234567890",
						},
					},
				},
			},
			want: func(a args) *entityAuth.UserSudir {
				return &entityAuth.UserSudir{
					CloudID:    "testCloudID",
					Login:      "testLogin",
					Email:      "testEmail",
					FIO:        "testFIO",
					LastName:   "testLastName",
					FirstName:  "testFirstName",
					MiddleName: "testMiddleName",
					SNILS:      "000-000-000-00",
					Portals: []*entityAuth.Portal{
						{
							ID:    1,
							Name:  "testPortalName",
							URL:   "testURL",
							Image: "testURL",
						},
					},
					Employees: []*entityAuth.EmployeeInfo{
						{
							Inn:   "1234567890",
							OrgID: "1234567890",
						},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}

			want := tt.want(tt.args)
			au := NewAuthMapper(f.tu)
			got := au.UserToEntity(tt.args.user)

			assert.Equal(t, want, got)
		})
	}
}

func Test_mapperAuth_SessionToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		session *authv1.Session
	}

	testUUID := uuid.New()
	testTime := time.Now()
	testTimePb := timestamppb.New(testTime)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entityAuth.Session
	}{
		{
			name: "nil",
			args: args{
				session: nil,
			},
			want: func(a args, f fields) *entityAuth.Session {
				return nil
			},
		},
		{
			name: "correct without params",
			args: args{
				session: &authv1.Session{
					UserType:           authv1.UserType_USER_TYPE_ANON,
					UserIp:             "127.0.0.1",
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
				},
			},
			want: func(a args, f fields) *entityAuth.Session {
				f.tu.EXPECT().TimestampToTime(testTimePb).Return(&testTime)
				return &entityAuth.Session{
					ID:                 (*entityAuth.SessionID)(&uuid.Nil),
					UserAuthType:       entityAuth.UserAuthTypeAnon,
					UserIP:             net.ParseIP("127.0.0.1"),
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTimePb.AsTime(),
					RefreshExpiredTime: testTimePb.AsTime(),
					CreatedTime:        testTimePb.AsTime(),
					RefreshedTime:      testTimePb.AsTime(),
					IsActive:           false,
				}
			},
		},
		{
			name: "correct with user",
			args: args{
				session: &authv1.Session{
					User: &authv1.Session_User{
						Id:        testUUID.String(),
						CloudId:   "testCloudID",
						LogonName: "testLogonName",
						Email:     "test@test.com",
						Snils:     "000-000-000-00",
					},
					UserType:           authv1.UserType_USER_TYPE_ANON,
					UserIp:             "127.0.0.1",
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
				},
			},
			want: func(a args, f fields) *entityAuth.Session {
				f.tu.EXPECT().TimestampToTime(testTimePb).Return(&testTime)
				return &entityAuth.Session{
					ID: (*entityAuth.SessionID)(&uuid.Nil),
					User: &entityAuth.User{
						ID:      testUUID,
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Login:   "testLogonName",
						Email:   "test@test.com",
					},
					UserAuthType:       entityAuth.UserAuthTypeAnon,
					UserIP:             net.ParseIP("127.0.0.1"),
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTimePb.AsTime(),
					RefreshExpiredTime: testTimePb.AsTime(),
					CreatedTime:        testTimePb.AsTime(),
					RefreshedTime:      testTimePb.AsTime(),
					IsActive:           false,
				}
			},
		},
		{
			name: "correct",
			args: args{
				session: &authv1.Session{
					Id: testUUID.String(),
					User: &authv1.Session_User{
						Id:        testUUID.String(),
						CloudId:   "testCloudID",
						LogonName: "testLogin",
						Email:     "testEmail",
						Snils:     "000-000-000-00",
						Portal: &authv1.ActivePortal{
							Id:   1,
							Name: "testPortalName",
							Url:  "testURL",
							Sid:  "testPortalSID",
						},
						Employee: &authv1.Session_User_Employee{
							Id: testUUID.String(),
						},
						Person: &authv1.Session_User_Person{
							Id: testUUID.String(),
						},
					},
					UserType: authv1.UserType_USER_TYPE_ANON,
					UserIp:   "127.0.0.1",
					Device: &authv1.Device{
						UserAgent: "testUserAgent",
					},
					SudirInfo: &authv1.SudirInfo{
						Sid:      "testSudirID",
						ClientId: "testClientID",
					},
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
					IsActive:           true,
				},
			},
			want: func(a args, f fields) *entityAuth.Session {
				f.tu.EXPECT().TimestampToTime(testTimePb).Return(&testTime)
				return &entityAuth.Session{
					ID: (*entityAuth.SessionID)(&testUUID),
					User: &entityAuth.User{
						ID:      testUUID,
						CloudID: "testCloudID",
						Login:   "testLogin",
						Email:   "testEmail",
						SNILS:   "000-000-000-00",
						Employee: &entityAuth.Employee{
							ExtID: testUUID.String(),
						},
						Person: &entityAuth.Person{
							ExtID: testUUID.String(),
						},
					},
					UserAuthType: entityAuth.UserAuthTypeAnon,
					UserIP:       net.ParseIP("127.0.0.1"),
					Device: &entityAuth.Device{
						UserAgent: "testUserAgent",
						SudirInfo: &entityAuth.SudirInfo{
							SID:      "testSudirID",
							ClientID: "testClientID",
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							ID:         1,
							Name:       "testPortalName",
							URL:        "testURL",
							IsSelected: true,
						},
						SID: "testPortalSID",
					},
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTimePb.AsTime(),
					RefreshExpiredTime: testTimePb.AsTime(),
					CreatedTime:        testTimePb.AsTime(),
					RefreshedTime:      testTimePb.AsTime(),
					IsActive:           true,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}

			want := tt.want(tt.args, f)
			au := NewAuthMapper(f.tu)
			got := au.SessionToEntity(tt.args.session)

			assert.Equal(t, want, got)
		})
	}
}

func Test_mapperAuth_UserTypeToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		userType authv1.UserType
	}

	tests := []struct {
		name string
		args args
		want func(a args) entityAuth.UserAuthType
	}{
		{
			name: "invalid",
			args: args{
				userType: -1,
			},
			want: func(a args) entityAuth.UserAuthType {
				return entityAuth.UserAuthTypeInvalid
			},
		},
		{
			name: "anon",
			args: args{
				userType: authv1.UserType_USER_TYPE_ANON,
			},
			want: func(a args) entityAuth.UserAuthType {
				return entityAuth.UserAuthTypeAnon
			},
		},
		{
			name: "auth",
			args: args{
				userType: authv1.UserType_USER_TYPE_AUTH,
			},
			want: func(a args) entityAuth.UserAuthType {
				return entityAuth.UserAuthTypeAuth
			},
		},
		{
			name: "oldauth",
			args: args{
				userType: authv1.UserType_USER_TYPE_OLD_AUTH,
			},
			want: func(a args) entityAuth.UserAuthType {
				return entityAuth.UserAuthTypeOldAuth
			},
		},
		{
			name: "service",
			args: args{
				userType: authv1.UserType_USER_TYPE_SERVICE,
			},
			want: func(a args) entityAuth.UserAuthType {
				return entityAuth.UserAuthTypeService
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}

			want := tt.want(tt.args)
			sm := NewAuthMapper(f.tu)
			got := sm.UserTypeToEntity(tt.args.userType)

			assert.Equal(t, want, got)
		})
	}
}

func Test_mapperSessions_SessionToPb(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		session *entityAuth.Session
	}

	testUUID := uuid.New()
	testTime := time.Now()
	testTimePb := timestamppb.New(testTime)
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *authv1.Session
	}{
		{
			name: "nil session",
			args: args{
				session: nil,
			},
			want: func(a args, _ fields) *authv1.Session {
				return nil
			},
		},
		{
			name: "correct without params",
			args: args{
				session: &entityAuth.Session{
					ID:                 (*entityAuth.SessionID)(&testUUID),
					UserAuthType:       entityAuth.UserAuthTypeAnon,
					UserIP:             net.ParseIP("127.0.0.1"),
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTime,
					RefreshExpiredTime: testTime,
					CreatedTime:        testTime,
					RefreshedTime:      testTime,
				},
			},
			want: func(a args, f fields) *authv1.Session {
				f.tu.EXPECT().TimeToTimestamp(&testTime).Return(testTimePb)
				return &authv1.Session{
					Id:                 testUUID.String(),
					UserIp:             "127.0.0.1",
					UserType:           authv1.UserType_USER_TYPE_ANON,
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
				}
			},
		},
		{
			name: "correct with user",
			args: args{
				session: &entityAuth.Session{
					ID: (*entityAuth.SessionID)(&testUUID),
					User: &entityAuth.User{
						ID:      testUUID,
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Login:   "testLogonName",
						Email:   "test@test.com",
					},
					UserAuthType:       entityAuth.UserAuthTypeAnon,
					UserIP:             net.ParseIP("127.0.0.1"),
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTime,
					RefreshExpiredTime: testTime,
					CreatedTime:        testTime,
					RefreshedTime:      testTime,
				},
			},
			want: func(a args, f fields) *authv1.Session {
				f.tu.EXPECT().TimeToTimestamp(&testTime).Return(testTimePb)
				return &authv1.Session{
					Id: testUUID.String(),
					User: &authv1.Session_User{
						Id:        testUUID.String(),
						CloudId:   "testCloudID",
						LogonName: "testLogonName",
						Email:     "test@test.com",
						Snils:     "000-000-000-00",
					},
					UserIp:             "127.0.0.1",
					UserType:           authv1.UserType_USER_TYPE_ANON,
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
				}
			},
		},
		{
			name: "correct with params",
			args: args{
				session: &entityAuth.Session{
					ID: (*entityAuth.SessionID)(&testUUID),
					User: &entityAuth.User{
						ID:      testUUID,
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Login:   "testLogonName",
						Email:   "test@test.com",
						Employee: &entityAuth.Employee{
							ExtID: testUUID.String(),
						},
						Person: &entityAuth.Person{
							ExtID: testUUID.String(),
						},
					},
					UserAuthType: entityAuth.UserAuthTypeAnon,
					UserIP:       net.ParseIP("127.0.0.1"),
					Device: &entityAuth.Device{
						UserAgent: "testUserAgent",
						SudirInfo: &entityAuth.SudirInfo{
							SID:      "testSudirID",
							ClientID: "testClientID",
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						SID: "testPortalSID",
						Portal: entityAuth.Portal{
							ID:   0,
							Name: "testPortalName",
							URL:  "testPortalURL",
						},
					},
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     &testTime,
					AccessExpiredTime:  testTime,
					RefreshExpiredTime: testTime,
					CreatedTime:        testTime,
					RefreshedTime:      testTime,
				},
			},
			want: func(a args, f fields) *authv1.Session {
				f.tu.EXPECT().TimeToTimestamp(&testTime).Return(testTimePb)
				return &authv1.Session{
					Id: testUUID.String(),
					Device: &authv1.Device{
						DeviceId:   "",
						Type:       authv1.DeviceType_DEVICE_TYPE_INVALID,
						OsType:     authv1.OSType_OS_TYPE_INVALID,
						AppVersion: "",
						UserAgent:  "testUserAgent",
					},
					User: &authv1.Session_User{
						Id:        testUUID.String(),
						CloudId:   "testCloudID",
						LogonName: "testLogonName",
						Email:     "test@test.com",
						Snils:     "000-000-000-00",
						Portal: &authv1.ActivePortal{
							Id:   0,
							Name: "testPortalName",
							Url:  "testPortalURL",
							Sid:  "testPortalSID",
						},
						Employee: &authv1.Session_User_Employee{
							Id: testUUID.String(),
						},
						Person: &authv1.Session_User_Person{
							Id: testUUID.String(),
						},
					},
					UserIp:   "127.0.0.1",
					UserType: authv1.UserType_USER_TYPE_ANON,
					SudirInfo: &authv1.SudirInfo{
						Sid:      "testSudirID",
						ClientId: "testClientID",
					},
					Issuer:             "testIssuer",
					Subject:            "testSubject",
					LastActiveTime:     testTimePb,
					AccessExpiredTime:  testTimePb,
					RefreshExpiredTime: testTimePb,
					CreatedTime:        testTimePb,
					RefreshedTime:      testTimePb,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}

			want := tt.want(tt.args, f)
			sm := NewAuthMapper(f.tu)
			got := sm.SessionToPb(tt.args.session)

			assert.Equal(t, want, got)
		})
	}
}

func Test_mapperSessions_UserTypeToEntity(t *testing.T) {
	type fields struct {
		tu *timeUtils.MockTimeUtils
	}
	type args struct {
		userType entityAuth.UserAuthType
	}

	tests := []struct {
		name string
		args args
		want func(a args) authv1.UserType
	}{
		{
			name: "invalid",
			args: args{
				userType: -1,
			},
			want: func(a args) authv1.UserType {
				return authv1.UserType_USER_TYPE_INVALID
			},
		},
		{
			name: "anon",
			args: args{
				userType: entityAuth.UserAuthTypeAnon,
			},
			want: func(a args) authv1.UserType {
				return authv1.UserType_USER_TYPE_ANON
			},
		},
		{
			name: "auth",
			args: args{
				userType: entityAuth.UserAuthTypeAuth,
			},
			want: func(a args) authv1.UserType {
				return authv1.UserType_USER_TYPE_AUTH
			},
		},
		{
			name: "oldauth",
			args: args{
				userType: entityAuth.UserAuthTypeOldAuth,
			},
			want: func(a args) authv1.UserType {
				return authv1.UserType_USER_TYPE_OLD_AUTH
			},
		},
		{
			name: "service",
			args: args{
				userType: entityAuth.UserAuthTypeService,
			},
			want: func(a args) authv1.UserType {
				return authv1.UserType_USER_TYPE_SERVICE
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				tu: timeUtils.NewMockTimeUtils(ctrl),
			}

			want := tt.want(tt.args)
			sm := NewAuthMapper(f.tu)
			got := sm.UserAuthTypeToPb(tt.args.userType)

			assert.Equal(t, want, got)
		})
	}
}
