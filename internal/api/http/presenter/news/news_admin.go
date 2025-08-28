package news

import (
	"strings"
	"time"

	"github.com/google/uuid"

	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

type newsAdminPresenter struct {
	categoryPresenter
	commentsPresenter
}

func NewNewsAdminPresenter() *newsAdminPresenter {
	return &newsAdminPresenter{}
}

func (p *newsAdminPresenter) NewNewsToDTO(news *viewNews.NewNews) *dtoNews.NewNews {
	participants := make([]*uuid.UUID, 0, len(news.ParticipantsIDs))
	for _, id := range news.ParticipantsIDs {
		if id != nil && *id != uuid.Nil {
			participants = append(participants, id)
		}
	}

	return &dtoNews.NewNews{
		Title:           news.Title,
		Slug:            news.Slug,
		ImageID:         news.GetImageIDPtr(),
		CategoryID:      news.CategoryID,
		OrganizationID:  news.GetProviderOrganizationIDPtr(),
		ProductID:       news.GetProviderProductIDPtr(),
		Participants:    participants,
		Status:          entityNews.NewsStatusDraft,
		Author:          dtoNews.Author{},
		OnMain:          news.Properties.OnMainPage,
		Pinned:          news.Properties.MainPagePinned,
		CanDisplayViews: news.Properties.ViewsEnabled,
		CanReacts:       news.Properties.LikesEnabled,
		CanCommented:    news.Properties.CommentsEnabled,
		PublicationAt:   news.GetPublishDatePtr(),
	}
}

func (p *newsAdminPresenter) UpdateNewsToDTO(updateNews *viewNews.UpdateNews) *dtoNews.UpdateNews {
	var title *string
	if updateNews.Title != "" {
		title = &updateNews.Title
	}

	var slug *string
	if updateNews.Slug != "" {
		slug = &updateNews.Slug
	}

	var imageID *uuid.UUID
	if updateNews.ImageID != nil {
		updateImageID := *updateNews.ImageID
		imageID = &uuid.Nil
		if updateImageID != "" {
			if id, err := uuid.Parse(updateImageID); err == nil {
				imageID = &id
			} else {
				imageID = nil
			}
		}
	}

	var categoryID *uuid.UUID
	if updateNews.CategoryID != nil {
		updateCategoryID := *updateNews.CategoryID
		categoryID = &uuid.Nil
		if updateCategoryID != "" {
			if id, err := uuid.Parse(updateCategoryID); err == nil {
				categoryID = &id
			} else {
				categoryID = nil
			}
		}
	}

	var organizationID *uuid.UUID
	if updateNews.ProviderOrganizationId != nil {
		updateOrganizationID := *updateNews.ProviderOrganizationId
		organizationID = &uuid.Nil
		if updateOrganizationID != "" {
			if id, err := uuid.Parse(updateOrganizationID); err == nil {
				organizationID = &id
			} else {
				organizationID = nil
			}
		}
	}

	var productID *uuid.UUID
	if updateNews.ProviderProductId != nil {
		updateProductID := *updateNews.ProviderProductId
		productID = &uuid.Nil
		if updateProductID != "" {
			if id, err := uuid.Parse(updateProductID); err == nil {
				productID = &id
			} else {
				productID = nil
			}
		}
	}

	var participants []*uuid.UUID
	if updateNews.ParticipantsIDs != nil {
		participants = make([]*uuid.UUID, 0, len(updateNews.ParticipantsIDs))
		for _, id := range updateNews.ParticipantsIDs {
			if id != nil && *id != uuid.Nil {
				participants = append(participants, id)
			}
		}
	}

	return &dtoNews.UpdateNews{
		Title:          title,
		Slug:           slug,
		ImageID:        imageID,
		CategoryID:     categoryID,
		OrganizationID: organizationID,
		ProductID:      productID,
		Participants:   participants,
		// Выставляем статус по умолчанию NewsStatusInvalid, так как обновление статуса новости не поддерживается
		Status:          entityNews.NewsStatusInvalid,
		Body:            updateNews.Content,
		OnMain:          &updateNews.Properties.OnMainPage,
		Pinned:          &updateNews.Properties.MainPagePinned,
		CanDisplayViews: &updateNews.Properties.ViewsEnabled,
		CanReacts:       &updateNews.Properties.LikesEnabled,
		CanCommented:    &updateNews.Properties.CommentsEnabled,
		PublicationAt:   updateNews.PublishDate,
	}
}

func (p *newsAdminPresenter) SearchNewsToDTO(search *viewNews.SearchNewsRequest) *dtoNews.SearchNews {

	orderBy := dtoNews.SearchNewsOrderByTitle
	switch search.OrderBy {
	case "title":
		orderBy = dtoNews.SearchNewsOrderByTitle
	case "createDate":
		orderBy = dtoNews.SearchNewsOrderByCreatedAt
	}

	orderDirection := dto.OrderDirectionDesc
	if search.OrderType == "ASC" {
		orderDirection = dto.OrderDirectionAsc
	}

	return &dtoNews.SearchNews{
		Query: search.Query,
		Filter: &dtoNews.SearchNewsFilter{
			Status:                   p.StatusToEntity(search.Filters.Status),
			ProviderOrganizationsIds: search.Filters.ProviderOrganizationsIds,
			ProviderProductsNames:    search.Filters.ProviderProductsNames,
			CategoriesNames:          search.Filters.CategoriesNames,
			AuthorsNames:             search.Filters.AuthorsNames,
			OnMainPage:               search.Filters.OnMainPage,
			IsPinnedOnMainPage:       search.Filters.IsPinnedOnMainPage,
		},
		Pagination: dtoNews.SearchNewsPagination{
			Page:  search.Page,
			Limit: search.Limit,
		},
		Order: dtoNews.SearchNewsOrder{
			By:        orderBy,
			Direction: orderDirection,
		},
	}
}

func (p *newsAdminPresenter) StatusToEntity(status viewNews.NewsStatus) entityNews.NewsStatus {
	switch status {
	case viewNews.NewsStatusDraft:
		return entityNews.NewsStatusDraft
	case viewNews.NewsStatusWaitingPublish:
		return entityNews.NewsStatusWaitingPublish
	case viewNews.NewsStatusPublished:
		return entityNews.NewsStatusPublished
	case viewNews.NewsStatusUnpublished:
		return entityNews.NewsStatusUnpublished
	default:
		return entityNews.NewsStatusInvalid
	}
}

func (p *newsAdminPresenter) FullNewsToView(n *entityNews.NewsFull) *viewNews.News {
	if n == nil {
		return nil
	}

	news := &viewNews.News{
		Id:           n.ID,
		Slug:         n.Slug,
		Title:        n.Title,
		ImageID:      n.GetImageIDPtr(),
		Participants: p.ParticipantsToView(n.Participants),
		Category:     p.NewsCategoryToView(n.GetCategoryPtr()),
		Content:      n.Body,
		Author:       p.AuthorToView(n.Author),
		Status:       p.StatusToView(n.Status),
		Properties: viewNews.NewsProperties{
			ViewsEnabled:    n.CanDisplayViews,
			LikesEnabled:    n.CanReacts,
			CommentsEnabled: n.CanCommented,
			OnMainPage:      n.OnMain,
			MainPagePinned:  n.Pinned,
		},
		ProviderOrganization: p.NewsOrganizationToView(n.GetOrganizationPtr()),
		ProviderProduct:      p.NewsProductToView(n.GetProductPtr()),
		CreateDate:           n.GetCreatedAtPtr(),
		UpdatedAt:            n.GetUpdatedAtPtr(),
	}

	if !n.GetPublicationAt().IsZero() {
		news.PublishDate = n.GetPublicationAtPtr()
	}

	return news
}

func (p *newsAdminPresenter) NewsCategoryToView(category *entityNews.Category) *viewNews.NewsCategory {
	if category == nil {
		return nil
	}
	if category.ID == uuid.Nil {
		return nil
	}
	return &viewNews.NewsCategory{
		ID:   category.ID,
		Name: category.Name,
	}
}

func (p *newsAdminPresenter) UpdateFlagsToDTO(n *viewNews.UpdateNewsFlags) *dtoNews.UpdateFlags {
	var onMainPage *bool
	if n.OnMainPage != nil {
		onMainPage = n.OnMainPage
	}

	var isPinnedOnMain *bool
	if n.IsPinnedOnMainPage != nil {
		isPinnedOnMain = n.IsPinnedOnMainPage
	}

	var updatedAt *time.Time
	if n.UpdatedAt != nil {
		updatedAt = n.UpdatedAt
	}

	return &dtoNews.UpdateFlags{
		OnMain:    onMainPage,     // Флаг на главной странице
		Pinned:    isPinnedOnMain, // Флаг закрепления на главной странице
		UpdatedAt: updatedAt,      // Время обновления новости
	}

}

func (p *newsAdminPresenter) NewsOrganizationToView(organization *entityNews.NewsOrganization) *viewNews.NewsOrganization {
	if organization == nil {
		return nil
	}

	return &viewNews.NewsOrganization{
		ID:   organization.ID,
		Name: organization.Name,
	}
}

func (p *newsAdminPresenter) NewsProductToView(product *entityNews.NewsProduct) *viewNews.NewsProduct {
	if product == nil {
		return nil
	}

	return &viewNews.NewsProduct{
		ID:   product.ID,
		Name: product.Name,
	}
}

func (p *newsAdminPresenter) AuthorToView(author entityNews.Author) viewNews.Author {
	fullname := []string{author.LastName, author.FirstName}
	if author.GetMiddleNamePtr() != nil {
		fullname = append(fullname, author.GetMiddleName())
	}

	return viewNews.Author{
		ID:      author.GetIDPtr(),
		Name:    strings.Join(fullname, " "),
		ImageID: author.GetImageIDPtr(),
	}
}

func (p *newsAdminPresenter) StatusToView(status entityNews.NewsStatus) viewNews.NewsStatus {
	switch status {
	case entityNews.NewsStatusDraft:
		return viewNews.NewsStatusDraft
	case entityNews.NewsStatusWaitingPublish:
		return viewNews.NewsStatusWaitingPublish
	case entityNews.NewsStatusPublished:
		return viewNews.NewsStatusPublished
	case entityNews.NewsStatusUnpublished:
		return viewNews.NewsStatusUnpublished
	default:
		return viewNews.NewsStatusDraft
	}
}

func (p *newsAdminPresenter) ParticipantsToView(participants []*entityNews.Participant) []*viewNews.NewsParticipants {
	if len(participants) == 0 {
		return []*viewNews.NewsParticipants{}
	}

	participantsView := make([]*viewNews.NewsParticipants, 0, len(participants))
	for _, participant := range participants {
		if i := p.ParticipantToView(participant); i != nil {
			participantsView = append(participantsView, i)
		}
	}
	return participantsView
}

func (p *newsAdminPresenter) ParticipantToView(participant *entityNews.Participant) *viewNews.NewsParticipants {
	if participant == nil {
		return nil
	}

	fullname := []string{participant.LastName, participant.FirstName}
	if participant.GetMiddleNamePtr() != nil {
		fullname = append(fullname, participant.GetMiddleName())
	}

	return &viewNews.NewsParticipants{
		Id:      participant.GetIDPtr(),
		Name:    strings.Join(fullname, " "),
		ImageID: participant.GetImageIDPtr(),
	}
}

func (p *newsAdminPresenter) FullNewsToSearchItems(n []*entityNews.NewsFull) []*viewNews.SearchNewsResponseItem {
	r := make([]*viewNews.SearchNewsResponseItem, 0, len(n))
	for _, item := range n {
		r = append(r, p.FullNewsToSearchItem(item))
	}
	return r
}

func (p *newsAdminPresenter) FullNewsToSearchItem(n *entityNews.NewsFull) *viewNews.SearchNewsResponseItem {
	if n == nil {
		return nil
	}

	news := &viewNews.SearchNewsResponseItem{
		Id:                   n.ID,
		Slug:                 n.Slug,
		Title:                n.Title,
		ImageID:              n.GetImageIDPtr(),
		Category:             p.NewsCategoryToView(n.GetCategoryPtr()),
		ProviderOrganization: p.NewsOrganizationToView(n.GetOrganizationPtr()),
		ProviderProduct:      p.NewsProductToView(n.GetProductPtr()),
		Author:               p.AuthorToView(n.Author),
		Flags: viewNews.SearchNewsResponseFlags{
			OnMainPage:         n.OnMain,
			IsPinnedOnMainPage: n.Pinned,
		},
		Status:   p.StatusToView(n.Status),
		CreateAt: n.GetCreatedAtPtr(),
	}

	if !n.GetPublicationAt().IsZero() {
		news.PublishedAt = n.GetPublicationAtPtr()
	}

	return news
}
