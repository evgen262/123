package news

import (
	"context"
	"errors"
	"fmt"

	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

type commentInteractor struct {
	commentsRepository CommentsRepository
	employeeRepository EmployeesRepository
	logger             ditzap.Logger
}

func NewCommentInteractor(
	commentsRepository CommentsRepository,
	employeeRepository EmployeesRepository,
	logger ditzap.Logger,
) *commentInteractor {
	return &commentInteractor{
		commentsRepository: commentsRepository,
		employeeRepository: employeeRepository,
		logger:             logger,
	}
}

func (i *commentInteractor) Create(ctx context.Context, in dtoNews.NewComment) (uuid.UUID, int, error) {
	logger := ditzap.WithFields(i.logger, zap.Any("comment", in))

	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		logger.Debug("commentInteractor.Create: no session", zap.Error(err))
		return uuid.Nil, 0, fmt.Errorf("commentInteractor.Create: can't get session: %w", err)
	}

	employee, err := i.employeeRepository.GetByExtIDAndPortalID(ctx, session.GetUser().GetEmployee().GetExtID(), session.ActivePortal.GetPortalID())
	if err != nil {
		logger.Error("commentInteractor.Create: can't get author", zap.Error(err))
		return uuid.Nil, 0, fmt.Errorf("commentInteractor.Create: can't get author: %w", err)
	}

	in.AuthorID = &employee.ID

	// 2) Валидация входа (после того как установили AuthorID)
	if errAuthor := in.Validate(); errAuthor != nil {
		logger.Debug("commentInteractor.Create: invalid request", zap.Error(errAuthor))
		return uuid.Nil, 0, fmt.Errorf("commentInteractor.Create: invalid request: %w", errAuthor)
	}

	// 3) Создание комментария в репозитории
	createdID, count, err := i.commentsRepository.Create(ctx, in)
	if err != nil {
		var (
			alreadyExists diterrors.AlreadyExistsError
			validationErr diterrors.ValidationError
		)
		// TODO- обогатить логи createdId и ExtId
		switch {
		case errors.As(err, &alreadyExists):
			logger.Warn("commentInteractor.Create: already exists", zap.Error(err))
		case errors.As(err, &validationErr):
			logger.Error("commentInteractor.Create: validation failed", zap.Error(err))
		default:
			logger.Error("commentInteractor.Create: repo error", zap.Error(err))
		}

		return uuid.Nil, 0, fmt.Errorf("commentInteractor.Create: can't create comment: %w", err)
	}

	return createdID, count, nil
}

func (i *commentInteractor) List(ctx context.Context, params *dtoNews.FilterComments) ([]*entityNews.NewsComment, int, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		i.logger.Debug("commentInteractor.List: can't get session", zap.Error(err))
		return nil, 0, fmt.Errorf("commentInteractor.List: can't get session: %w", err)
	}

	params.Visitor = &entityNews.Visitor{
		PortalID: session.ActivePortal.GetPortalID(),
	}

	comments, total, err := i.commentsRepository.List(ctx, params)
	if err != nil {
		i.logger.Error("commentInteractor.Get: can't get comments", zap.Error(err), zap.String("news_id", params.NewsID.String()))
		return nil, 0, fmt.Errorf("commentInteractor.Get: can't get comments: %w", err)
	}

	return comments, total, nil
}
