package auth

import (
	viewAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type authPresenter struct {
}

func NewAuthPresenter() *authPresenter {
	return &authPresenter{}
}

func (p authPresenter) AuthToView(authInfo *entityAuth.Auth) *viewAuth.AuthResponse {
	return viewAuth.NewAuthResponse(authInfo.PortalSession, p.PortalsToView(authInfo.Portals))
}

func (p authPresenter) PortalsToView(portals []*entityAuth.Portal) []*viewAuth.Portal {
	portalsView := make([]*viewAuth.Portal, len(portals))
	hasSelectedPortal := false

	for i, portal := range portals {
		portalsView[i] = p.PortalToView(portal)
		if portal.IsSelected {
			hasSelectedPortal = true
		}
	}

	// Если нет выбранного портала, то по умолчанию ставим активным первый портал
	if !hasSelectedPortal && len(portalsView) > 0 {
		portalsView[0].IsActive = true
	}

	return portalsView
}

func (p authPresenter) PortalToView(portal *entityAuth.Portal) *viewAuth.Portal {
	return &viewAuth.Portal{
		ID:       portal.ID,
		Name:     portal.Name,
		URL:      portal.URL,
		Image:    portal.Image,
		IsActive: portal.IsSelected,
	}
}
