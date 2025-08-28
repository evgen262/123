package auth

import (
	"net"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type authMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewAuthMapper(timeUtils timeUtils.TimeUtils) *authMapper {
	return &authMapper{
		timeUtils: timeUtils,
	}
}

func (am authMapper) UserToEntity(userPb *authv1.User) *entityAuth.UserSudir {
	return &entityAuth.UserSudir{
		CloudID:    userPb.GetCloudId(),
		Login:      userPb.GetInfo().GetLogonName(),
		Email:      userPb.GetInfo().GetEmail(),
		FIO:        userPb.GetInfo().GetFio(),
		LastName:   userPb.GetInfo().GetLastName(),
		FirstName:  userPb.GetInfo().GetFirstName(),
		MiddleName: userPb.GetInfo().GetMiddleName(),
		SNILS:      userPb.GetSnils(),
		SID:        userPb.GetSudirSid(),
		Portals:    am.PortalsToEntity(userPb.GetPortals()),
		Employees:  am.EmployeesToEntity(userPb.GetEmployees()),
	}
}

func (am authMapper) PortalsToEntity(portalsPb []*authv1.Portal) []*entityAuth.Portal {
	portals := make([]*entityAuth.Portal, 0, len(portalsPb))
	for _, portalPb := range portalsPb {
		portals = append(portals, &entityAuth.Portal{
			ID:    int(portalPb.GetId()),
			Name:  portalPb.GetShortName(), // portalPb.GetFullName()
			URL:   portalPb.GetUrl(),
			Image: portalPb.GetLogoUrl(),
		})
	}
	return portals
}

func (am authMapper) EmployeesToEntity(employeesPb []*authv1.Employee) []*entityAuth.EmployeeInfo {
	employees := make([]*entityAuth.EmployeeInfo, 0, len(employeesPb))
	for _, employeePb := range employeesPb {
		employees = append(employees, &entityAuth.EmployeeInfo{
			Inn:   employeePb.GetInn(),
			OrgID: employeePb.GetOrgId(),
		})
	}
	return employees
}

func (am authMapper) SessionToEntity(session *authv1.Session) *entityAuth.Session {
	if session == nil {
		return nil
	}

	sessionID := entityAuth.SessionID(am.idToUUID(session.GetId()))
	return &entityAuth.Session{
		ID:                 &sessionID,
		User:               am.userToEntity(session.GetUser()),
		UserAuthType:       am.UserTypeToEntity(session.GetUserType()),
		UserIP:             net.ParseIP(session.GetUserIp()),
		Device:             am.DeviceToEntity(session.GetDevice(), session.GetSudirInfo()),
		ActivePortal:       am.ActivePortalToEntity(session.GetUser().GetPortal()),
		Issuer:             session.GetIssuer(),
		Subject:            session.GetSubject(),
		LastActiveTime:     am.timeUtils.TimestampToTime(session.GetLastActiveTime()),
		AccessExpiredTime:  session.GetAccessExpiredTime().AsTime(),
		RefreshExpiredTime: session.GetRefreshExpiredTime().AsTime(),
		CreatedTime:        session.GetCreatedTime().AsTime(),
		RefreshedTime:      session.GetRefreshedTime().AsTime(),
		IsActive:           session.GetIsActive(),
	}
}

// idToUUID преобразует идентификатор в виде строки в uuid.UUID
func (am authMapper) idToUUID(id string) uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil
	}

	return uid
}

// userToEntity преобразует данные пользователя в сущность User
func (am authMapper) userToEntity(userPb *authv1.Session_User) *entityAuth.User {
	if userPb == nil {
		return nil
	}

	return &entityAuth.User{
		ID:       am.idToUUID(userPb.GetId()),
		CloudID:  userPb.GetCloudId(),
		Login:    userPb.GetLogonName(),
		Email:    userPb.GetEmail(),
		SNILS:    userPb.GetSnils(),
		Employee: am.employeeToEntity(userPb.GetEmployee()),
		Person:   am.personToEntity(userPb.GetPerson()),
	}
}

// employeeToEntity преобразует данные сотрудника в сущность Employee
func (am authMapper) employeeToEntity(employeePb *authv1.Session_User_Employee) *entityAuth.Employee {
	if employeePb == nil {
		return nil
	}

	return &entityAuth.Employee{
		ExtID: employeePb.GetId(),
	}
}

// personToEntity преобразует данные физического лица в сущность Person
func (am authMapper) personToEntity(personPb *authv1.Session_User_Person) *entityAuth.Person {
	if personPb == nil {
		return nil
	}

	return &entityAuth.Person{
		ExtID: personPb.GetId(),
	}
}

// DeviceToEntity преобразует данные устройства пользователя в сущность Device
func (am authMapper) DeviceToEntity(devicePb *authv1.Device, sudirInfo *authv1.SudirInfo) *entityAuth.Device {
	if devicePb == nil && sudirInfo == nil {
		return nil
	}

	device := new(entityAuth.Device)
	if devicePb != nil {
		device.UserAgent = devicePb.GetUserAgent()
	}
	if sudirInfo != nil {
		device.SudirInfo = &entityAuth.SudirInfo{
			SID:      sudirInfo.GetSid(),
			ClientID: sudirInfo.GetClientId(),
		}
	}

	return device
}

