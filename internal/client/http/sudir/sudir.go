package sudir

import (
	"context"
	"fmt"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//go:generate mockgen -source=sudir.go -destination=./sudir_mock.go -package=sudir

const (
	idToken = "id_token"
)

const (
	endpointAuthPath  = "/blitz/oauth/ae"
	endpointTokenPath = "/blitz/oauth/te"
)

type client struct {
	oauth        *oauth2.Config
	parser       *jwt.Parser
	logger       ditzap.Logger
	clientID     string
	clientSecret string
	baseURL      string
}

func NewClient(baseURL, clientID, secret string, logger ditzap.Logger) *client {
	cli := &client{
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  baseURL + endpointAuthPath,
				TokenURL: baseURL + endpointTokenPath,
			},
			// RedirectURL: redirect,
			Scopes: []string{ScopeOpenId, ScopeProfile, ScopeEmail, ScopeUserInfo, ScopeEmployee, ScopeGroups},
		},
		parser:       jwt.NewParser(),
		logger:       logger,
		clientID:     clientID,
		clientSecret: secret,
		baseURL:      baseURL,
	}

	return cli
}

// AuthURL
//
//	URL для авторизации пользователя в СУДИР
func (c *client) AuthURL(options AuthURLOptions) string {
	authOptions := make([]oauth2.AuthCodeOption, 0, 5)

	if options.IsOffline {
		authOptions = append(authOptions, oauth2.AccessTypeOffline)
	} else {
		authOptions = append(authOptions, oauth2.AccessTypeOnline)
	}

	if options.ClientID != "" {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("client_id", options.ClientID))
	}

	if options.RedirectURI != "" {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("redirect_uri", options.RedirectURI))
	}

	if options.CodeChallengeMethod != "" {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("code_challenge_method", options.CodeChallengeMethod))
	}

	if options.CodeChallenge != "" {
		authOptions = append(authOptions, oauth2.SetAuthURLParam("code_challenge", options.CodeChallenge))
	}

	authURL := c.oauth.AuthCodeURL(options.State, authOptions...)
	c.logger.Debug(
		"sudir AuthURL:",
		zap.String("state", options.State),
		zap.String("callback_url", options.RedirectURI),
		zap.String("return_url", authURL),
	)
	return authURL
}

// CodeExchange
//
//	проверка кода авторизации и обмен его на токены
func (c *client) CodeExchange(ctx context.Context, code string, options CodeExchangeOptions) (*OAuthResponse, error) {
	exchangeOptions := make([]oauth2.AuthCodeOption, 0, 2)

	if options.RedirectURI != "" {
		exchangeOptions = append(exchangeOptions, oauth2.SetAuthURLParam("redirect_uri", options.RedirectURI))
	}

	if options.CodeVerifier != "" {
		exchangeOptions = append(exchangeOptions, oauth2.SetAuthURLParam("code_verifier", options.CodeVerifier))
	}

	cli := *c.oauth
	if options.ClientID != "" {
		cli.ClientID = options.ClientID
		cli.ClientSecret = options.ClientSecret
	}

	c.logger.Debug(
		"sudir CodeExchange:",
		zap.String("sudir_code", code),
		zap.String("redirect_uri", options.RedirectURI),
		zap.String("code_verifier", options.CodeVerifier),
		zap.String("client_id", cli.ClientID),
		zap.String("client_secret", cli.ClientSecret),
	)

	oauthToken, err := cli.Exchange(ctx, code, exchangeOptions...)
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			c.logger.Debug(
				"sudir CodeExchange: ошибка обмена кода",
				zap.String("error_code", rErr.ErrorCode),
				zap.String("error_description", rErr.ErrorDescription),
				zap.String("error_uri", rErr.ErrorURI),
				zap.ByteString("error_body", rErr.Body),
			)

			// проверка типа ошибки oauth (неверный код авторизации)
			if rErr.ErrorCode == "invalid_grant" {
				return nil, ErrInvalidGrant
			}
		}
		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	c.logger.Debug(
		"sudir CodeExchange:",
		zap.String("access_token", oauthToken.AccessToken),
		zap.String("refresh_token", oauthToken.RefreshToken),
		zap.Time("expiry", oauthToken.Expiry),
	)

	response := &OAuthResponse{
		AccessToken:  oauthToken.AccessToken,
		RefreshToken: oauthToken.RefreshToken,
		Expiry:       &oauthToken.Expiry,
	}

	tk, ok := oauthToken.Extra(idToken).(string)
	if !ok || tk == "" {
		c.logger.Error("в ответе СУДИР отсутствует IDToken",
			entity.LogModuleUE,
			entity.LogCode("UE_010"),
			zap.Error(err),
		)
		return nil, ErrNoJWTToken
	}

	c.logger.Debug(
		"sudir CodeExchange:",
		zap.String("id_token", tk),
	)

	response.IDToken = tk
	return response, nil
}

