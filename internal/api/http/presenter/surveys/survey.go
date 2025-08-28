package surveys

import (
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
)

var respondentTypeToEntityDictionary = map[int]entitySurvey.RespondentType{
	0: entitySurvey.RespondentTypeUnknown,
	1: entitySurvey.RespondentTypeAnonymous,
	2: entitySurvey.RespondentTypeUser,
}

var questionTypeToEntityDictionary = map[string]entitySurvey.QuestionType{
	"text":        entitySurvey.QuestionTypeText,
	"inputNumber": entitySurvey.QuestionTypeNumber,
	"checkbox":    entitySurvey.QuestionTypeCheckbox,
	"radio":       entitySurvey.QuestionTypeRadio,
	"checkboxImg": entitySurvey.QuestionTypeCheckboxImg,
	"radioImg":    entitySurvey.QuestionTypeRadioImg,
}

var contentTypeToEntityDictionary = map[int]entitySurvey.ContentType{
	1: entitySurvey.ContentTypeText,
	2: entitySurvey.ContentTypeDigit,
}

type surveysPresenter struct {
	logger ditzap.Logger
}

func NewSurveysPresenter(logger ditzap.Logger) *surveysPresenter {
	return &surveysPresenter{logger: logger}
}

func (sp surveysPresenter) ToNewEntity(survey *viewSurveys.NewSurvey) *entitySurvey.Survey {
	newSurvey := &entitySurvey.Survey{
		Title:             survey.Title,
		Description:       survey.Description,
		ActivePeriodStart: survey.ActivePeriodStart,
		ActivePeriodEnd:   survey.ActivePeriodEnd,
		Author:            survey.Author,
		IsPublished:       survey.IsPublished,
		Respondent: &entitySurvey.SurveyRespondent{
			Type: sp.respondentTypeToEntity(survey.RespondentType),
		},
		Questions: sp.questionsToNewEntity(survey.Questions),
	}

	if len(survey.Respondents) != 0 && entitySurvey.RespondentType(survey.RespondentType) == entitySurvey.RespondentTypeUser {
		newSurvey.Respondent.Ids = sp.respondentsToIDs(survey.Respondents)
	}

	return newSurvey
}

func (sp surveysPresenter) respondentsToIDs(respondents []*viewSurveys.SurveyRespondent) entitySurvey.RespondentIDs {
	respondentsIDs := make(entitySurvey.RespondentIDs, 0, len(respondents))
	for _, respondent := range respondents {
		respondentsIDs = append(respondentsIDs, entitySurvey.RespondentID(respondent.ID))
	}

	return respondentsIDs
}

func (sp surveysPresenter) respondentTypeToEntity(respondentType int) entitySurvey.RespondentType {
	value, ok := respondentTypeToEntityDictionary[respondentType]
	if !ok {
		sp.logger.Warn("unknown respondent type, switched to default type")
		return entitySurvey.RespondentTypeAll
	}

	return value
}

func (sp surveysPresenter) questionsToNewEntity(questions []*viewSurveys.NewSurveyQuestion) []*entitySurvey.Question {
	newQuestions := make([]*entitySurvey.Question, 0, len(questions))
	for _, question := range questions {
		newQuestion := &entitySurvey.Question{
			Text:       question.Text,
			Type:       sp.questionTypeToEntity(question.Type),
			IsRequired: question.IsRequired,
			IsActive:   question.IsActive,
			Weight:     question.Weight,
			Answers:    sp.answerVariantsToNewEntity(question.AnswerVariants),
		}
		if rules := question.Rules; rules != nil {
			newRules := &entitySurvey.QuestionRules{}
			if rules.PickMinCount != nil {
				newRules.PickMinCount = rules.PickMinCount
			}
			if rules.PickMaxCount != nil {
				newRules.PickMaxCount = rules.PickMaxCount
			}
			newQuestion.Rules = newRules
		}
		newQuestions = append(newQuestions, newQuestion)
	}

	return newQuestions
}

func (sp surveysPresenter) questionTypeToEntity(questionType string) entitySurvey.QuestionType {
	value, ok := questionTypeToEntityDictionary[questionType]
	if !ok {
		return entitySurvey.QuestionTypeInvalid
	}

	return value
}

