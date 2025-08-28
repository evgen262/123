package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

type authUseCase struct {
	sudir              SudirClient
	kadry              KadryClient
	stateRepository    StateRepository
	tokenRepository    TokenRepository
	employeeRepository EmployeeRepository
	logger             ditzap.Logger
}

func NewAuthUseCase(
	sudir SudirClient,
	kadry KadryClient,
	stateRepo StateRepository,
	tokenRepo TokenRepository,
	employeeRepo EmployeeRepository,
	logger ditzap.Logger,
) *authUseCase {
	return &authUseCase{
		sudir:              sudir,
		kadry:              kadry,
		stateRepository:    stateRepo,
		tokenRepository:    tokenRepo,
		employeeRepository: employeeRepo,
		logger:             logger,
	}
}

// GetAuthURL
//
//	URL для перенаправления пользователя для авторизации в СУДИР
func (a *authUseCase) GetAuthURL(ctx context.Context, callbackURL, clientID, clientSecret string) (string, error) {

	options := sudir.AuthURLOptions{
		IsOffline:   true,
		RedirectURI: callbackURL,
	}

	// Если переданы client_id и client_secret
	// проводим аутентификацию мобильного приложения
	if clientID != "" && clientSecret != "" {
		cv, err := sudir.NewCodeVerifier("")
		if err != nil {
			return "", err
		}

		stateOptions := &entity.StateOptions{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			CodeVerifier: cv.Value,
		}

		if deviceId, ok := ctx.Value("deviceid").(string); ok {
			stateOptions.DeviceID = deviceId
		}

		if ua, ok := ctx.Value("x-cfc-useragent").(string); ok {
			stateOptions.UserAgent = ua
		}

		state, err := a.stateRepository.New(ctx, stateOptions)
		if err != nil {
			return "", err
		}

		options.State = state.ID
		options.ClientID = state.ClientID
		options.CodeChallengeMethod = "S256"
		options.CodeChallenge = cv.CodeChallengeS256()
	} else {
		stateOptions := &entity.StateOptions{
			CallbackURL: callbackURL,
		}
		state, err := a.stateRepository.New(ctx, stateOptions)
		if err != nil {
			return "", err
		}
		options.State = state.ID
	}

	return a.sudir.AuthURL(options), nil
}