func (c *client) LoginCredentials(ctx context.Context, options LoginOptions) (*OAuthResponse, error) {
	cli := &clientcredentials.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		TokenURL:     c.baseURL + endpointTokenPath,
		Scopes:       []string{ScopeOpenId, ScopeProfile, ScopeEmail, ScopeUserInfo, ScopeEmployee},
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	oauthToken, err := cli.Token(ctx)
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			c.logger.Debug(
				"sudir LoginCredentials: ошибка входа",
				zap.String("error_code", rErr.ErrorCode),
				zap.String("error_description", rErr.ErrorDescription),
				zap.String("error_uri", rErr.ErrorURI),
				zap.ByteString("error_body", rErr.Body),
			)

			switch rErr.ErrorCode {
			case "invalid_client":
				return nil, ErrInvalidClient
			case "access_denied":
				c.logger.Warn("sudir LoginCredentials: доступ по динамическим реквезитам запрещён",
					entity.LogModuleUE,
					entity.LogCode("UE_013"),
					zap.String("client_id", options.ClientID),
					zap.Error(err),
				)
				return nil, ErrAccessDenied
			}
		} else {
			c.logger.Error("sudir LoginCredentials: ошибка входа по динамическим реквезитам доступа",
				entity.LogModuleUE,
				entity.LogCode("UE_013"),
				zap.String("client_id", options.ClientID),
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	c.logger.Debug(
		"sudir LoginCredentials:",
		zap.String("access_token", oauthToken.AccessToken),
		zap.Time("expiry", oauthToken.Expiry),
	)

	response := &OAuthResponse{
		AccessToken: oauthToken.AccessToken,
		Expiry:      &oauthToken.Expiry,
	}

	return response, nil
}

// RefreshToken
//
//	проверка валидности токена и обмен его на новый
func (c *client) RefreshToken(ctx context.Context, refreshToken string) (*OAuthResponse, error) {
	oauthToken, err := c.oauth.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			c.logger.Debug(
				"sudir RefreshToken:",
				zap.String("error_code", rErr.ErrorCode),
				zap.String("error_description", rErr.ErrorDescription),
				zap.String("error_uri", rErr.ErrorURI),
				zap.ByteString("error_body", rErr.Body),
			)

			// проверка типа ошибки oauth (неверный refresh token)
			if rErr.ErrorCode == "invalid_grant" {
				return nil, ErrInvalidGrant
			}
		} else {
			c.logger.Error("sudir RefreshToken: ошибка обновления access_token", zap.Error(err))
		}

		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	c.logger.Debug(
		"sudir RefreshToken:",
		zap.String("access_token", oauthToken.AccessToken),
		zap.String("refresh_token", oauthToken.RefreshToken),
		zap.Time("expiry", oauthToken.Expiry),
	)

	return &OAuthResponse{
		AccessToken:  oauthToken.AccessToken,
		RefreshToken: oauthToken.RefreshToken,
		Expiry:       &oauthToken.Expiry,
	}, nil
}

