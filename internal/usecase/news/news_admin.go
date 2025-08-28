package news

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"go.uber.org/zap"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

func NewNewsAdminInteractor(newsRepository NewsRepository, employeeRepository EmployeesRepository, logger ditzap.Logger) *newsAdminInteractor {
	return &newsAdminInteractor{
		newsRepository:     newsRepository,
		employeeRepository: employeeRepository,
		logger:             logger,
	}
}

type newsAdminInteractor struct {
	newsRepository     NewsRepository
	employeeRepository EmployeesRepository
	logger             ditzap.Logger
}

func (i *newsAdminInteractor) Create(ctx context.Context, news *dtoNews.NewNews) (uuid.UUID, error) {
	logger := ditzap.WithFields(i.logger, zap.Any("news", news))
	err := news.Validate()
	if err != nil {
		logger.Debug("newsAdminInteractor.Create: invalid request", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("newsAdminInteractor.Create: invalid request: %w", err)
	}

	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		return uuid.UUID{}, fmt.Errorf("newsAdminInteractor.Create: can't get session: %w", err)
	}

	news.Visibility = &entityNews.NewsVisibility{}
	/**
	Получаем ID портала из сессии, т.к. получения Категории новостей еще нет (https://oblako.mos.ru/jira/browse/TECH-591).
	Портал является видимостью новости, а на первом этапе видимость берется от Категории.
	*/
	if userPortal := session.ActivePortal.GetPortalID(); userPortal != 0 {
		news.Visibility.PortalsIDs = []int{userPortal}
	}

	employee, err := i.employeeRepository.GetByExtIDAndPortalID(ctx, session.GetUser().GetEmployee().GetExtID(), session.ActivePortal.GetPortalID())
	if err != nil {
		logger.Error("newsAdminInteractor.Create: can't get author", zap.Error(err))
		return uuid.UUID{}, fmt.Errorf("newsAdminInteractor.Create: can't get author: %w", err)
	}
	news.Author.ID = employee.ID
	news.Author.LastName = employee.Person.LastName
	news.Author.FirstName = employee.Person.FirstName
	if employee.Person.MiddleName != "" {
		news.Author.MiddleName = &employee.Person.MiddleName
	}
	if employee.Person.ImageID != "" {
		imageID, err := uuid.Parse(employee.Person.ImageID)
		if err != nil {
			logger.Warn("newsAdminInteractor.Create: can't parse author imageID", zap.String("author_imageID", employee.Person.ImageID))
		} else {
			news.Author.ImageID = &imageID
		}
	}

	createdID, err := i.newsRepository.Create(ctx, news)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.AlreadyExistsError)):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.Create: can't create news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.Create: can't create news", zap.Error(err))
		}
		return uuid.UUID{}, fmt.Errorf("newsAdminInteractor.Create: can't create news: %w", err)
	}

	return createdID, nil
}

func (i *newsAdminInteractor) Update(ctx context.Context, id uuid.UUID, news *dtoNews.UpdateNews) (*entityNews.News, error) {
	logger := ditzap.WithFields(i.logger, ditzap.UUID("news_id", id), zap.Any("news", news))
	if id == uuid.Nil {
		return nil, fmt.Errorf("newsAdminInteractor.Update: invalid request: empty ID")
	}

	err := news.Validate()
	if err != nil {
		logger.Debug("newsAdminInteractor.Update: invalid request", zap.Error(err))
		return nil, fmt.Errorf("newsAdminInteractor.Update: invalid request: %w", err)
	}

	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		return nil, fmt.Errorf("newsAdminInteractor.Update: can't get session: %w", err)
	}

	news.Visibility = &entityNews.NewsVisibility{}
	/**
	Получаем ID портала из сессии, т.к. получения Категории новостей еще нет (https://oblako.mos.ru/jira/browse/TECH-591).
	Портал является видимостью новости, а на первом этапе видимость берется от Категории.
	*/
	if userPortal := session.ActivePortal.GetPortalID(); userPortal != 0 {
		news.Visibility.PortalsIDs = []int{userPortal}
	}

	// Проверяем, если фронт передал null в publicationDate, то устанавливаем дату публикации в текущее время publishDate,
	// то устанавливаем дату публикации в "нулевое" время для очистки даты публикации.
	if news.GetPublicationAtPtr() == nil {
		news.PublicationAt = &time.Time{}
	}

	updatedNews, err := i.newsRepository.Update(ctx, id, news)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.AlreadyExistsError)):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.Update: can't update news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.Update: can't update news", zap.Error(err))
		}
		return nil, fmt.Errorf("newsAdminInteractor.Update: can't update news: %w", err)
	}

	return updatedNews, nil
}

