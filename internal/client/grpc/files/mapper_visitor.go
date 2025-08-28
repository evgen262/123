package files

import (
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/shared/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type visitorMapper struct{}

func NewVisitorMapper() *visitorMapper {
	return &visitorMapper{}
}

func (v *visitorMapper) SessionToVisitorPb(session *auth.Session) *sharedv1.Visitor {
	switch session.UserAuthType {
	case auth.UserAuthTypeAuth, auth.UserAuthTypeOldAuth:
		user := session.GetUser()
		if user == nil {
			return nil
		}
		return &sharedv1.Visitor{
			Visitor: &sharedv1.Visitor_User_{
				User: &sharedv1.Visitor_User{
					Id: user.ID.String(),
				},
			},
		}
	default:
		return nil
	}
}
