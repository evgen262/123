package sudir

// ErrorResponse
//
//	https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
type ErrorResponse struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
}
