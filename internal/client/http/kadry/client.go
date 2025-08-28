package kadry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

const bodyMaxSize = 1024000

type client struct {
	subscriberID string
	userID       string
	secret       string
	baseURL      string

	logger     ditzap.Logger
	httpClient *http.Client
}

func NewClient(baseURL, subscriberID, secret, userID string, logger ditzap.Logger) *client {
	return &client{
		subscriberID: subscriberID,
		userID:       userID,
		secret:       secret,
		baseURL:      baseURL,

		logger:     logger,
		httpClient: http.DefaultClient,
	}
}

func (c *client) getServiceURL(path string) (string, error) {
	urlPath, err := url.Parse(c.baseURL + path)
	if err != nil {
		return "", err
	}

	params := &url.Values{}
	params.Add("SubscriberID", c.subscriberID)
	params.Add("UserID", c.userID)
	urlPath.RawQuery = params.Encode()

	return urlPath.String(), nil
}

func (c *client) newHTTPRequest(method, methodURL string, data any) (*http.Request, error) {
	srvURL, err := c.getServiceURL(methodURL)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if data != nil {
		jsonData, errJSON := json.Marshal(data)
		if errJSON != nil {
			return nil, errJSON
		}
		body := bytes.NewReader(jsonData)
		req, err = http.NewRequest(method, srvURL, body)
	} else {
		req, err = http.NewRequest(method, srvURL, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(url.QueryEscape(c.subscriberID), url.QueryEscape(c.secret))

	return req, nil
}

func (c *client) sksMobileApp(ctx context.Context, data *request) (*response, error) {
	httpReq, err := c.newHTTPRequest(http.MethodPost, "/hs/frontend_api/execute/sksMobileApp", data)
	if err != nil {
		c.logger.Error("kadry sksMobileApp: ошибка при создании запроса",
			entity.LogModuleUE,
			entity.LogCode("UE_020"),
			zap.Object("request", data),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%s", ErrStatusTextBadRequest)
	}
	defer httpReq.Body.Close()

	httpResp, err := c.httpClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		c.logger.Error("kadry sksMobileApp: ошибка запроса к серверу СКС",
			entity.LogModuleUE,
			entity.LogCode("UE_021"),
			zap.Object("request", data),
			zap.Error(err),
		)
		return nil, fmt.Errorf("%s: %w", ErrStatusTextErrorRequest, err)
	}
	defer httpResp.Body.Close()

	switch httpResp.StatusCode {
	case http.StatusBadRequest:
		fallthrough
	case http.StatusOK:
		resp := &response{}
		if err := json.NewDecoder(io.LimitReader(httpResp.Body, bodyMaxSize)).Decode(resp); err != nil {
			c.logger.Debug("kadry sksMobileApp: ошибка десериализации",
				entity.LogModuleUE,
				entity.LogCode("UE_021"),
				zap.Object("request", data),
				zap.Error(err),
			)
			return nil, fmt.Errorf("%s: %w", ErrStatusTextServiceError, err)
		}

		c.logger.Debug("kadry sksMobileApp:", zap.Object("request", data), zap.Object("response", resp))
		return resp, nil
	default:
		c.logger.Error("неизвестный ответ от сервера СКС",
			entity.LogModuleUE,
			entity.LogCode("UE_021"),
			zap.Object("request", data),
			zap.String("status", httpResp.Status),
		)
	}

	return nil, fmt.Errorf("%s код: %d", ErrStatusTextServiceError, httpResp.StatusCode)
}

func (c *client) GetEmployeesInfo(ctx context.Context, cloudId string, attributes ...AttributeName) ([]entity.EmployeeInfo, error) {
	req := newRequest([]string{cloudId}, attributes...)

	resp, err := c.sksMobileApp(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.RequestExecuted {
		return nil, fmt.Errorf("%s %s", ErrStatusTextServiceError, resp.ErrorDescription)
	}

	employees := make([]entity.EmployeeInfo, 0, len(resp.Body.MobileApp))

	for _, data := range resp.Body.MobileApp {
		employees = append(employees, entity.EmployeeInfo{
			CloudID: data.PersonID,
			Inn:     data.InnOrg,
			OrgID:   data.OrgID,
			FIO:     data.FIOPerson,
			SNILS:   data.SNILS,
		})
	}

	return employees, nil
}