func (sp surveysPresenter) answerVariantsToNewEntity(answerVariants []*viewSurveys.NewSurveyAnswerVariant) []*entitySurvey.AnswerVariant {
	newAnswerVariants := make([]*entitySurvey.AnswerVariant, 0, len(answerVariants))
	for _, answerVariant := range answerVariants {
		newAnswerVariant := &entitySurvey.AnswerVariant{
			Text:        answerVariant.Text,
			WithContent: answerVariant.WithContent,
			Weight:      answerVariant.Weight,
		}
		if answerVariant.WithContent && answerVariant.ContentType != nil {
			newAnswerVariant.ContentType = sp.contentTypeToEntity(*answerVariant.ContentType)
		}

		if answerVariant.WithContent && answerVariant.Rules != nil {
			newAnswerVariant.Rules = &entitySurvey.AnswerRules{
				ContentMinLength: answerVariant.Rules.ContentMinLength,
				ContentMaxLength: answerVariant.Rules.ContentMaxLength,
				ContentMinDigit:  answerVariant.Rules.ContentMinDigit,
				ContentMaxDigit:  answerVariant.Rules.ContentMaxDigit,
			}
		}

		if answerVariant.Image != nil {
			newAnswerVariant.Image = &entitySurvey.Image{
				ID: entitySurvey.ImageID(answerVariant.Image.ID),
				ExternalImageInfo: &entitySurvey.ExternalProperties{
					ID:       answerVariant.Image.ImageExternalID,
					FileName: answerVariant.Image.ImageExternalFileName,
					URL:      answerVariant.Image.ImageExternalURL,
					Size:     int64(answerVariant.Image.ImageExternalSize),
				},
			}
		}
		newAnswerVariants = append(newAnswerVariants, newAnswerVariant)
	}

	return newAnswerVariants
}

func (sp surveysPresenter) contentTypeToEntity(contentType int) entitySurvey.ContentType {
	value, ok := contentTypeToEntityDictionary[contentType]
	if !ok {
		return entitySurvey.ContentTypeInvalid
	}

	return value
}

func (sp surveysPresenter) ToEntity(survey *viewSurveys.Survey) *entitySurvey.Survey {
	surveyID := entitySurvey.SurveyID(survey.ID)
	newSurvey := &entitySurvey.Survey{
		ID:                &surveyID,
		Title:             survey.Title,
		Description:       survey.Description,
		ActivePeriodStart: survey.ActivePeriodStart,
		ActivePeriodEnd:   survey.ActivePeriodEnd,
		Author:            survey.Author,
		IsPublished:       survey.IsPublished,
		CreatedAt:         &survey.CreatedAt,
		Respondent: &entitySurvey.SurveyRespondent{
			Type: sp.respondentTypeToEntity(survey.RespondentType),
		},
		Questions: sp.questionsToEntity(survey.Questions),
	}

	if len(survey.Respondents) != 0 && entitySurvey.RespondentType(survey.RespondentType) == entitySurvey.RespondentTypeUser {
		newSurvey.Respondent.Ids = sp.respondentsToIDs(survey.Respondents)
	}
	if survey.UpdatedAt != nil {
		newSurvey.UpdatedAt = survey.UpdatedAt
	}
	if survey.DeletedAt != nil {
		newSurvey.DeletedAt = survey.DeletedAt
	}

	return newSurvey
}

func (sp surveysPresenter) questionsToEntity(questions []*viewSurveys.SurveyQuestion) []*entitySurvey.Question {
	newQuestions := make([]*entitySurvey.Question, 0, len(questions))
	for _, question := range questions {
		qID := entitySurvey.QuestionID(question.ID)
		sID := entitySurvey.SurveyID(question.SurveyID)
		newQuestion := &entitySurvey.Question{
			ID:         &qID,
			SurveyID:   &sID,
			Text:       question.Text,
			Type:       sp.questionTypeToEntity(question.Type),
			IsRequired: question.IsRequired,
			IsActive:   question.IsActive,
			Weight:     question.Weight,
			Answers:    sp.answerVariantsToEntity(question.AnswerVariants),
		}
		if rules := question.Rules; rules != nil {
			newRules := &entitySurvey.QuestionRules{}
			if rules.PickMinCount != nil {
				newRules.PickMinCount = rules.PickMinCount
			}
			if rules.PickMaxCount != nil {
				newRules.PickMaxCount = rules.PickMaxCount
			}
			newQuestion.Rules = newRules
		}

		if question.DeletedAt != nil {
			newQuestion.DeletedAt = question.DeletedAt
		}

		newQuestions = append(newQuestions, newQuestion)
	}

	return newQuestions
}

