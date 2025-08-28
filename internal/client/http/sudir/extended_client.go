package sudir

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const bodyMaxSize = 1024000

const (
	endpointUserPath   = "/blitz/oauth/me"
	endpointIntrospect = "/blitz/oauth/introspect"
	endpointRegister   = "/blitz/oauth/register"
)

type extendedClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewExtendedClient(baseURL string) *extendedClient {
	return &extendedClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *extendedClient) newHTTPRequest(method, methodURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, methodURL, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *extendedClient) retrieveError(resp *http.Response, body []byte) *oauth2.RetrieveError {
	rErr := &oauth2.RetrieveError{
		Response: resp,
		Body:     body,
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return rErr
	}

	rErr.ErrorCode = errResp.ErrorCode
	rErr.ErrorDescription = errResp.ErrorDescription
	rErr.ErrorURI = errResp.ErrorURI

	return rErr
}

func (c *extendedClient) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	httpReq, err := c.newHTTPRequest(http.MethodGet, c.baseURL+endpointUserPath, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot get user info: %v", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(httpResp.Body, bodyMaxSize))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch user info: %v", err)
	}
	if code := httpResp.StatusCode; code < 200 || code > 299 {
		return nil, c.retrieveError(httpResp, body)
	}

	var uInfo UserInfo
	if err := json.Unmarshal(body, &uInfo); err != nil {
		return nil, fmt.Errorf("cannot parse user info response: %v", err)
	}

	return &uInfo, nil
}

func (c *extendedClient) ValidateToken(ctx context.Context, clientID, clientSecret, accessToken string) (*ValidationInfo, error) {
	params := url.Values{}
	params.Set("token", accessToken)
	params.Set("token_type_hint", "access_token")

	httpReq, err := c.newHTTPRequest(http.MethodPost, c.baseURL+endpointIntrospect, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.SetBasicAuth(clientID, clientSecret)
	httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("cannot get validation info: %v", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(httpResp.Body, bodyMaxSize))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch validation info: %v", err)
	}
	if code := httpResp.StatusCode; code < 200 || code > 299 {
		return nil, c.retrieveError(httpResp, body)
	}

	// ValidationInfo
	var vInfo ValidationInfo
	if err := json.Unmarshal(body, &vInfo); err != nil {
		return nil, fmt.Errorf("cannot parse validation info response: %v", err)
	}

	return &vInfo, nil
}

func (c *extendedClient) Logout(ctx context.Context, clientID, registrationToken string) error {
	httpReq, err := c.newHTTPRequest(http.MethodDelete, c.baseURL+endpointRegister+"/"+clientID, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+registrationToken)
	httpReq.WithContext(ctx)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("cannot logout: %v", err)
	}
	defer httpResp.Body.Close()

	if code := httpResp.StatusCode; code < 200 || code > 299 {
		body, err := io.ReadAll(io.LimitReader(httpResp.Body, bodyMaxSize))
		if err != nil {
			return fmt.Errorf("cannot fetch logout response: %v", err)
		}

		return c.retrieveError(httpResp, body)
	}

	return nil
}
