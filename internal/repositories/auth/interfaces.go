package auth

import (
	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

//go:generate mockgen -source=interfaces.go -destination=./auth_mock.go -package=auth

type MapperAuth interface {
	UserToEntity(userPb *authv1.User) *entityAuth.UserSudir
	PortalsToEntity(portalsPb []*authv1.Portal) []*entityAuth.Portal
	SessionToEntity(session *authv1.Session) *entityAuth.Session
	SessionToPb(session *entityAuth.Session) *authv1.Session
}

type RedirectSessionMapper interface {
	UserInfoToCreateRequestUserUnfoPb(
		userInfo *entityAuth.RedirectSessionUserInfo,
	) *redirectsessionv1.CreateSessionRequest_UserInfo
}