// Auth аутентификация пользователя
//
//	метод возвращает информацию о пользователе в СУДИР
//	и oauth2 токены
func (a *authUseCase) Auth(ctx context.Context, code, stateID, callbackURL string) (*entity.AuthInfo, error) {
	if err := a.stateRepository.IsExists(ctx, stateID); err != nil {
		return nil, err
	}

	state, err := a.stateRepository.Get(ctx, stateID)
	if err != nil {
		return nil, err
	}

	defer a.stateRepository.Delete(context.WithoutCancel(ctx), stateID)

	if state.CallbackURL != "" && state.CallbackURL != callbackURL {
		a.logger.Warn("callback url mismatch",
			entity.LogModuleUE,
			entity.LogCode("UE_018"),
			zap.String("state-id", stateID),
			zap.String("state_callback_url", state.CallbackURL),
			zap.String("callback_url", callbackURL),
		)
		return nil, ErrCallbackURLMismatch
	}

	options := sudir.CodeExchangeOptions{
		RedirectURI:  callbackURL,
		ClientID:     state.ClientID,
		ClientSecret: state.ClientSecret,
		CodeVerifier: state.CodeVerifier,
	}

	if state.ClientID != "" && state.ClientSecret != "" {
		a.logger.Debug("state from cache",
			zap.String("callback_url", callbackURL),
			zap.String("client-id", state.ClientID),
			zap.String("client_secret", state.ClientSecret),
			zap.String("code_verifier", state.CodeVerifier),
			zap.String("user-agent", state.UserAgent))
	}

	oauth, err := a.sudir.CodeExchange(ctx, code, options)
	if err != nil {
		return nil, err
	}

	user, err := a.sudir.ParseToken(oauth.IDToken)
	if err != nil {
		return nil, err
	}

	info := &entity.AuthInfo{
		OAuth: &entity.OAuth{
			AccessToken:  oauth.AccessToken,
			RefreshToken: oauth.RefreshToken,
			Expiry:       oauth.Expiry,
		},
		User: &entity.User{
			CloudID: entity.CloudID(user.CloudGUID),
			Info: &entity.UserInfo{
				SessionID:  user.SID,
				Sub:        user.Subject,
				LastName:   user.FamilyName,
				FirstName:  user.Name,
				MiddleName: user.MiddleName,
				LogonName:  user.LogonName,
				Company:    user.Company,
				Department: user.Department,
				Position:   user.Position,
				Email:      user.Email,
			},
		},
	}

	if state.DeviceID != "" || state.ClientID != "" {
		info.Device = &entity.Device{
			ID:        state.DeviceID,
			ClientID:  state.ClientID,
			UserAgent: state.UserAgent,
		}
	}

	params := entity.EmployeeGetParams{Key: entity.EmployeeGetKeyCloudID, CloudID: user.CloudGUID}
	if user.CloudGUID == "" {
		a.logger.Warn("в payload JWT отсутствует CloudGUID",
			entity.LogModuleUE,
			entity.LogCode("UE_012"),
			zap.String("logon-name", user.LogonName),
			zap.String("email", user.Email),
		)
		params = entity.EmployeeGetParams{Key: entity.EmployeeGetKeyEmail, Email: user.Email}
		if user.Email == "" {
			a.logger.Error("в payload JWT отсутствует Email",
				entity.LogModuleUE,
				entity.LogCode("UE_022"),
				zap.String("logon-name", user.LogonName),
			)
			return info, nil
		}
	}

	// В качестве ключа (tokenID) используется cloudGUID, а при его отсутствии email
	// Если вход осуществляется в WEB, то в качестве ключа используется tokenID + идентификатор сессии СУДИР
	// если через мобильное приложение - то динамический client_id
	tokenID := params.ParamByKey()
	if user.SID != "" {
		tokenID = tokenID + ":" + user.SID
	}
	if state.ClientID != "" {
		tokenID = state.ClientID
	}

	if err = a.tokenRepository.Save(ctx, tokenID, oauth.RefreshToken); err != nil {
		a.logger.Warn("Ошибка при сохранении refresh_token",
			zap.String("id", tokenID),
			zap.Error(err),
		)
	}

	// Получаем сотрудников из кэша
	employees, err := a.employeeRepository.Get(ctx, params.ParamByKey())
	if err == nil && len(employees) > 0 {
		info.User.Employees = employees
		return info, nil
	}

	// TODO подумать над переносом функционала этого кейса в Auth-facade
	var errEmp error
	switch {
	case params.Key == entity.EmployeeGetKeyCloudID:
		var employeesInfo []entity.EmployeeInfo
		employeesInfo, err = a.kadry.GetEmployeesInfo(ctx, user.CloudGUID, kadry.PersonID, kadry.InnOrg, kadry.SNILS)
		if err != nil {
			a.logger.Error("ошибка запроса к серверу СКС",
				entity.LogModuleUE,
				entity.LogCode("UE_023"),
				zap.String("logon-name", user.LogonName),
				zap.String("email", user.Email),
				zap.String("cloud-id", user.CloudGUID),
				zap.Error(err),
			)
		} else if len(employeesInfo) > 0 {
			userEmployees := make([]entity.EmployeeInfo, 0, len(employeesInfo))
			for _, employee := range employeesInfo {
				if user.CloudGUID != employee.CloudID {
					a.logger.Warn("неверный CloudID из системы кадров",
						zap.String("sudir-id", user.CloudGUID),
						zap.String("cloud-id", employee.CloudID),
					)
					continue
				}
				userEmployees = append(userEmployees, employee)
				info.User.Employees = userEmployees
			}
			break
		}

		if user.Email == "" {
			errEmp = kadry.ErrEmployeesService
			break
		}

		fallthrough
	case params.Key == entity.EmployeeGetKeyEmail:
		// Получаем сотрудников из сервиса сотрудников
		var personID uuid.UUID
		if personID, err = a.employeeRepository.GetPersonIDByEmployeeEmail(ctx, user.Email); err != nil {
			a.logger.Error("ошибка получения personID из сервиса сотрудников",
				entity.LogModuleUE,
				entity.LogCode("UE_024"),
				zap.String("logon-name", user.LogonName),
				zap.String("email", user.Email),
				zap.Error(err),
			)
			errEmp = fmt.Errorf("ошибка получения personID из сервиса сотрудников: %w", err)
			break
		}
		if info.User.Employees, err = a.employeeRepository.GetEmployeesInfoByPersonID(ctx, personID); err != nil {
			a.logger.Error("ошибка получения данных о сотрудниках из сервиса сотрудников",
				entity.LogModuleUE,
				entity.LogCode("UE_025"),
				ditzap.UUID("person-id", personID),
				zap.String("logon-name", user.LogonName),
				zap.String("email", user.Email),
				zap.Error(err),
			)
			errEmp = fmt.Errorf("ошибка получения данных о сотрудниках из сервиса сотрудников: %w", err)
		}
	}

	if errEmp != nil {
		return info, errEmp
	}

	if len(info.User.Employees) == 0 {
		logCode := "UE_026"
		msg := "данные о сотрудниках не найдены в СКС"
		if params.Key == entity.EmployeeGetKeyEmail {
			logCode = "UE_027"
			msg = "данные о сотрудниках не найдены в сервисе сотрудников"
		}
		a.logger.Warn(msg,
			entity.LogModuleUE,
			entity.LogCode(logCode),
			zap.String(params.Key.String(), params.ParamByKey()),
		)
		return info, nil
	}

	// Сохраняем сотрудников в кэш
	if err = a.employeeRepository.Save(ctx, params.ParamByKey(), info.User.Employees); err != nil {
		a.logger.Warn("ошибка при сохранении данных о сотрудниках",
			zap.String(params.Key.String(), params.ParamByKey()),
			zap.Error(err),
		)
	}

	return info, nil
}

