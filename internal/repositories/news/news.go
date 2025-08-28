package news

import (
	"context"
	"fmt"

	newsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/news/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/shared/v1"
	organizationsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsv2/organizations/v1"
	portalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsv2/portals/v1"
	portalsSharedv2 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsv2/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

func NewNewsRepository(
	newsApi newsv1.NewsAPIClient,
	portalsApi portalsv2.PortalsAPIClient,
	organizationsApi organizationsv2.OrganizationsAPIClient,
	newsMapper NewsMapper,
	sharedMapper repositories.SharedMapper,
	logger ditzap.Logger,
) *newsRepository {
	return &newsRepository{
		newsApi:          newsApi,
		portalsApi:       portalsApi,
		organizationsApi: organizationsApi,
		newsMapper:       newsMapper,
		sharedMapper:     sharedMapper,
		logger:           logger,
	}
}

type newsRepository struct {
	newsApi          newsv1.NewsAPIClient
	portalsApi       portalsv2.PortalsAPIClient
	organizationsApi organizationsv2.OrganizationsAPIClient
	newsMapper       NewsMapper
	sharedMapper     repositories.SharedMapper
	logger           ditzap.Logger
}

func (r *newsRepository) Create(ctx context.Context, news *dtoNews.NewNews) (uuid.UUID, error) {
	if news == nil {
		return uuid.UUID{}, fmt.Errorf("newsRepository.Create: %w", diterrors.ErrInputEmpty)
	}
	var imageID *wrapperspb.StringValue
	if news.GetImageIDPtr() != nil {
		imageID = &wrapperspb.StringValue{
			Value: news.GetImageID().String(),
		}
	}

	var visibility *newsv1.Visibility
	if news.GetVisibilityPtr() != nil {
		visibility = &newsv1.Visibility{
			PortalIds: r.sharedMapper.IntSliceToInt32(news.GetVisibility().PortalsIDs),
		}
	}

	var organizationID string
	if orgID := news.GetOrganizationID(); orgID != uuid.Nil {
		organizationID = orgID.String()
	}

	var productID string
	if prodID := news.GetProductID(); prodID != uuid.Nil {
		productID = prodID.String()
	}

	resp, err := r.newsApi.Create(ctx, &newsv1.CreateRequest{
		Title:           news.GetTitle(),
		Url:             news.GetSlug(),
		ImageId:         imageID,
		CategoryId:      news.GetCategoryID().String(),
		OrganizationId:  organizationID,
		ProductId:       productID,
		Participants:    r.newsMapper.ParticipantsIdsToPb(news.GetParticipants()),
		Status:          r.newsMapper.StatusToPb(news.GetValidStatus()),
		Body:            news.GetBody(),
		Author:          news.GetAuthor().ID.String(),
		OnMain:          news.GetOnMain(),
		Pinned:          news.GetPinned(),
		CanDisplayViews: news.GetCanDisplayViews(),
		CanReacts:       news.GetCanReacts(),
		CanCommented:    news.GetCanCommented(),
		Visibility:      visibility,
		PublicationTime: r.sharedMapper.TimeToTimestamp(news.GetPublicationAtPtr()),
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("newsRepository.Create: can't create news: %w", diterrors.GrpcErrorToError(err))
	}

	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("newsRepository.Create: can't parse news id [%s]: %w", resp.GetId(), err)
	}

	return id, nil
}

func (r *newsRepository) Update(ctx context.Context, id uuid.UUID, news *dtoNews.UpdateNews) (*entityNews.News, error) {
	if news == nil {
		return nil, fmt.Errorf("newsRepository.Update: %w", diterrors.ErrInputEmpty)
	}

	var participants *newsv1.UpdateRequest_Participants
	if news.Participants != nil {
		participants = &newsv1.UpdateRequest_Participants{
			Ids: r.newsMapper.ParticipantsIdsToPb(news.GetParticipants()),
		}
	}

	updateResponse, err := r.newsApi.Update(ctx, &newsv1.UpdateRequest{
		Id:              id.String(),
		Title:           r.sharedMapper.StringValue(news.GetTitlePtr()),
		Url:             r.sharedMapper.StringValue(news.GetSlugPtr()),
		ImageId:         r.sharedMapper.UUIDStringValue(news.GetImageIDPtr()),
		CategoryId:      r.sharedMapper.UUIDStringValue(news.GetCategoryIDPtr()),
		OrganizationId:  r.sharedMapper.UUIDStringValue(news.GetOrganizationIDPtr()),
		ProductId:       r.sharedMapper.UUIDStringValue(news.GetProductIDPtr()),
		Participants:    participants,
		Status:          r.newsMapper.StatusToPb(news.GetValidStatus()),
		Body:            r.sharedMapper.BytesValue(news.GetBody()),
		OnMain:          r.sharedMapper.BoolValue(news.GetOnMainPtr()),
		Pinned:          r.sharedMapper.BoolValue(news.GetPinnedPtr()),
		CanDisplayViews: r.sharedMapper.BoolValue(news.GetCanDisplayViewsPtr()),
		CanReacts:       r.sharedMapper.BoolValue(news.GetCanReactsPtr()),
		CanCommented:    r.sharedMapper.BoolValue(news.GetCanCommentedPtr()),
		PublicationTime: r.sharedMapper.TimeToTimestamp(news.GetPublicationAtPtr()),
	})
	if err != nil {
		return nil, fmt.Errorf("newsRepository.Update: can't update news: %w", diterrors.GrpcErrorToError(err))
	}

	updatedNews := r.newsMapper.OnceNewsToEntity(updateResponse.GetNews())
	if updatedNews == nil || updatedNews.ID == uuid.Nil {
		return nil, fmt.Errorf("newsRepository.Update: can't parse updated news: %w", diterrors.ErrInputEmpty)
	}

	return updatedNews, nil
}

func (r *newsRepository) Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error) {
	if search == nil {
		found, err := r.newsApi.Filter(ctx, &newsv1.FilterRequest{})
		if err != nil {
			return nil, fmt.Errorf("newsRepository.Search: can't search news without params: %w", diterrors.GrpcErrorToError(err))
		}
		// TODO: POKAZ Total для показа берется длина новостей
		return &dtoNews.SearchNewsResult{
			News:  r.newsMapper.NewsFullToEntity(found.GetNews()),
			Total: len(found.GetNews()),
		}, nil
	}

	// TODO: Должен работать с news-search
	found, err := r.newsApi.Filter(ctx, &newsv1.FilterRequest{
		Filters: &newsv1.FilterRequest_Title{Title: search.Query},
		Options: &newsv1.FilterRequest_Options{
			WithComments: false,
			OnlyMain:     false,
		},
		// TODO: Сделать пагинацию
		Pagination: &sharedv1.ScrollPaginationRequest{
			Limit: int32(search.Pagination.Limit),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("newsRepository.Search: can't search news: %w", diterrors.GrpcErrorToError(err))
	}

	sNews := r.newsMapper.NewsFullToEntity(found.GetNews())

	portalsIds := make([]int32, 0, len(sNews))
	organizationsIds := make(uuid.UUIDs, 0, len(sNews))
	for _, n := range sNews {
		if n.Organization == nil {
			continue
		}
		organizationsIds = append(organizationsIds, n.Organization.ID)
		if n.GetVisibilityPtr() == nil {
			continue
		}
		for _, portal := range n.GetVisibilityPtr().Portals {
			portalsIds = append(portalsIds, int32(portal.ID))
		}
	}

	portals := make(map[int]entityNews.NewsPortal, 0)
	if p, err := r.portalsApi.Filter(ctx, &portalsv2.FilterRequest{
		Filters: &portalsv2.FilterRequest_Filters{
			Ids: portalsIds,
		},
		Options: &portalsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
	}); err == nil {
		for _, portal := range p.GetPortals() {
			portals[int(portal.GetPortal().GetId())] = entityNews.NewsPortal{
				ID:   int(portal.GetPortal().GetId()),
				Name: portal.GetPortal().GetName(),
			}
		}
	} else {
		r.logger.Warn("newsRepository.Search: can't get portals", zap.Any("portal_ids", portalsIds), zap.Error(err))
	}

	orgs := make(map[uuid.UUID]entityNews.NewsOrganization, 0)
	if o, err := r.organizationsApi.Filter(ctx, &organizationsv2.FilterRequest{
		Filters: &organizationsv2.FilterRequest_Filters{
			Ids: organizationsIds.Strings(),
		},
		Options: &organizationsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
		Pagination: &portalsSharedv2.PaginationRequest{
			Limit: -1,
		},
	}); err == nil {
		for _, org := range o.GetOrganizations() {
			if orgId, err := uuid.Parse(org.GetId()); err == nil {
				orgs[orgId] = entityNews.NewsOrganization{
					ID:   orgId,
					Name: org.GetName(),
				}
			}
		}
	} else {
		r.logger.Warn("newsRepository.Search: can't get organizations", zap.Strings("org_ids", organizationsIds.Strings()), zap.Error(err))
	}

	for i, n := range sNews {
		if n.Organization == nil {
			continue
		}
		if _, ok := orgs[n.Organization.ID]; ok {
			sNews[i].Organization.Name = orgs[n.Organization.ID].Name
		}
		vis := n.GetVisibilityPtr()
		if vis == nil {
			continue
		}
		for j, portal := range vis.Portals {
			if _, ok := portals[portal.ID]; ok {
				vis.Portals[j].Name = portals[portal.ID].Name
			} else {
				vis.Portals[j] = nil
			}
		}
	}
	// TODO: POKAZ Total для показа берется длина новостей
	return &dtoNews.SearchNewsResult{
		News:  sNews,
		Total: len(sNews),
	}, nil
}

func (r *newsRepository) Get(ctx context.Context, id uuid.UUID) (*entityNews.NewsFull, error) {
	found, err := r.newsApi.Get(ctx, &newsv1.GetRequest{
		By: &newsv1.GetRequest_Id{
			Id: id.String(),
		},
		Options: &newsv1.GetRequest_Options{
			WithComments: true,
			OnlyMain:     false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("newsRepository.Get: can't get news: %w", diterrors.GrpcErrorToError(err))
	}

	if found.GetNews() == nil {
		return nil, fmt.Errorf("newsRepository.Get: news not found: %w", diterrors.ErrNotFound)
	}

	sNews := r.newsMapper.OnceNewsFullToEntity(found.GetNews())

	portalsIds := make([]int32, 0, len(sNews.GetVisibilityPtr().Portals))
	for _, p := range sNews.GetVisibilityPtr().Portals {
		portalsIds = append(portalsIds, int32(p.ID))
	}

	portals := make(map[int]entityNews.NewsPortal, 0)
	if p, err := r.portalsApi.Filter(ctx, &portalsv2.FilterRequest{
		Filters: &portalsv2.FilterRequest_Filters{
			Ids: portalsIds,
		},
		Options: &portalsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
	}); err == nil {
		for _, portal := range p.GetPortals() {
			portals[int(portal.GetPortal().GetId())] = entityNews.NewsPortal{
				ID:   int(portal.GetPortal().GetId()),
				Name: portal.GetPortal().GetName(),
			}
		}
	} else {
		r.logger.Warn("newsRepository.Get: can't get portals", zap.Any("portal_ids", portalsIds), zap.Error(err))
	}

	var nOrg *entityNews.NewsOrganization
	if o, err := r.organizationsApi.Filter(ctx, &organizationsv2.FilterRequest{
		Filters: &organizationsv2.FilterRequest_Filters{
			Ids: []string{sNews.ID.String()},
		},
		Options: &organizationsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
		Pagination: &portalsSharedv2.PaginationRequest{
			Limit: -1,
		},
	}); err == nil {
		for _, org := range o.GetOrganizations() {
			if orgId, err := uuid.Parse(org.GetId()); err == nil {
				nOrg = &entityNews.NewsOrganization{
					ID:   orgId,
					Name: org.GetName(),
				}
			}
		}
		if nOrg != nil {
			sNews.Organization.Name = nOrg.Name
		}
		/*
			else {
				sNews.Organization = nil
			}
		*/
	} else {
		r.logger.Warn("newsRepository.Get: can't get organizations", zap.String("org_id", sNews.ID.String()), zap.Error(err))
	}

	if sNews.GetVisibilityPtr() != nil {
		vis := sNews.GetVisibilityPtr()
		for i, portal := range vis.Portals {
			if _, ok := portals[portal.ID]; ok {
				vis.Portals[i].Name = portals[portal.ID].Name
			} else {
				vis.Portals[i] = nil
			}
		}
	}

	// TODO: Было реализовано костылями для показа. Реализовать на стороне фасада обогащение полями участников новости

	return sNews, nil
}

func (r *newsRepository) GetBySlug(ctx context.Context, slug string) (*entityNews.NewsFull, error) {
	found, err := r.newsApi.Get(ctx, &newsv1.GetRequest{
		By: &newsv1.GetRequest_Slug{
			Slug: slug,
		},
		Options: &newsv1.GetRequest_Options{
			WithComments: true,
			OnlyMain:     false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("newsRepository.GetBySlug: can't get news: %w", diterrors.GrpcErrorToError(err))
	}

	if found.GetNews() == nil {
		return nil, fmt.Errorf("newsRepository.GetBySlug: news not found: %w", diterrors.ErrNotFound)
	}

	sNews := r.newsMapper.OnceNewsFullToEntity(found.GetNews())

	portalsIds := make([]int32, 0, len(sNews.GetVisibilityPtr().Portals))
	for _, p := range sNews.GetVisibilityPtr().Portals {
		portalsIds = append(portalsIds, int32(p.ID))
	}

	portals := make(map[int]entityNews.NewsPortal, 0)
	if p, err := r.portalsApi.Filter(ctx, &portalsv2.FilterRequest{
		Filters: &portalsv2.FilterRequest_Filters{
			Ids: portalsIds,
		},
		Options: &portalsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
	}); err == nil {
		for _, portal := range p.GetPortals() {
			portals[int(portal.GetPortal().GetId())] = entityNews.NewsPortal{
				ID:   int(portal.GetPortal().GetId()),
				Name: portal.GetPortal().GetName(),
			}
		}
	} else {
		r.logger.Warn("newsRepository.Get: can't get portals", zap.Any("portal_ids", portalsIds), zap.Error(err))
	}

	var nOrg *entityNews.NewsOrganization
	if o, err := r.organizationsApi.Filter(ctx, &organizationsv2.FilterRequest{
		Filters: &organizationsv2.FilterRequest_Filters{
			Ids: []string{sNews.ID.String()},
		},
		Options: &organizationsv2.FilterRequest_Options{
			WithDeleted:  true,
			WithDisabled: true,
		},
		Pagination: &portalsSharedv2.PaginationRequest{
			Limit: -1,
		},
	}); err == nil {
		for _, org := range o.GetOrganizations() {
			if orgId, err := uuid.Parse(org.GetId()); err == nil {
				nOrg = &entityNews.NewsOrganization{
					ID:   orgId,
					Name: org.GetName(),
				}
			}
		}
		if nOrg != nil {
			sNews.Organization.Name = nOrg.Name
		} /*
			else {
				sNews.Organization = nil
			}
		*/
	} else {
		r.logger.Warn("newsRepository.Get: can't get organizations", zap.String("org_id", sNews.ID.String()), zap.Error(err))
	}

	if sNews.GetVisibilityPtr() != nil {
		vis := sNews.GetVisibilityPtr()
		for i, portal := range vis.Portals {
			if _, ok := portals[portal.ID]; ok {
				vis.Portals[i].Name = portals[portal.ID].Name
			} else {
				vis.Portals[i] = nil
			}
		}
	}

	return sNews, nil
}

func (r *newsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Опускаем результат, т.к. нет нужды.
	if _, err := r.newsApi.Delete(ctx, &newsv1.DeleteRequest{
		Id: id.String(),
	}); err != nil {
		return fmt.Errorf("newsRepository.Delete: can't delete news: %w", diterrors.GrpcErrorToError(err))
	}
	return nil
}