func (sp surveysPresenter) answerVariantsToEntity(answerVariants []*viewSurveys.SurveyAnswerVariant) []*entitySurvey.AnswerVariant {
	newAnswerVariants := make([]*entitySurvey.AnswerVariant, 0, len(answerVariants))
	for _, answerVariant := range answerVariants {
		avID := entitySurvey.AnswerID(answerVariant.ID)
		qID := entitySurvey.QuestionID(answerVariant.QuestionID)
		newAnswerVariant := &entitySurvey.AnswerVariant{
			ID:          &avID,
			QuestionID:  &qID,
			Text:        answerVariant.Text,
			WithContent: answerVariant.WithContent,
			Weight:      answerVariant.Weight,
		}
		if answerVariant.WithContent && answerVariant.ContentType != nil {
			newAnswerVariant.ContentType = sp.contentTypeToEntity(*answerVariant.ContentType)
		}
		if answerVariant.DeletedAt != nil {
			newAnswerVariant.DeletedAt = answerVariant.DeletedAt
		}

		if answerVariant.WithContent && answerVariant.Rules != nil {
			newAnswerVariant.Rules = &entitySurvey.AnswerRules{
				ContentMinLength: answerVariant.Rules.ContentMinLength,
				ContentMaxLength: answerVariant.Rules.ContentMaxLength,
				ContentMinDigit:  answerVariant.Rules.ContentMinDigit,
				ContentMaxDigit:  answerVariant.Rules.ContentMaxDigit,
			}
		}

		if answerVariant.Image != nil {
			newAnswerVariant.Image = &entitySurvey.Image{
				ID: entitySurvey.ImageID(answerVariant.Image.ID),
				ExternalImageInfo: &entitySurvey.ExternalProperties{
					ID:       answerVariant.Image.ImageExternalID,
					FileName: answerVariant.Image.ImageExternalFileName,
					URL:      answerVariant.Image.ImageExternalURL,
					Size:     int64(answerVariant.Image.ImageExternalSize),
				},
			}
		}
		newAnswerVariants = append(newAnswerVariants, newAnswerVariant)
	}

	return newAnswerVariants
}

func (sp surveysPresenter) ToView(survey *entitySurvey.Survey) *viewSurveys.Survey {
	viewSurvey := &viewSurveys.Survey{
		ID:                uuid.UUID(*survey.GetID()),
		Title:             survey.Title,
		Description:       survey.Description,
		ActivePeriodStart: survey.ActivePeriodStart,
		ActivePeriodEnd:   survey.ActivePeriodEnd,
		Author:            survey.Author,
		IsPublished:       survey.IsPublished,
		RespondentType:    int(survey.Respondent.Type),
		CreatedAt:         *survey.GetCreatedAt(),
		Respondents:       sp.respondentsToView(survey.GetRespondent().Ids),
		Questions:         sp.questionsToView(survey.Questions),
	}
	if survey.GetUpdatedAt() != nil {
		viewSurvey.UpdatedAt = survey.GetUpdatedAt()
	}
	if survey.GetDeletedAt() != nil {
		viewSurvey.DeletedAt = survey.GetDeletedAt()
	}

	return viewSurvey
}

func (sp surveysPresenter) respondentsToView(respondents entitySurvey.RespondentIDs) []*viewSurveys.SurveyRespondent {
	respondentsIDs := make([]*viewSurveys.SurveyRespondent, 0, len(respondents))
	for _, respondent := range respondents {
		respondentsIDs = append(respondentsIDs, &viewSurveys.SurveyRespondent{ID: uuid.UUID(respondent)})
	}

	return respondentsIDs
}