func (a *authUseCase) LoginByCredentials(ctx context.Context, clientID, clientSecret string) (*entity.AuthInfo, error) {
	opts := sudir.LoginOptions{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	oauth, err := a.sudir.LoginCredentials(ctx, opts)
	if err != nil {
		return nil, err
	}

	info := &entity.AuthInfo{
		OAuth: &entity.OAuth{
			AccessToken: oauth.AccessToken,
			Expiry:      oauth.Expiry,
		},
	}

	return info, nil
}

// RefreshToken обновление токена в СУДИР
//
//	 по идентификатору пользователя cloud_id + идентификатор сессии СУДИР
//	 или по идентификатору инстанса мобильного приложения client_id
//		метод возвращает обновлённый access_token
func (a *authUseCase) RefreshToken(ctx context.Context, id, sessionID string) (string, error) {
	if sessionID != "" {
		id = id + ":" + sessionID
	}

	refreshToken, err := a.tokenRepository.Get(ctx, id)
	if err != nil {
		return "", err
	}

	oauth, err := a.sudir.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	// продолжаем работу так как ошибка не мешает бизнес-логике (отдаём новый access token)
	if err := a.tokenRepository.Save(ctx, id, oauth.RefreshToken); err != nil {
		a.logger.Warn("Ошибка при сохранении refresh_token",
			zap.String("id", id),
			zap.Error(err),
		)
	}

	return oauth.AccessToken, nil
}

// GetUserInfo получение информации о пользователе
//
//	в обмен на access_token
func (a *authUseCase) GetUserInfo(ctx context.Context, accessToken string) (*entity.UserInfo, error) {
	uInfo, err := a.sudir.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return &entity.UserInfo{
		Sub:        uInfo.Sub,
		LastName:   uInfo.FamilyName,
		FirstName:  uInfo.Name,
		MiddleName: uInfo.MiddleName,
		LogonName:  uInfo.LogonName,
		Company:    uInfo.Company,
		Department: uInfo.Department,
		Position:   uInfo.Position,
		Email:      uInfo.Email,
	}, nil
}

// GetEmployees получение информации о пользователе как о сотруднике
//
//	в обмен на cloud_id или email
func (a *authUseCase) GetEmployees(ctx context.Context, params entity.EmployeeGetParams) ([]entity.EmployeeInfo, error) {
	// Получаем сотрудников из кэша
	employees, err := a.employeeRepository.Get(ctx, params.ParamByKey())
	if err == nil && len(employees) > 0 {
		return employees, nil
	}

	var userEmployees []entity.EmployeeInfo
	switch params.Key {
	// Получаем сотрудников из СКС
	case entity.EmployeeGetKeyCloudID:
		employeesInfo, getErr := a.kadry.GetEmployeesInfo(ctx, params.CloudID, kadry.PersonID, kadry.InnOrg, kadry.SNILS)
		if getErr != nil {
			a.logger.Error("ошибка запроса к серверу СКС",
				entity.LogModuleUE,
				entity.LogCode("UE_028"),
				zap.String("cloud-id", params.CloudID),
				zap.Error(getErr),
			)
			return nil, kadry.ErrEmployeesService
		}

		userEmployees = make([]entity.EmployeeInfo, 0, len(employeesInfo))
		for _, employee := range employeesInfo {
			if params.CloudID != employee.CloudID {
				a.logger.Warn("неверный CloudID из системы кадров",
					zap.String("sudir-id", params.CloudID),
					zap.String("cloud-id", employee.CloudID),
				)
				continue
			}
			userEmployees = append(userEmployees, employee)
		}
	// Получаем сотрудников из сервиса сотрудников
	// TODO подумать над переносом функционала этого кейса в Auth-facade
	case entity.EmployeeGetKeyEmail:
		var personID uuid.UUID
		personID, err = a.employeeRepository.GetPersonIDByEmployeeEmail(ctx, params.Email)
		if err != nil {
			a.logger.Error("ошибка получения personID из сервиса сотрудников",
				entity.LogModuleUE,
				entity.LogCode("UE_029"),
				zap.String("email", params.Email),
				zap.Error(err),
			)
			return nil, fmt.Errorf("ошибка получения personID из сервиса сотрудников: %w", err)
		}

		userEmployees, err = a.employeeRepository.GetEmployeesInfoByPersonID(ctx, personID)
		if err != nil {
			a.logger.Error("ошибка получения данных о сотрудниках из сервиса сотрудников",
				entity.LogModuleUE,
				entity.LogCode("UE_030"),
				ditzap.UUID("person-id", personID),
				zap.String("email", params.Email),
				zap.Error(err),
			)
			return nil, fmt.Errorf("ошибка получения данных о сотрудниках из сервиса сотрудников: %w", err)
		}
	default:
		return nil, diterrors.NewValidationError(ErrInvalidKeyType)
	}

	if len(userEmployees) == 0 {
		logCode := "UE_031"
		msg := "данные о сотрудниках не найдены в СКС"
		if params.Key == entity.EmployeeGetKeyEmail {
			logCode = "UE_032"
			msg = "данные о сотрудниках не найдены в сервисе сотрудников"
		}
		a.logger.Warn(msg,
			entity.LogModuleUE,
			entity.LogCode(logCode),
			zap.String(params.Key.String(), params.ParamByKey()),
		)
		return userEmployees, nil
	}

	// Сохраняем сотрудников в кэш
	if err = a.employeeRepository.Save(ctx, params.ParamByKey(), userEmployees); err != nil {
		a.logger.Warn("ошибка при сохранении данных о сотрудниках",
			zap.String(params.Key.String(), params.ParamByKey()),
			zap.Error(err),
		)
	}

	return userEmployees, nil
}

// IsValidToken валидация access_token
func (a *authUseCase) IsValidToken(ctx context.Context, accessToken string) (*entity.TokenInfo, error) {
	vInfo, err := a.sudir.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if !vInfo.IsActive {
		return nil, ErrTokenIsExpire
	}

	scopesInfo := strings.Split(vInfo.Scope, " ")
	scopes := make([]entity.ScopeType, 0, len(scopesInfo))

	for _, scopeInfo := range scopesInfo {
		scope := entity.ScopeUnknown
		switch scopeInfo {
		case sudir.ScopeOpenId:
			scope = entity.ScopeOpenId
		case sudir.ScopeEmail:
			scope = entity.ScopeEmail
		case sudir.ScopeEmployee:
			scope = entity.ScopeEmployee
		case sudir.ScopeGroups:
			scope = entity.ScopeGroups
		case sudir.ScopeProfile:
			scope = entity.ScopeProfile
		case sudir.ScopeUserInfo:
			scope = entity.ScopeUserInfo
		}
		scopes = append(scopes, scope)
	}

	return &entity.TokenInfo{
		Subject:        vInfo.Sub,
		Scopes:         scopes,
		TokenType:      vInfo.TokenType,
		ClientID:       vInfo.ClientID,
		IsActive:       vInfo.IsActive,
		ExpirationTime: time.Unix(vInfo.Exp, 0),
		IssuedAt:       time.Unix(vInfo.Iat, 0),
	}, nil
}

// Logout очистка параметров авторизации пользователя
func (a *authUseCase) Logout(ctx context.Context, id, sessionID, registrationToken string) error {
	if sessionID != "" {
		// для WEB достаточно удалить refresh_token из хранилища
		a.tokenRepository.Delete(ctx, id+":"+sessionID)
		return nil
	}

	a.tokenRepository.Delete(ctx, id)

	// если вход был выполнен с помощью динамических реквизитов - необходимо произвести выход из СУДИР
	if registrationToken != "" {
		if err := a.sudir.Logout(ctx, id, registrationToken); err != nil {
			return err
		}
	}

	return nil
}
