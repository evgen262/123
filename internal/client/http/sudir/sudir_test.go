package sudir

/*
import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
)

type sudirSuite struct {
	suite.Suite

	oauth       *MockOAuth2
	tokenSource *MockTokenSource

	baseUrl         string
	testState       string
	testCode        string
	testCallbackURL string

	logger *ditzap.MockLogger

	ctx context.Context
}

func (ss *sudirSuite) SetupTest() {
	ctrl := gomock.NewController(ss.T())

	ss.oauth = NewMockOAuth2(ctrl)
	ss.tokenSource = NewMockTokenSource(ctrl)

	ss.logger = ditzap.NewMockLogger(ctrl)

	ss.baseUrl = "https://some/service/url/"
	ss.testState = "71bac977-ea7f-4156-9504-60f7d443ab62"
	ss.testCode = "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q"
	ss.testCallbackURL = "https://callback/url"
	ss.ctx = context.TODO()

	if location, err := time.LoadLocation("Europe/Moscow"); err == nil {
		time.Local = location
	}
}

func (ss *sudirSuite) Test_AuthURL() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	want := ss.baseUrl +
		endpointAuthPath +
		"?access_type=offline&client_id=clientID&response_type=code&scope=openid%2Bprofile%2Bemail%2Buserinfo%2Bemployee%2Bgroups&state=" + ss.testState +
		"&redirect_uri=https%3A%2F%2Fcallback%2Furl"
	ss.oauth.EXPECT().AuthCodeURL(ss.testState, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", ss.testCallbackURL)).Return(want)
	ss.logger.EXPECT().Debug(
		"sudir AuthURL:",
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	)

	options := AuthURLOptions{
		IsOffline:   true,
		RedirectURI: ss.testCallbackURL,
		State:       ss.testState,
	}

	got := c.AuthURL(options)
	ss.Equal(want, got)
}

func (ss *sudirSuite) Test_CodeExchange_ExchangeErr() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	testErr := errors.New("some test error")

	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any())
	ss.oauth.EXPECT().Exchange(ss.ctx, ss.testCode, oauth2.SetAuthURLParam("redirect_uri", ss.testCallbackURL)).Return(nil, testErr)

	got, err := c.CodeExchange(ss.ctx, ss.testCode, ss.testCallbackURL)

	ss.Nil(got)
	ss.EqualError(err, fmt.Errorf("%s: %w", ErrStatusTextServiceError, testErr).Error())
}

func (ss *sudirSuite) Test_CodeExchange_NoIDToken() {
	testExpiry := time.Now().Local().Add(time.Hour)
	testOauthToken := &oauth2.Token{
		AccessToken:  "wjsTCFTERhLa86xYLy4mPZjvx7RSHV9oUwQ8V3zxKMs1MDMxMTZlMS04Yjk1LTRmNDEtOGYwZi1jNzEwOGI4NzY4ZWU",
		TokenType:    "Bearer",
		RefreshToken: "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg",
		Expiry:       testExpiry,
	}
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any())
	ss.oauth.EXPECT().Exchange(ss.ctx, ss.testCode, oauth2.SetAuthURLParam("redirect_uri", ss.testCallbackURL)).Return(testOauthToken, nil)
	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any(), gomock.Any())

	got, err := c.CodeExchange(ss.ctx, ss.testCode, ss.testCallbackURL)
	ss.ErrorIs(err, ErrNoJWTToken)
	ss.Nil(got)
}

func (ss *sudirSuite) Test_CodeExchange_InvalidGrant() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any())
	ss.oauth.EXPECT().Exchange(ss.ctx, ss.testCode, oauth2.SetAuthURLParam("redirect_uri", ss.testCallbackURL)).Return(nil, &oauth2.RetrieveError{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
	})
	ss.logger.EXPECT().Debug("sudir CodeExchange: ошибка обмена кода", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	got, err := c.CodeExchange(ss.ctx, ss.testCode, ss.testCallbackURL)
	ss.ErrorIs(err, ErrInvalidGrant)
	ss.Nil(got)
}

func (ss *sudirSuite) Test_CodeExchange_Ok() {
	testIDToken := "eyJraWQiOiJkZWZhdWx0IiwiYWxnIjoiUlMyNTYifQ.eyJlbWFpbCI6ImRldmVsb3BlckBpdC5tb3MucnUiLCJsb2dvbm5hbWUiOiJEZXZlbG9wZXIiLCJ1YV9pZCI6IlNIQTI1Nl91TVcwN0xaWWUtcmVvbHhwbjVqMGRkNmhUcm5zT09DYmJuVlNTVVlnZ3VjIiwianRpIjoiSzVIaU1nblNCYW0tSDR6dmhzM3NROHFDRjdaU0pTek9ZTUJWS2pheDJzYyIsImV4cCI6MTY4ODQxNjEwMCwiY2xvdWRHVUlEIjoiM2M1Y2JiMTYtMDExYS0zMTBlLTk3ZTItNTY1NDAwYTI2NTA2IiwiaWF0IjoxNjg4MzgzNzAwLCJhdWQiOlsiY2ZjLXp2LmRpdGNsb3VkLnJ1Il0sImFtciI6WyJwYXNzd29yZCJdLCJpc3MiOiJodHRwczovL3N1ZGlyLXRlc3QubW9zLnJ1IiwiY29tcGFueSI6ItCT0JrQoyDQmNC90YTQvtCz0L7RgNC-0LQiLCJkZXBhcnRtZW50Ijoi0J7RgtC00LXQuyDQvNC-0LHQuNC70YzQvdC-0Lkg0YDQsNC30YDQsNCx0L7RgtC60LgiLCJwb3NpdGlvbiI6ItC_0YDQvtCz0YDQsNC80LzQuNGB0YIiLCJzdWIiOiJEZXZlbG9wZXJAaHEuY29ycC5tb3MucnUiLCJjcmlkIjoiMCIsInNpZCI6IjUxZTEwMzYxLThhZTUtNGYwMS04ZjRmLWM3YjE4ZTA4ODc2ZSJ9.XOrKFMHSJFYbo-Fq8I5Yd9znJh6prA2t86JX89FrrRGO6r-0n-2T7VTeHwd0TP7loDQQOCeBd2NG-4wkSAnDRtBBQqk0jgQjcdQwpYXBaZonnWYwRFyBI6nYlxo5Iq5DSasEKt7kYJ6PpMF7Pcp8jfauYB4wGPGmFrf_PkpXUZZqftDYqWRGeCaguoPKyUoOGEtUNfyDo7pK5T2RmUdBoxu63qs-Z9ot0ZUzJ0ZxrwbIDDGv4PQH_edj4Wtix4oWP0HKxkALzAPHkQEuUv-H6gYwxSg0qXgLLOoh_Zd8alIwLhMUfjg71rpkLkN6ZRYdrb3s8dCB_FyesI_DjquY0w"
	testExpiry := time.Now().Local().Add(time.Hour)
	testOauthToken := &oauth2.Token{
		AccessToken:  "wjsTCFTERhLa86xYLy4mPZjvx7RSHV9oUwQ8V3zxKMs1MDMxMTZlMS04Yjk1LTRmNDEtOGYwZi1jNzEwOGI4NzY4ZWU",
		TokenType:    "Bearer",
		RefreshToken: "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg",
		Expiry:       testExpiry,
	}
	testOauthToken = testOauthToken.WithExtra(
		map[string]interface{}{
			"id_token": testIDToken,
		},
	)
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth
	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any())
	ss.oauth.EXPECT().Exchange(ss.ctx, ss.testCode, oauth2.SetAuthURLParam("redirect_uri", ss.testCallbackURL)).Return(testOauthToken, nil)
	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any(), gomock.Any(), gomock.Any())
	ss.logger.EXPECT().Debug("sudir CodeExchange:", gomock.Any())
	want := &OAuthResponse{
		IDToken:      testIDToken,
		AccessToken:  testOauthToken.AccessToken,
		RefreshToken: testOauthToken.RefreshToken,
		Expiry:       &testOauthToken.Expiry,
	}

	got, err := c.CodeExchange(ss.ctx, ss.testCode, ss.testCallbackURL)
	ss.NoError(err)
	ss.Equal(want, got)
}

func (ss *sudirSuite) Test_RefreshToken_InvalidGrant() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	refreshToken := "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg"

	ss.logger.EXPECT().Debug("sudir RefreshToken:", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	ss.oauth.EXPECT().TokenSource(ss.ctx, &oauth2.Token{RefreshToken: refreshToken}).Return(ss.tokenSource)
	ss.tokenSource.EXPECT().Token().Return(nil, &oauth2.RetrieveError{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
	})

	wantErr := ErrInvalidGrant

	got, err := c.RefreshToken(ss.ctx, refreshToken)

	ss.Nil(got)
	ss.EqualError(err, wantErr.Error())
}

func (ss *sudirSuite) Test_RefreshToken_Err() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	testErr := errors.New("some test error")
	refreshToken := "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg"

	ss.oauth.EXPECT().TokenSource(ss.ctx, &oauth2.Token{RefreshToken: refreshToken}).Return(ss.tokenSource)
	ss.tokenSource.EXPECT().Token().Return(nil, testErr)

	wantErr := fmt.Errorf("%s: %w", ErrStatusTextServiceError, testErr)

	got, err := c.RefreshToken(ss.ctx, refreshToken)

	ss.Nil(got)
	ss.EqualError(err, wantErr.Error())
}

func (ss *sudirSuite) Test_RefreshToken_Ok() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	refreshToken := "WNjZilQssfRetWT81slSAp-KTJFYWNRjK9yFMg3YGmySKT64TEiKGuBX_kRaTRImpzt98llHXwpCl9C1HyHALg"

	testExpiry := time.Now().Local().Add(time.Hour)
	testToken := &oauth2.Token{
		AccessToken:  "wjsTCFTERhLa86xYLy4mPZjvx7RSHV9oUwQ8V3zxKMs1MDMxMTZlMS04Yjk1LTRmNDEtOGYwZi1jNzEwOGI4NzY4ZWU",
		TokenType:    "Bearer",
		RefreshToken: "OpE_Tvh73mHkh4nBIEPMNQA7vuhPt9yiSBBn_vHBD1LxjZbJ-EvsR1StWzxVRYyYCkmhYLLdhBrayxDKOC7FqA",
		Expiry:       testExpiry,
	}

	want := &OAuthResponse{
		AccessToken:  testToken.AccessToken,
		RefreshToken: testToken.RefreshToken,
		Expiry:       &testToken.Expiry,
	}

	ss.oauth.EXPECT().TokenSource(ss.ctx, &oauth2.Token{RefreshToken: refreshToken}).Return(ss.tokenSource)
	ss.tokenSource.EXPECT().Token().Return(testToken, nil)
	ss.logger.EXPECT().Debug("sudir RefreshToken:", gomock.Any(), gomock.Any(), gomock.Any())
	got, err := c.RefreshToken(ss.ctx, refreshToken)

	ss.NoError(err)
	ss.Equal(want, got)
}

func (ss *sudirSuite) Test_ParseToken_Err() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	wantErr := fmt.Errorf("%s: %s", ErrStatusTextServiceError, "token is malformed: token contains an invalid number of segments")

	got, err := c.ParseToken("")

	ss.Nil(got)
	ss.EqualError(err, wantErr.Error())
}

func (ss *sudirSuite) Test_ParseToken_Ok() {
	c := NewClient(ss.baseUrl, "clientID", "securityPass", ss.logger)
	c.oauth = ss.oauth

	testIDToken := "eyJraWQiOiJkZWZhdWx0IiwiYWxnIjoiUlMyNTYifQ.eyJlbWFpbCI6ImRldmVsb3BlckBpdC5tb3MucnUiLCJsb2dvbm5hbWUiOiJEZXZlbG9wZXIiLCJ1YV9pZCI6IlNIQTI1Nl91TVcwN0xaWWUtcmVvbHhwbjVqMGRkNmhUcm5zT09DYmJuVlNTVVlnZ3VjIiwianRpIjoiSzVIaU1nblNCYW0tSDR6dmhzM3NROHFDRjdaU0pTek9ZTUJWS2pheDJzYyIsImV4cCI6MTY4ODQxNjEwMCwiY2xvdWRHVUlEIjoiM2M1Y2JiMTYtMDExYS0zMTBlLTk3ZTItNTY1NDAwYTI2NTA2IiwiaWF0IjoxNjg4MzgzNzAwLCJhdWQiOlsiY2ZjLXp2LmRpdGNsb3VkLnJ1Il0sImFtciI6WyJwYXNzd29yZCJdLCJpc3MiOiJodHRwczovL3N1ZGlyLXRlc3QubW9zLnJ1IiwiY29tcGFueSI6ItCT0JrQoyDQmNC90YTQvtCz0L7RgNC-0LQiLCJkZXBhcnRtZW50Ijoi0J7RgtC00LXQuyDQvNC-0LHQuNC70YzQvdC-0Lkg0YDQsNC30YDQsNCx0L7RgtC60LgiLCJwb3NpdGlvbiI6ItC_0YDQvtCz0YDQsNC80LzQuNGB0YIiLCJzdWIiOiJEZXZlbG9wZXJAaHEuY29ycC5tb3MucnUiLCJjcmlkIjoiMCIsInNpZCI6IjUxZTEwMzYxLThhZTUtNGYwMS04ZjRmLWM3YjE4ZTA4ODc2ZSJ9.XOrKFMHSJFYbo-Fq8I5Yd9znJh6prA2t86JX89FrrRGO6r-0n-2T7VTeHwd0TP7loDQQOCeBd2NG-4wkSAnDRtBBQqk0jgQjcdQwpYXBaZonnWYwRFyBI6nYlxo5Iq5DSasEKt7kYJ6PpMF7Pcp8jfauYB4wGPGmFrf_PkpXUZZqftDYqWRGeCaguoPKyUoOGEtUNfyDo7pK5T2RmUdBoxu63qs-Z9ot0ZUzJ0ZxrwbIDDGv4PQH_edj4Wtix4oWP0HKxkALzAPHkQEuUv-H6gYwxSg0qXgLLOoh_Zd8alIwLhMUfjg71rpkLkN6ZRYdrb3s8dCB_FyesI_DjquY0w"

	want := &JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://sudir-test.mos.ru",
			Subject:   "Developer@hq.corp.mos.ru",
			Audience:  jwt.ClaimStrings{"cfc-zv.ditcloud.ru"},
			ExpiresAt: jwt.NewNumericDate(time.Date(2023, time.July, 3, 23, 28, 20, 0, time.Local)),
			IssuedAt:  jwt.NewNumericDate(time.Date(2023, time.July, 3, 14, 28, 20, 0, time.Local)),
			ID:        "K5HiMgnSBam-H4zvhs3sQ8qCF7ZSJSzOYMBVKjax2sc",
		},
		UserClaims: UserClaims{
			SID:        "51e10361-8ae5-4f01-8f4f-c7b18e08876e",
			CloudGUID:  "3c5cbb16-011a-310e-97e2-565400a26506",
			Company:    "ГКУ Инфогород",
			Department: "Отдел мобильной разработки",
			Email:      "developer@it.mos.ru",
			LogonName:  "Developer",
			Position:   "программист",
		},
	}

	got, err := c.ParseToken(testIDToken)

	ss.NoError(err)
	ss.Equal(want, got)
}

func TestSync(t *testing.T) {
	suite.Run(t, &sudirSuite{})
}
*/
