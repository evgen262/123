package presenter

import (
	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

type authPresenter struct {
	timeUtils timeUtils.TimeUtils
}

func NewAuthPresenter() *authPresenter {
	return &authPresenter{}
}

func (ap *authPresenter) UserToPb(user *entity.User) *authv1.User {
	if user == nil {
		return nil
	}

	pb := &authv1.User{
		CloudId:   string(user.CloudID),
		Employees: ap.EmployeesToPb(user.Employees),
		Info:      ap.UserInfoToPb(user.Info),
	}

	// Добавляем идентификатор сессии СУДИР
	if user.Info != nil {
		pb.SudirSid = user.Info.SessionID
	}

	if len(user.Employees) > 0 {
		pb.Snils = user.Employees[0].SNILS
	}

	return pb
}

func (ap *authPresenter) EmployeesToPb(entities []entity.EmployeeInfo) []*authv1.Employee {
	employees := make([]*authv1.Employee, 0, len(entities))
	for _, employee := range entities {
		emp := employee
		employees = append(employees, &authv1.Employee{
			Inn:   emp.Inn,
			OrgId: emp.OrgID,
			Fio:   emp.FIO,
		})
	}
	return employees
}

func (ap *authPresenter) UserInfoToPb(user *entity.UserInfo) *authv1.UserInfo {
	if user == nil {
		return nil
	}

	pb := &authv1.UserInfo{
		Email:           user.Email,
		LogonName:       user.LogonName,
		Sub:             user.Sub,
		LastName:        user.LastName,
		FirstName:       user.FirstName,
		MiddleName:      user.MiddleName,
		SudirCompany:    user.Company,
		SudirDepartment: user.Department,
		SudirPosition:   user.Position,
	}

	return pb
}

func (ap *authPresenter) DeviceToPb(device *entity.Device) *authv1.AuthResponse_Device {
	if device == nil {
		return nil
	}

	return &authv1.AuthResponse_Device{
		Id:        device.ID,
		ClientId:  device.ClientID,
		UserAgent: device.UserAgent,
	}
}

func (ap *authPresenter) TokenInfoToPb(info *entity.TokenInfo) *authv1.TokenInfo {
	if info == nil {
		return nil
	}

	pb := &authv1.TokenInfo{
		Subject:    info.Subject,
		Scopes:     ap.ScopePb(info.Scopes),
		TokenType:  info.TokenType,
		ClientId:   info.ClientID,
		IsActive:   info.IsActive,
		ExpTime:    ap.timeUtils.TimeToTimestamp(&info.ExpirationTime),
		IssuedTime: ap.timeUtils.TimeToTimestamp(&info.IssuedAt),
	}

	return pb
}

func (ap *authPresenter) ScopePb(scopes []entity.ScopeType) []authv1.Scope {
	pb := make([]authv1.Scope, 0, len(scopes))

	for _, scope := range scopes {
		psCope := authv1.Scope_SCOPE_INVALID
		switch scope {
		case entity.ScopeOpenId:
			psCope = authv1.Scope_SCOPE_OPEN_ID
		case entity.ScopeProfile:
			psCope = authv1.Scope_SCOPE_PROFILE
		case entity.ScopeEmail:
			psCope = authv1.Scope_SCOPE_EMAIL
		case entity.ScopeUserInfo:
			psCope = authv1.Scope_SCOPE_USER_INFO
		case entity.ScopeEmployee:
			psCope = authv1.Scope_SCOPE_EMPLOYEE
		case entity.ScopeGroups:
			psCope = authv1.Scope_SCOPE_GROUPS
		}
		pb = append(pb, psCope)
	}
	return pb
}
