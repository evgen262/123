package sudir

type UserInfo struct {
	Sub        string `json:"sub,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	Name       string `json:"name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LogonName  string `json:"logonname,omitempty"`
	Company    string `json:"company,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	Email      string `json:"email,omitempty"`
}

type ValidationInfo struct {
	Sub       string `json:"sub,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Jti       string `json:"jti,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	IsActive  bool   `json:"active,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
}