// ActivePortalToEntity преобразует активный портал пользователя в сущность ActivePortal
func (am authMapper) ActivePortalToEntity(portal *authv1.ActivePortal) *entityAuth.ActivePortal {
	if portal == nil {
		return nil
	}

	return &entityAuth.ActivePortal{
		Portal: entityAuth.Portal{
			ID:         int(portal.GetId()),
			Name:       portal.GetName(),
			URL:        portal.GetUrl(),
			IsSelected: true,
		},
		SID: portal.GetSid(),
	}
}

// UserTypeToEntity преобразует тип пользователя Pb в сущность UserAuthType
func (am authMapper) UserTypeToEntity(userType authv1.UserType) entityAuth.UserAuthType {
	switch userType {
	case authv1.UserType_USER_TYPE_ANON:
		return entityAuth.UserAuthTypeAnon
	case authv1.UserType_USER_TYPE_AUTH:
		return entityAuth.UserAuthTypeAuth
	case authv1.UserType_USER_TYPE_OLD_AUTH:
		return entityAuth.UserAuthTypeOldAuth
	case authv1.UserType_USER_TYPE_SERVICE:
		return entityAuth.UserAuthTypeService
	default:
		return entityAuth.UserAuthTypeInvalid
	}
}

func (am authMapper) SessionToPb(session *entityAuth.Session) *authv1.Session {
	if session == nil {
		return nil
	}

	return &authv1.Session{
		Id:                 session.GetID().String(),
		User:               am.UserToPb(session.GetUser(), session.GetActivePortal()),
		UserType:           am.UserAuthTypeToPb(session.UserAuthType),
		UserIp:             session.UserIP.String(),
		Device:             am.DeviceToPb(session.GetDevice()),
		SudirInfo:          am.SudirInfoToPb(session.GetDevice().GetSudirInfo()),
		Issuer:             session.Issuer,
		Subject:            session.Subject,
		LastActiveTime:     am.timeUtils.TimeToTimestamp(session.LastActiveTime),
		AccessExpiredTime:  timestamppb.New(session.AccessExpiredTime),
		RefreshExpiredTime: timestamppb.New(session.RefreshExpiredTime),
		CreatedTime:        timestamppb.New(session.CreatedTime),
		RefreshedTime:      timestamppb.New(session.RefreshedTime),
		IsActive:           session.IsActive,
	}
}

func (am authMapper) UserToPb(user *entityAuth.User, activePortal *entityAuth.ActivePortal) *authv1.Session_User {
	if user == nil {
		return nil
	}

	userPb := &authv1.Session_User{
		Id:        user.ID.String(),
		CloudId:   user.CloudID,
		LogonName: user.Login,
		Email:     user.Email,
		Snils:     user.SNILS,
		Portal:    am.ActivePortalToPb(activePortal),
		Employee:  am.EmployeeToPb(user.GetEmployee()),
		Person:    am.PersonToPb(user.GetPerson()),
	}

	return userPb
}

func (am authMapper) EmployeeToPb(employee *entityAuth.Employee) *authv1.Session_User_Employee {
	if employee == nil {
		return nil
	}

	return &authv1.Session_User_Employee{
		Id: employee.ExtID,
	}
}

func (am authMapper) PersonToPb(person *entityAuth.Person) *authv1.Session_User_Person {
	if person == nil {
		return nil
	}

	return &authv1.Session_User_Person{
		Id: person.ExtID,
	}
}

func (am authMapper) ActivePortalToPb(activePortal *entityAuth.ActivePortal) *authv1.ActivePortal {
	if activePortal == nil {
		return nil
	}

	return &authv1.ActivePortal{
		Id:   int32(activePortal.Portal.ID),
		Name: activePortal.Portal.Name,
		Url:  activePortal.Portal.URL,
		Sid:  activePortal.SID,
	}
}

func (am authMapper) DeviceToPb(device *entityAuth.Device) *authv1.Device {
	if device == nil {
		return nil
	}

	return &authv1.Device{
		UserAgent: device.UserAgent,
	}
}

func (am authMapper) SudirInfoToPb(info *entityAuth.SudirInfo) *authv1.SudirInfo {
	if info == nil {
		return nil
	}

	return &authv1.SudirInfo{
		Sid:      info.SID,
		ClientId: info.ClientID,
	}
}

func (am authMapper) UserAuthTypeToPb(userType entityAuth.UserAuthType) authv1.UserType {
	switch userType {
	case entityAuth.UserAuthTypeAnon:
		return authv1.UserType_USER_TYPE_ANON
	case entityAuth.UserAuthTypeAuth:
		return authv1.UserType_USER_TYPE_AUTH
	case entityAuth.UserAuthTypeOldAuth:
		return authv1.UserType_USER_TYPE_OLD_AUTH
	case entityAuth.UserAuthTypeService:
		return authv1.UserType_USER_TYPE_SERVICE
	default:
		return authv1.UserType_USER_TYPE_INVALID
	}
}
