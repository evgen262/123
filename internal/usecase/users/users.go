package users

import (
	"context"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"

	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
)

type usersInteractor struct {
	employeesRepository EmployeesRepository
	logger              ditzap.Logger
}

func NewUsersInteractor(employeesRepository EmployeesRepository, logger ditzap.Logger) *usersInteractor {
	return &usersInteractor{
		employeesRepository: employeesRepository,
		logger:              logger,
	}
}

// GetMe
// Метод возвращает следующую информацию о пользователе:
//   - ФИО
//   - Пол
//   - Идентификатор фото
//   - Идентификатор физического лица в 1с
//
// TODO: В настоящее время нет сервиса пользователей
// в связи с этим информацию о пользователе получаем следующим способом:
//  1. Из сервиса сессии забираем идентификатор физического в 1с
//  2. По данному идентификатору и идентификатору портала на котором залогинен пользователь
//     в сервисе employees получаем информацию о физлице + сотрудник
//  3. Если по каким-либо причинам не можем получить информацию - выдаём пустые поля
func (ui *usersInteractor) GetMe(ctx context.Context) (*entityUser.UserInfo, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		ui.logger.Error(usecase.ErrGetSessionFromContext.Error(),
			zap.Error(err),
		)
		return nil, usecase.ErrGetSessionFromContext
	}

	employeeExtID := session.GetUser().GetEmployee().GetExtID()
	if employeeExtID == "" {
		ui.logger.Error(ErrEmptySessionEmployeeExtID.Error(),
			zap.Object("session", session),
		)
		return nil, ErrEmptySessionEmployeeExtID
	}

	activePortal := session.GetActivePortal()
	if activePortal == nil || activePortal.Portal.ID == 0 {
		ui.logger.Error(ErrZeroSessionActivePortalID.Error(),
			zap.Object("session", session),
		)
		return nil, ErrZeroSessionActivePortalID
	}

	userInfo := new(entityUser.UserInfo)
	employee, err := ui.employeesRepository.GetByExtIDAndPortalID(ctx, employeeExtID, activePortal.Portal.ID)
	if err != nil {
		ui.logger.Error("can't get composite employee in repository",
			zap.Object("session", session),
			zap.Error(err),
		)
		return userInfo, nil
	}

	if employee == nil {
		ui.logger.Error("received nil composite employee from repository",
			zap.Object("session", session),
		)
		return userInfo, nil
	}

	userInfo.User = entityUser.ShortUser{
		LastName:   employee.Person.LastName,
		FirstName:  employee.Person.FirstName,
		MiddleName: employee.Person.MiddleName,
		ImageID:    employee.Person.ImageID,
		Gender:     employee.Person.Gender,
		PortalData: entityUser.PortalData{
			PersonID:   employee.Person.ExtID,
			EmployeeID: employeeExtID,
		},
	}

	return userInfo, nil
}
