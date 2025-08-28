package auth

type AuthResponse struct {
	PortalSession string `json:"portalSession,omitempty"`
	// User          User      `json:"user"`
	Portals []*Portal `json:"portals"`
}

func NewAuthResponse(portalSession string, portals []*Portal) *AuthResponse {
	return &AuthResponse{
		PortalSession: portalSession,
		Portals:       portals,
	}
}