// GetUserInfo
//
//	Обмен access_token на информацию о пользователе
func (c *client) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	cli := NewExtendedClient(c.baseURL)

	userInfo, err := cli.GetUserInfo(ctx, accessToken)
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			c.logger.Error(
				"sudir GetUserInfo: ошибка получения информации о пользователе",
				entity.LogModuleUE,
				entity.LogCode("UE_015"),
				zap.String("error_code", rErr.ErrorCode),
				zap.String("error_description", rErr.ErrorDescription),
				zap.String("error_uri", rErr.ErrorURI),
				zap.ByteString("error_body", rErr.Body),
			)
		} else {
			c.logger.Error(
				"sudir GetUserInfo: ошибка получения информации о пользователе",
				entity.LogModuleUE,
				entity.LogCode("UE_015"),
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	c.logger.Debug(
		"sudir GetUserInfo:",
		zap.String("sub", userInfo.Sub),
		zap.String("logonname", userInfo.LogonName),
		zap.String("family_name", userInfo.FamilyName),
		zap.String("name", userInfo.Name),
		zap.String("middle_name", userInfo.MiddleName),
		zap.String("company", userInfo.Company),
		zap.String("department", userInfo.Department),
		zap.String("position", userInfo.Position),
	)

	return userInfo, nil
}

func (c *client) Logout(ctx context.Context, clientID, registrationToken string) error {
	cli := NewExtendedClient(c.baseURL)

	err := cli.Logout(ctx, clientID, registrationToken)
	if err != nil {
		if rErr, ok := err.(*oauth2.RetrieveError); ok {
			c.logger.Error(
				"sudir Logout: ошибка при выходе",
				entity.LogModuleUE,
				entity.LogCode("UE_017"),
				zap.String("client_id", clientID),
				zap.String("error_code", rErr.ErrorCode),
				zap.String("error_description", rErr.ErrorDescription),
				zap.String("error_uri", rErr.ErrorURI),
				zap.ByteString("error_body", rErr.Body),
			)

			// проверка типа ошибки oauth (неверный refresh token)
			if rErr.ErrorCode == "access_denied" {
				return ErrAccessDenied
			}
		} else {
			c.logger.Error(
				"sudir Logout: ошибка при выходе",
				entity.LogModuleUE,
				entity.LogCode("UE_017"),
				zap.String("client_id", clientID),
				zap.Error(err),
			)
		}
		return fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	return nil
}

// ValidateToken
//
//	валидация access_token
func (c *client) ValidateToken(ctx context.Context, accessToken string) (*ValidationInfo, error) {
	cli := NewExtendedClient(c.baseURL)

	vInfo, err := cli.ValidateToken(ctx, c.clientID, c.clientSecret, accessToken)
	if err != nil {
		c.logger.Error(
			"sudir ValidateToken: ошибка ополучения информации об access_token",
			entity.LogModuleUE,
			entity.LogCode("UE_016"),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}

	c.logger.Debug(
		"sudir ValidateToken:",
		zap.String("sub", vInfo.Sub),
		zap.String("scope", vInfo.Scope),
		zap.String("jti", vInfo.Jti),
		zap.String("token_type", vInfo.TokenType),
		zap.String("client_id", vInfo.ClientID),
		zap.Bool("active", vInfo.IsActive),
		zap.Time("iat", time.Unix(vInfo.Iat, 0)),
		zap.Time("exp", time.Unix(vInfo.Exp, 0)),
	)

	return vInfo, nil
}

// ParseToken
//
//	получение payload из jwt токена
func (c *client) ParseToken(token string) (*JWTPayload, error) {
	payload := &JWTPayload{}
	_, _, err := c.parser.ParseUnverified(token, payload)
	if err != nil {
		c.logger.Error("ошибка парсинга IDToken",
			entity.LogModuleUE,
			entity.LogCode("UE_011"),
			zap.Error(err),
			zap.String("id_token", token),
		)
		return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
	}
	return payload, nil
}
