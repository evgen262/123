package auth

import (
	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type redirectSessionMapper struct {
}

func NewRedirectSessionMapper() *redirectSessionMapper {
	return &redirectSessionMapper{}
}

func (m redirectSessionMapper) UserInfoToEntity(userInfo *redirectsessionv1.UserInfo) *entitySession.RedirectSessionUserInfo {
	if userInfo == nil {
		return nil
	}

	return &entitySession.RedirectSessionUserInfo{
		SessionID: userInfo.GetSessionId(),
		Email:     userInfo.GetEmail(),
		SNILS:     userInfo.GetSnils(),
		PortalURL: userInfo.GetPortalUrl(),
		TargetURL: userInfo.GetTargetUrl(),
		UserAgent: userInfo.GetUseragent(),
		IP:        userInfo.GetIp(),
	}
}

func (m redirectSessionMapper) UserInfoToCreateRequestUserUnfoPb(
	userInfo *entitySession.RedirectSessionUserInfo,
) *redirectsessionv1.CreateSessionRequest_UserInfo {
	if userInfo == nil {
		return nil
	}

	return &redirectsessionv1.CreateSessionRequest_UserInfo{
		Email:     userInfo.Email,
		Snils:     userInfo.SNILS,
		PortalUrl: userInfo.PortalURL,
		TargetUrl: userInfo.TargetURL,
		Useragent: userInfo.UserAgent,
		Ip:        userInfo.IP,
	}
}
