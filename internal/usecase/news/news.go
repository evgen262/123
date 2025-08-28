package news

import (
	"context"
	"fmt"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
)

func NewNewsInteractor(newsRepository NewsRepository, employeeRepository EmployeesRepository, logger ditzap.Logger) *newsInteractor {
	return &newsInteractor{
		newsRepository:     newsRepository,
		employeeRepository: employeeRepository,
		logger:             logger,
	}
}

type newsInteractor struct {
	newsRepository     NewsRepository
	employeeRepository EmployeesRepository
	logger             ditzap.Logger
}

func (i *newsInteractor) Get(ctx context.Context, slug string) (*entityNews.NewsFull, error) {
	news, err := i.newsRepository.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("newsInteractor.Get: can't get news: %w", err)
	}
	return news, nil
}

func (i *newsInteractor) Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		return nil, fmt.Errorf("newsInteractor.Update: can't get session: %w", err)
	}

	// TODO:пока не реализованы клиентские методы в сервисе news-search, фильтрация производится в данном методе
	if search.GetFilterPtr() == nil {
		search.Filter = &dtoNews.SearchNewsFilter{}
	}
	// Поиск только опубликованных новостей
	search.Filter.Status = entityNews.NewsStatusPublished

	search.Visitor = &entityNews.Visitor{
		PortalID: session.ActivePortal.GetPortalID(),
	}

	result, err := i.newsRepository.Search(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("newsInteractor.Search: can't search news: %w", err)
	}

	// TODO:пока не реализованы клиентские методы в сервисе news-search, фильтрация производится в данном методе
	if result == nil || len(result.News) == 0 {
		return result, nil
	}

	news := make([]*entityNews.NewsFull, 0, len(result.News))
	for _, n := range result.News {
		if n.Status != entityNews.NewsStatusPublished {
			continue
		}

		if search.GetFilter().OnMainPage && !n.OnMain {
			continue
		}

		news = append(news, n)
	}

	result.News = news
	result.Total = len(news)

	return result, nil
}
