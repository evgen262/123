package entity

type State struct {
	ID           string `json:"ID"`
	ClientID     string `json:"ClientID,omitempty"`
	ClientSecret string `json:"ClientSecret,omitempty"`
	CodeVerifier string `json:"CodeVerifier,omitempty"`
	DeviceID     string `json:"DeviceID,omitempty"`
	UserAgent    string `json:"UserAgent,omitempty"`
	CallbackURL  string `json:"CallbackURL,omitempty"`
}

type StateOptions struct {
	ClientID     string `json:"ClientID,omitempty"`
	ClientSecret string `json:"ClientSecret,omitempty"`
	CodeVerifier string `json:"CodeVerifier,omitempty"`
	DeviceID     string `json:"DeviceID,omitempty"`
	UserAgent    string `json:"UserAgent,omitempty"`
	CallbackURL  string `json:"CallbackURL,omitempty"`
}