func (sp surveysPresenter) questionsToView(questions []*entitySurvey.Question) []*viewSurveys.SurveyQuestion {
	viewQuestions := make([]*viewSurveys.SurveyQuestion, 0, len(questions))
	for _, question := range questions {
		viewQuestion := &viewSurveys.SurveyQuestion{
			ID:             uuid.UUID(*question.GetID()),
			SurveyID:       uuid.UUID(*question.GetSurveyID()),
			Text:           question.Text,
			Type:           string(question.Type),
			IsRequired:     question.IsRequired,
			IsActive:       question.IsActive,
			Weight:         question.Weight,
			AnswerVariants: sp.answerVariantsToView(question.Answers),
		}
		if rules := question.GetRules(); rules != nil {
			newRules := &viewSurveys.SurveyQuestionRules{}
			if pickMinCount := rules.GetPickMinCount(); pickMinCount != nil {
				newRules.PickMinCount = pickMinCount
			}
			if pickMaxCount := rules.GetPickMaxCount(); pickMaxCount != nil {
				newRules.PickMaxCount = pickMaxCount
			}
			viewQuestion.Rules = newRules
		}
		if question.GetDeletedAt() != nil {
			viewQuestion.DeletedAt = question.DeletedAt
		}

		viewQuestions = append(viewQuestions, viewQuestion)
	}

	return viewQuestions
}

func (sp surveysPresenter) answerVariantsToView(answerVariants []*entitySurvey.AnswerVariant) []*viewSurveys.SurveyAnswerVariant {
	viewVariants := make([]*viewSurveys.SurveyAnswerVariant, 0, len(answerVariants))
	for _, variant := range answerVariants {
		contentType := int(variant.ContentType)
		newVariant := &viewSurveys.SurveyAnswerVariant{
			ID:          uuid.UUID(*variant.GetID()),
			QuestionID:  uuid.UUID(*variant.GetQuestionID()),
			Text:        variant.Text,
			WithContent: variant.WithContent,
			ContentType: &contentType,
			Weight:      variant.Weight,
		}
		if variant.GetDeletedAt() != nil {
			newVariant.DeletedAt = variant.GetDeletedAt()
		}

		if rules := variant.GetRules(); rules != nil {
			newVariant.Rules = &viewSurveys.SurveyAnswerVariantRules{
				ContentMinLength: rules.GetContentMinLength(),
				ContentMaxLength: rules.GetContentMaxLength(),
				ContentMinDigit:  rules.GetContentMinDigit(),
				ContentMaxDigit:  rules.GetContentMaxDigit(),
			}
		}

		if image := variant.GetImage(); image != nil {
			newImage := &viewSurveys.SurveyImage{
				ID: string(image.ID),
			}
			if extImage := image.GetExternalImageInfo(); extImage != nil {
				newImage.ImageExternalID = extImage.ID
				newImage.ImageExternalFileName = extImage.FileName
				newImage.ImageExternalURL = extImage.URL
				newImage.ImageExternalSize = int(extImage.Size)
			}
			newVariant.Image = newImage
		}

		viewVariants = append(viewVariants, newVariant)
	}

	return viewVariants
}

func (sp surveysPresenter) ToShortView(survey *entitySurvey.Survey) *viewSurveys.SurveyInfo {
	viewSurvey := &viewSurveys.SurveyInfo{
		Title:             survey.Title,
		Description:       survey.Description,
		ActivePeriodStart: survey.ActivePeriodStart,
		ActivePeriodEnd:   survey.ActivePeriodEnd,
		Questions:         sp.questionsToShortView(survey.Questions),
	}

	return viewSurvey
}

func (sp surveysPresenter) questionsToShortView(questions []*entitySurvey.Question) []*viewSurveys.SurveyQuestionInfo {
	viewQuestions := make([]*viewSurveys.SurveyQuestionInfo, 0, len(questions))
	for _, question := range questions {
		viewQuestion := &viewSurveys.SurveyQuestionInfo{
			ID:             uuid.UUID(*question.GetID()),
			Text:           question.Text,
			Type:           string(question.Type),
			IsRequired:     question.IsRequired,
			Weight:         question.Weight,
			AnswerVariants: sp.answerVariantsToShortView(question.Answers),
		}
		if rules := question.GetRules(); rules != nil {
			newRules := &viewSurveys.SurveyQuestionRules{}
			if pickMinCount := rules.GetPickMinCount(); pickMinCount != nil {
				newRules.PickMinCount = pickMinCount
			}
			if pickMaxCount := rules.GetPickMaxCount(); pickMaxCount != nil {
				newRules.PickMaxCount = pickMaxCount
			}
			viewQuestion.Rules = newRules
		}

		viewQuestions = append(viewQuestions, viewQuestion)
	}

	return viewQuestions
}

