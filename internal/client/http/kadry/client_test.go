package kadry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

type kadrySuite struct {
	suite.Suite

	subscriberID string
	userID       string
	secret       string
	baseURL      string

	ctx context.Context

	logger  *ditzap.MockLogger
	handler *http.ServeMux
}

func (ks *kadrySuite) SetupTest() {
	ctrl := gomock.NewController(ks.T())

	ks.logger = ditzap.NewMockLogger(ctrl)

	ks.ctx = context.TODO()

	ks.handler = http.NewServeMux()

	ks.baseURL = "https://test/url"
	ks.subscriberID = "test_subscriber"
	ks.userID = "test_user"
	ks.secret = "security_pass"
}

func (ks *kadrySuite) Test_getServiceUrl_Err() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	wantErr := fmt.Errorf("parse \"%s\\x13\": net/url: invalid control character in URL", ks.baseURL)
	got, err := c.getServiceURL("\u0013")

	ks.Empty(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_getServiceUrl_Ok() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	want := fmt.Sprintf("%s/method/path?SubscriberID=%s&UserID=%s", ks.baseURL, ks.subscriberID, ks.userID)
	got, err := c.getServiceURL("/method/path")

	ks.NoError(err)
	ks.Equal(want, got)
}

func (ks *kadrySuite) Test_newHTTPRequest_ErrURL() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	wantErr := fmt.Errorf("parse \"%s\\x12\": net/url: invalid control character in URL", ks.baseURL)
	got, err := c.newHTTPRequest(http.MethodGet, "\u0012", nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_newHTTPRequest_ErrMethod() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	wantErr := fmt.Errorf("net/http: invalid method \",\"")
	got, err := c.newHTTPRequest(",", "/method/path", nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_newHTTPRequest_ErrData() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	wantErr := fmt.Errorf("json: unsupported type: func()")
	got, err := c.newHTTPRequest(http.MethodGet, "/method/path", func() {})

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_newHTTPRequest_Ok() {
	c := NewClient(ks.baseURL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	want, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/method/path?SubscriberID=%s&UserID=%s", ks.baseURL, ks.subscriberID, ks.userID), nil)
	want.Header.Set("Content-Type", "application/json")
	want.Header.Set("Authorization", "Basic dGVzdF9zdWJzY3JpYmVyOnNlY3VyaXR5X3Bhc3M=")

	got, err := c.newHTTPRequest(http.MethodGet, "/method/path", nil)

	ks.NoError(err)
	ks.Equal(want, got)
}

func (ks *kadrySuite) Test_sksMobileApp_ErrData() {
	c := NewClient("\u0011", ks.subscriberID, ks.secret, ks.userID, ks.logger)

	// methodErr := errors.New("parse \"\\x11/hs/frontend_api/execute/sksMobileApp\": net/url: invalid control character in URL")
	ks.logger.EXPECT().Error("kadry sksMobileApp: ошибка при создании запроса",
		entity.LogModuleUE,
		entity.LogCode("UE_020"),
		gomock.Any(),
		gomock.Any(),
	)

	got, err := c.sksMobileApp(ks.ctx, nil)

	ks.Nil(got)
	ks.EqualError(err, ErrStatusTextBadRequest)
}

func (ks *kadrySuite) Test_sksMobileApp_ServerError() {
	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	ks.logger.EXPECT().Error("неизвестный ответ от сервера СКС",
		entity.LogModuleUE,
		entity.LogCode("UE_021"),
		gomock.Any(),
		zap.String("status", "500 Internal Server Error"),
	)
	wantErr := fmt.Errorf("%s код: %d", ErrStatusTextServiceError, 500)
	got, err := c.sksMobileApp(ks.ctx, nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_sksMobileApp_BadRequest() {
	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)
	ks.logger.EXPECT().Debug("kadry sksMobileApp: ошибка десериализации",
		entity.LogModuleUE,
		entity.LogCode("UE_021"),
		gomock.Any(),
		gomock.Any(),
	)
	wantErr := fmt.Errorf("%s: %s", ErrStatusTextServiceError, "EOF")
	got, err := c.sksMobileApp(ks.ctx, nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_sksMobileApp_BadBody() {
	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ks.Require().NoError(json.NewEncoder(w).Encode("<html></html>"))
	})
	s := httptest.NewServer(handler)

	ks.logger.EXPECT().Debug("kadry sksMobileApp: ошибка десериализации",
		entity.LogModuleUE,
		entity.LogCode("UE_021"),
		gomock.Any(),
		gomock.Any(),
	)
	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	wantErr := fmt.Errorf("%s: %s", ErrStatusTextServiceError, "json: cannot unmarshal string into Go value of type kadry.response")
	got, err := c.sksMobileApp(ks.ctx, nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_sksMobileApp_ReqErr() {
	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)
	ks.logger.EXPECT().Error("kadry sksMobileApp: ошибка запроса к серверу СКС",
		entity.LogModuleUE,
		entity.LogCode("UE_021"),
		gomock.Any(),
		gomock.Any(),
	)
	wantErr := fmt.Errorf(
		"%s: Post \"%s/hs/frontend_api/execute/sksMobileApp?SubscriberID=%s&UserID=%s\": context deadline exceeded",
		ErrStatusTextErrorRequest,
		s.URL,
		ks.subscriberID,
		ks.userID,
	)
	ctx, cancel := context.WithTimeout(ks.ctx, time.Nanosecond)
	defer cancel()
	got, err := c.sksMobileApp(ctx, nil)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_sksMobileApp_Ok() {
	testResponse := &response{
		commonResponse: commonResponse{
			MessageType:     "Result",
			RequestExecuted: false,
			RequestType:     "sksMobileApp",
		},
	}

	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ks.Require().NoError(json.NewEncoder(w).Encode(testResponse))
	})
	s := httptest.NewServer(handler)
	ks.logger.EXPECT().Debug("kadry sksMobileApp:", gomock.Any(), gomock.Any())
	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	got, err := c.sksMobileApp(ks.ctx, nil)

	ks.Equal(got, testResponse)
	ks.NoError(err)
}

func (ks *kadrySuite) Test_GetEmployeesInfo_Err() {
	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	ks.logger.EXPECT().Error("неизвестный ответ от сервера СКС",
		entity.LogModuleUE,
		entity.LogCode("UE_021"),
		zap.Object("request", &request{PersonIDArray: []string{"3c5cbb16-011a-310e-97e2-565400a26506"}}),
		zap.String("status", "500 Internal Server Error"),
	)
	wantErr := fmt.Errorf("%s код: %d", ErrStatusTextServiceError, 500)
	got, err := c.GetEmployeesInfo(ks.ctx, "3c5cbb16-011a-310e-97e2-565400a26506")

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_GetEmployeesInfo_BadResponse() {
	testCloudID := "3c5cbb16-011a-310e-97e2-565400a26506"
	testRequest := &request{
		PersonIDArray: []string{"3c5cbb16-011a-310e-97e2-565400a26506"},
	}
	testResponse := &response{
		commonResponse: commonResponse{
			MessageType:      "Result",
			RequestExecuted:  false,
			RequestType:      "sksMobileApp",
			ErrorDescription: "some bad parameter in request",
		},
	}

	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		ks.Require().NoError(json.NewEncoder(w).Encode(testResponse))
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	ks.logger.EXPECT().Debug("kadry sksMobileApp:", zap.Object("request", testRequest), zap.Object("response", testResponse))
	wantErr := fmt.Errorf("%s %s", ErrStatusTextServiceError, testResponse.ErrorDescription)
	got, err := c.GetEmployeesInfo(ks.ctx, testCloudID)

	ks.Nil(got)
	ks.EqualError(err, wantErr.Error())
}

func (ks *kadrySuite) Test_GetEmployeesInfo_Ok() {
	testCloudID := "3c5cbb16-011a-310e-97e2-565400a26506"
	testResponse := &response{
		commonResponse: commonResponse{
			MessageType:     "Result",
			RequestExecuted: true,
			RequestType:     "sksMobileApp",
		},
		Body: struct {
			MobileApp mobileApp `json:"MobileApp,omitempty"`
		}{
			MobileApp: []mobileAppInfo{
				{
					PersonID:  "3c5cbb16-011a-310e-97e2-565400a26506",
					InnOrg:    "770123456789",
					OrgID:     "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					FIOPerson: "Иванов Иван Иванович",
					SNILS:     "123-123-123 23",
				},
				{
					PersonID:  "3c5cbb16-011a-310e-97e2-565400a26506",
					InnOrg:    "779876543210",
					OrgID:     "71bac977-ea7f-4156-9504-60f7d443ab62",
					FIOPerson: "Иванов Иван Иванович",
					SNILS:     "123-123-123 23",
				},
			},
		},
	}

	want := []entity.EmployeeInfo{
		{
			CloudID: "3c5cbb16-011a-310e-97e2-565400a26506",
			Inn:     "770123456789",
			OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			FIO:     "Иванов Иван Иванович",
			SNILS:   "123-123-123 23",
		},
		{
			CloudID: "3c5cbb16-011a-310e-97e2-565400a26506",
			Inn:     "779876543210",
			OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
			FIO:     "Иванов Иван Иванович",
			SNILS:   "123-123-123 23",
		},
	}

	handler := ks.handler
	handler.HandleFunc("/hs/frontend_api/execute/sksMobileApp", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ks.Require().NoError(json.NewEncoder(w).Encode(testResponse))
	})
	s := httptest.NewServer(handler)

	c := NewClient(s.URL, ks.subscriberID, ks.secret, ks.userID, ks.logger)

	ks.logger.EXPECT().Debug("kadry sksMobileApp:", gomock.Any(), gomock.Any())
	got, err := c.GetEmployeesInfo(ks.ctx, testCloudID)

	ks.NoError(err)
	ks.Equal(want, got)
}

func TestSync(t *testing.T) {
	suite.Run(t, &kadrySuite{})
}
