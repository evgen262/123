package news

import (
	newsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/news/v1"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"github.com/google/uuid"

	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

func NewNewsMapper(sharedMapper SharedMapper) *newsMapper {
	return &newsMapper{
		sharedMapper: sharedMapper,
	}
}

type newsMapper struct {
	sharedMapper SharedMapper
}

func (m *newsMapper) NewsFullToEntity(newsPb []*newsv1.News) []*entityNews.NewsFull {
	if len(newsPb) == 0 {
		return nil
	}

	newsEntity := make([]*entityNews.NewsFull, 0, len(newsPb))
	for _, onceNewsPb := range newsPb {
		if n := m.OnceNewsFullToEntity(onceNewsPb); n != nil {
			newsEntity = append(newsEntity, n)
		}
	}
	return newsEntity
}

func (m *newsMapper) OnceNewsFullToEntity(newsPb *newsv1.News) *entityNews.NewsFull {
	id, err := uuid.Parse(newsPb.GetId())
	if err != nil {
		id = uuid.Nil
	}

	categoryID, err := uuid.Parse(newsPb.GetCategory().GetId())
	if err != nil {
		categoryID = uuid.Nil
	}

	portals := make([]*entityNews.NewsPortal, 0, len(newsPb.GetVisibility().GetPortalIds()))
	for _, portal := range newsPb.GetVisibility().GetPortalIds() {
		portals = append(portals, &entityNews.NewsPortal{
			ID:   int(portal),
			Name: "",
		})
	}

	var organization *entityNews.NewsOrganization
	if newsPb.GetOrganization().GetId() != "" {
		if id := m.sharedMapper.StringToUUIDPtr(newsPb.GetOrganization().GetId(), true); id != nil {
			organization = &entityNews.NewsOrganization{
				ID:   *id,
				Name: newsPb.GetOrganization().GetName(),
			}
		}
	}

	var product *entityNews.NewsProduct
	if newsPb.GetProduct() != nil {
		id, err := uuid.Parse(newsPb.GetProduct().GetId())
		if err == nil {
			product = &entityNews.NewsProduct{
				ID:   id,
				Name: newsPb.GetProduct().GetName(),
			}
		}
	}

	participants := make([]*entityNews.Participant, 0, len(newsPb.GetParticipants()))
	for _, participant := range newsPb.GetParticipants() {
		participants = append(participants, &entityNews.Participant{
			ID:         m.sharedMapper.StringToUUIDPtr(participant.GetId(), true),
			LastName:   participant.GetLastName(),
			FirstName:  participant.GetFirstName(),
			MiddleName: m.sharedMapper.StringValueToStringPtr(participant.GetMiddleName()),
			ImageID:    m.sharedMapper.StringValueToUUIDPtr(participant.GetImageId(), true),
		})
	}

	return &entityNews.NewsFull{
		ID:      id,
		Title:   newsPb.GetTitle(),
		Slug:    newsPb.GetUrl(),
		ImageID: m.sharedMapper.StringValueToUUIDPtr(newsPb.GetImageId(), true),
		Category: &entityNews.Category{
			ID:        categoryID,
			Name:      newsPb.GetCategory().GetName(),
			UpdatedAt: nil,
			Visibility: entityNews.CategoryVisibility{
				ComplexIDs: m.sharedMapper.Int32SliceToInt(newsPb.GetCategory().GetVisibility().GetComplexIds()),
				OIVs:       m.sharedMapper.Int32SliceToInt(newsPb.GetCategory().GetVisibility().GetPortalIds()),
			},
		},
		Organization: organization,
		Product:      product,
		Participants: participants,
		Status:       m.StatusToEntity(newsPb.GetStatus()),
		Body:         newsPb.GetBody(),
		Author: entityNews.Author{
			ID:         m.sharedMapper.StringToUUIDPtr(newsPb.GetAuthor().GetId(), true),
			LastName:   newsPb.GetAuthor().GetLastName(),
			FirstName:  newsPb.GetAuthor().GetFirstName(),
			MiddleName: m.sharedMapper.StringValueToStringPtr(newsPb.GetAuthor().GetMiddleName()),
			ImageID:    m.sharedMapper.StringValueToUUIDPtr(newsPb.GetAuthor().GetImageId(), true),
		},
		OnMain:          newsPb.GetOnMain(),
		Pinned:          newsPb.GetPinned(),
		CanDisplayViews: newsPb.GetCanDisplayViews(),
		Views:           int(newsPb.GetViewCount()),
		CanReacts:       newsPb.GetCanReacts(),
		Likes:           int(newsPb.GetLikesCount()),
		CanCommented:    newsPb.GetCanCommented(),
		Comments:        m.CommentsToEntity(newsPb.GetComments()),
		Visibility: &entityNews.NewsNamedVisibility{
			Portals: portals,
		},
		CreatedAt:     m.sharedMapper.TimestampToTime(newsPb.GetCreateTime()),
		UpdatedAt:     m.sharedMapper.TimestampToTime(newsPb.GetUpdateTime()),
		PublicationAt: m.sharedMapper.TimestampToTime(newsPb.GetPublicationTime()),
	}
}

func (m *newsMapper) NewsToEntity(newsPb []*newsv1.News) []*entityNews.News {
	if len(newsPb) == 0 {
		return nil
	}

	newsEntity := make([]*entityNews.News, 0, len(newsPb))
	for _, onceNewsPb := range newsPb {
		if n := m.OnceNewsToEntity(onceNewsPb); n != nil {
			newsEntity = append(newsEntity, n)
		}
	}
	return newsEntity
}

func (m *newsMapper) OnceNewsToEntity(newsPb *newsv1.News) *entityNews.News {
	id, err := uuid.Parse(newsPb.GetId())
	if err != nil {
		id = uuid.Nil
	}

	categoryID, err := uuid.Parse(newsPb.GetCategory().GetId())
	if err != nil {
		categoryID = uuid.Nil
	}

	var visibility *entityNews.NewsVisibility
	if newsPb.GetVisibility() != nil {
		visibility = &entityNews.NewsVisibility{
			ComplexIDs: m.sharedMapper.Int32SliceToInt(newsPb.GetVisibility().GetComplexIds()),
			PortalsIDs: m.sharedMapper.Int32SliceToInt(newsPb.GetVisibility().GetPortalIds()),
		}
	}

	return &entityNews.News{
		ID:             id,
		Title:          newsPb.GetTitle(),
		Slug:           newsPb.GetUrl(),
		ImageID:        m.sharedMapper.StringValueToUUIDPtr(newsPb.GetImageId(), true),
		CategoryID:     categoryID,
		OrganizationID: m.sharedMapper.StringToUUIDPtr(newsPb.GetOrganization().GetId(), true),
		ProductID:      m.sharedMapper.StringToUUIDPtr(newsPb.GetProduct().GetId(), true),
		Participants:   m.ParticipantsToUUID(newsPb.GetParticipants()),
		Status:         m.StatusToEntity(newsPb.GetStatus()),
		Body:           newsPb.GetBody(),
		Author: entityNews.Author{
			ID:         m.sharedMapper.StringToUUIDPtr(newsPb.GetAuthor().GetId(), true),
			LastName:   newsPb.GetAuthor().GetLastName(),
			FirstName:  newsPb.GetAuthor().GetFirstName(),
			MiddleName: m.sharedMapper.StringValueToStringPtr(newsPb.GetAuthor().GetMiddleName()),
			ImageID:    m.sharedMapper.StringValueToUUIDPtr(newsPb.GetAuthor().GetImageId(), true),
		},
		OnMain:          newsPb.GetOnMain(),
		Pinned:          newsPb.GetPinned(),
		CanDisplayViews: newsPb.GetCanDisplayViews(),
		CanReacts:       newsPb.GetCanReacts(),
		CanCommented:    newsPb.GetCanCommented(),
		Visibility:      visibility,
		CreatedAt:       m.sharedMapper.TimestampToTime(newsPb.GetCreateTime()),
		UpdatedAt:       m.sharedMapper.TimestampToTime(newsPb.GetUpdateTime()),
		PublicationAt:   m.sharedMapper.TimestampToTime(newsPb.GetPublicationTime()),
	}
}

func (m *newsMapper) ParticipantsIdsToPb(participants []*uuid.UUID) []string {
	if participants == nil {
		return nil
	}
	participantsPb := make([]string, 0, len(participants))
	for _, participant := range participants {
		if participant == nil {
			continue
		}
		participantsPb = append(participantsPb, participant.String())
	}
	return participantsPb
}

func (m *newsMapper) StatusToPb(status entityNews.NewsStatus) newsv1.NewsStatus {
	switch status {
	case entityNews.NewsStatusDraft:
		return newsv1.NewsStatus_NEWS_STATUS_DRAFT
	case entityNews.NewsStatusWaitingPublish:
		return newsv1.NewsStatus_NEWS_STATUS_WAITING_PUBLISH
	case entityNews.NewsStatusPublished:
		return newsv1.NewsStatus_NEWS_STATUS_PUBLISHED
	case entityNews.NewsStatusUnpublished:
		return newsv1.NewsStatus_NEWS_STATUS_UNPUBLISHED
	default:
		return newsv1.NewsStatus_NEWS_STATUS_INVALID
	}
}

func (m *newsMapper) StatusToEntity(statusPb newsv1.NewsStatus) entityNews.NewsStatus {
	switch statusPb {
	case newsv1.NewsStatus_NEWS_STATUS_DRAFT:
		return entityNews.NewsStatusDraft
	case newsv1.NewsStatus_NEWS_STATUS_WAITING_PUBLISH:
		return entityNews.NewsStatusWaitingPublish
	case newsv1.NewsStatus_NEWS_STATUS_PUBLISHED:
		return entityNews.NewsStatusPublished
	case newsv1.NewsStatus_NEWS_STATUS_UNPUBLISHED:
		return entityNews.NewsStatusUnpublished
	default:
		return entityNews.NewsStatusInvalid
	}
}

func (m *newsMapper) CommentsToEntity(commentsPb []*newsv1.Comment) []*entityNews.NewsComment {
	if len(commentsPb) == 0 {
		return nil
	}

	comments := make([]*entityNews.NewsComment, 0, len(commentsPb))
	for _, commentPb := range commentsPb {
		if c := m.CommentToEntity(commentPb); c != nil {
			comments = append(comments, c)
		}
	}
	return comments
}
func (m *newsMapper) CommentToEntity(comment *newsv1.Comment) *entityNews.NewsComment {
	id, err := uuid.Parse(comment.GetId())
	if err != nil {
		id = uuid.Nil
	}

	return &entityNews.NewsComment{
		ID:      id,
		Message: comment.GetText(),
		Author: entityNews.Author{
			ID:         m.sharedMapper.StringToUUIDPtr(comment.GetAuthor().GetId(), true),
			LastName:   comment.GetAuthor().GetLastName(),
			FirstName:  comment.GetAuthor().GetFirstName(),
			MiddleName: m.sharedMapper.StringValueToStringPtr(comment.GetAuthor().GetMiddleName()),
			ImageID:    m.sharedMapper.StringValueToUUIDPtr(comment.GetAuthor().GetImageId(), true),
		},
		CreatedAt: m.sharedMapper.TimestampToTime(comment.GetCreateTime()),
		UpdatedAt: m.sharedMapper.TimestampToTime(comment.GetUpdateTime()),
		DeletedAt: m.sharedMapper.TimestampToTime(comment.GetDeleteTime()),
	}
}

func (m *newsMapper) VisibilityToPb(vis *entityNews.NewsVisibility) *newsv1.Visibility {
	if m == nil || vis == nil {
		return nil
	}
	return &newsv1.Visibility{
		PortalIds:  m.sharedMapper.IntSliceToInt32(vis.PortalsIDs),
		ComplexIds: m.sharedMapper.IntSliceToInt32(vis.ComplexIDs),
	}
}
func (m *newsMapper) ParticipantsToUUID(participants []*newsv1.Participant) []*uuid.UUID {
	if participants == nil {
		return nil
	}
	participantsPb := make([]*uuid.UUID, 0, len(participants))
	for _, participant := range participants {
		if participant == nil || participant.GetId() == "" {
			continue
		}
		participantsPb = append(participantsPb, m.sharedMapper.StringToUUIDPtr(participant.GetId(), true))
	}
	return participantsPb
}
func (m *newsMapper) NewCommentToPb(comment *dtoNews.NewComment) *newsv1.Comment {
	if comment == nil {
		return nil
	}

	return &newsv1.Comment{
		Id:         "",
		Text:       "",
		Author:     nil,
		CreateTime: nil,
		UpdateTime: nil,
		DeleteTime: nil,
		ByUser:     false,
	}
	return nil
}