func (sp surveysPresenter) answerVariantsToShortView(answerVariants []*entitySurvey.AnswerVariant) []*viewSurveys.
	SurveyAnswerVariantInfo {
	viewVariants := make([]*viewSurveys.SurveyAnswerVariantInfo, 0, len(answerVariants))
	for _, variant := range answerVariants {
		contentType := int(variant.ContentType)
		newVariant := &viewSurveys.SurveyAnswerVariantInfo{
			ID:          uuid.UUID(*variant.GetID()),
			Text:        variant.Text,
			WithContent: variant.WithContent,
			ContentType: &contentType,
			Weight:      variant.Weight,
		}

		if rules := variant.GetRules(); rules != nil {
			newVariant.Rules = &viewSurveys.SurveyAnswerVariantRules{
				ContentMinLength: rules.GetContentMinLength(),
				ContentMaxLength: rules.GetContentMaxLength(),
				ContentMinDigit:  rules.GetContentMinDigit(),
				ContentMaxDigit:  rules.GetContentMaxDigit(),
			}
		}

		if variant.GetImage() != nil {
			image := &viewSurveys.SurveyImageInfo{
				ID: string(variant.GetImage().ID),
			}
			newVariant.Image = image
		}

		viewVariants = append(viewVariants, newVariant)
	}

	return viewVariants
}

func (sp surveysPresenter) IDsToEntities(ids []uuid.UUID) entitySurvey.SurveyIDs {
	surveyIDs := make(entitySurvey.SurveyIDs, 0, len(ids))
	for _, id := range ids {
		surveyIDs = append(surveyIDs, entitySurvey.SurveyID(id))
	}

	return surveyIDs
}

func (sp surveysPresenter) RespondentToEntity(respondent *viewSurveys.GetAllSurveysRespondent) *entitySurvey.SurveyRespondent {
	if respondent == nil {
		return &entitySurvey.SurveyRespondent{
			Type: entitySurvey.RespondentTypeAll,
		}
	}

	respondentsIDs := make(entitySurvey.RespondentIDs, 0, len(respondent.RespondentIDs))
	for _, respondent := range respondent.RespondentIDs {
		respondentsIDs = append(respondentsIDs, entitySurvey.RespondentID(respondent))
	}

	return &entitySurvey.SurveyRespondent{
		Type: sp.respondentTypeToEntity(respondent.RespondentType),
		Ids:  respondentsIDs,
	}
}

func (sp surveysPresenter) OptionsToEntity(options *viewSurveys.SurveysOptions) entitySurvey.SurveyFilterOptions {
	if options == nil {
		return entitySurvey.SurveyFilterOptions{}
	}

	return entitySurvey.SurveyFilterOptions{
		WithQuestions:         options.WithQuestions,
		WithAnswers:           options.WithAnswers,
		WithDeleted:           options.WithDeleted,
		WithInactiveQuestions: options.WithInactiveQuestions,
	}
}

func (sp surveysPresenter) PaginationToEntity(pagination *viewSurveys.GetAllSurveysPagination) entitySurvey.Pagination {
	result := entitySurvey.Pagination{
		Limit:    uint32(pagination.Limit),
		LastDate: pagination.LastDate,
	}

	if pagination.LastId != nil {
		lastID := entitySurvey.SurveyID(*pagination.LastId)
		result.LastId = &lastID
	}

	if pagination.Total != nil {
		result.Total = int64(*pagination.Total)
	}

	return result
}

func (sp surveysPresenter) SurveysWithPaginationToView(surveysWithPagination *entitySurvey.SurveysWithPagination) *viewSurveys.
	SurveysWithPagination {
	surveysResult := make([]*viewSurveys.Survey, 0, len(surveysWithPagination.Surveys))
	for _, survey := range surveysWithPagination.Surveys {
		surveysResult = append(surveysResult, sp.ToView(survey))
	}

	paginationResult := viewSurveys.SurveyPagination{
		Limit:    int(surveysWithPagination.Pagination.Limit),
		LastDate: surveysWithPagination.Pagination.LastDate,
		Total:    int(surveysWithPagination.Pagination.Total),
	}
	if surveysWithPagination.Pagination.LastId != nil {
		lastID := uuid.UUID(*surveysWithPagination.Pagination.LastId)
		paginationResult.LastId = &lastID
	}

	return &viewSurveys.SurveysWithPagination{
		Pagination: paginationResult,
		Surveys:    surveysResult,
	}
}

func (sp surveysPresenter) IDToView(ID *entitySurvey.SurveyID) *viewSurveys.IDResponse {
	return &viewSurveys.IDResponse{ID: ID.String()}
}
