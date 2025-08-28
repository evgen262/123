package users

import (
	viewUsers "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/users"
	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
)

type usersPresenter struct {
}

func NewUsersPresenter() *usersPresenter {
	return &usersPresenter{}
}

func (up *usersPresenter) ShortUserToView(shortUser *entityUser.ShortUser) *viewUsers.ShortUser {
	if shortUser == nil {
		return nil
	}

	u := &viewUsers.ShortUser{
		LastName:   shortUser.LastName,
		FirstName:  shortUser.FirstName,
		MiddleName: shortUser.MiddleName,
		Gender:     shortUser.Gender.String(),
		PortalData: viewUsers.PortalData{
			PersonID:   shortUser.PortalData.PersonID,
			EmployeeID: shortUser.PortalData.EmployeeID,
		},
		ImageID: shortUser.ImageID,
	}

	return u
}