func (i *newsAdminInteractor) UpdateFlags(ctx context.Context, id uuid.UUID, flags *dtoNews.UpdateFlags) (*entityNews.News, error) {
	logger := ditzap.WithFields(i.logger, ditzap.UUID("news_id", id), zap.Any("flags", flags))
	if id == uuid.Nil {
		return nil, fmt.Errorf("newsAdminInteractor.UpdateFlags: invalid request: empty ID")
	}
	if flags == nil {
		return nil, fmt.Errorf("newsAdminInteractor.UpdateFlags: invalid request: empty flags")
	}
	news := &dtoNews.UpdateNews{
		OnMain:    flags.OnMain,
		Pinned:    flags.Pinned,
		UpdatedAt: flags.UpdatedAt,
	}
	updatedNews, err := i.newsRepository.Update(ctx, id, news)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.UpdateFlags: can't update news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.UpdateFlags: can't update news", zap.Error(err))
		}
		return nil, fmt.Errorf("newsAdminInteractor.UpdateFlags: can't update news: %w", err)
	}
	return updatedNews, nil
}

func (i *newsAdminInteractor) ChangeStatus(ctx context.Context, id uuid.UUID, status entityNews.NewsStatus) (*entityNews.News, error) {
	logger := ditzap.WithFields(i.logger, ditzap.UUID("news_id", id), zap.Int("status", int(status)))

	news := &dtoNews.UpdateNews{
		Status: status,
	}

	updatedNews, err := i.newsRepository.Update(ctx, id, news)
	if err != nil {
		logger.Error("newsAdminInteractor.ChangeStatus: can't update news", zap.Error(err))
		return nil, fmt.Errorf("newsAdminInteractor.ChangeStatus: can't update news: %w", err)
	}

	return updatedNews, nil
}

func (i *newsAdminInteractor) Get(ctx context.Context, id uuid.UUID) (*entityNews.NewsFull, error) {
	logger := ditzap.WithFields(i.logger, ditzap.UUID("news_id", id))
	result, err := i.newsRepository.Get(ctx, id)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.Get: can't get news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.Get: can't get news", zap.Error(err))
		}
		return nil, fmt.Errorf("newsAdminInteractor.Get: can't get news: %w", err)
	}
	return result, nil
}

func (i *newsAdminInteractor) Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error) {
	logger := ditzap.WithFields(i.logger, zap.Any("search", search))
	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		logger.Debug("newsInteractor.Update: can't get session", zap.Error(err))
		return nil, fmt.Errorf("newsInteractor.Update: can't get session: %w", err)
	}

	search.Visitor = &entityNews.Visitor{
		PortalID: session.ActivePortal.GetPortalID(),
	}

	result, err := i.newsRepository.Search(ctx, search)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.Search: can't search news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.Search: can't search news", zap.Error(err))
		}
		return nil, fmt.Errorf("newsAdminInteractor.Search: can't search news: %w", err)
	}
	return result, nil
}

func (i *newsAdminInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	logger := ditzap.WithFields(i.logger, ditzap.UUID("news_id", id))
	if id == uuid.Nil {
		return fmt.Errorf("newsAdminInteractor.Delete: invalid request: empty ID")
	}

	if err := i.newsRepository.Delete(ctx, id); err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			logger.Debug("newsAdminInteractor.Delete: can't delete news", zap.Error(err))
		default:
			logger.Error("newsAdminInteractor.Delete: can't delete news", zap.Error(err))
		}
		return fmt.Errorf("newsAdminInteractor.Delete: can't delete news: %w", err)
	}
	return nil
}
