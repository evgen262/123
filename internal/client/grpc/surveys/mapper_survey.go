package surveys

import (
	"fmt"

	answervariantv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answervariant/v1"
	questionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/question/v1"
	respondentv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/respondent/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/shared/v1"
	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var questionTypeEntityToPbDictionary = map[surveys.QuestionType]questionv1.Question_QuestionType{
	surveys.QuestionTypeText:        questionv1.Question_QUESTION_TYPE_TEXT,
	surveys.QuestionTypeCheckboxImg: questionv1.Question_QUESTION_TYPE_CHECKBOX_IMG,
	surveys.QuestionTypeCheckbox:    questionv1.Question_QUESTION_TYPE_CHECKBOX,
	surveys.QuestionTypeRadioImg:    questionv1.Question_QUESTION_TYPE_RADIO_IMG,
	surveys.QuestionTypeRadio:       questionv1.Question_QUESTION_TYPE_RADIO,
	surveys.QuestionTypeNumber:      questionv1.Question_QUESTION_TYPE_INPUT_NUMBER,
}

var contentTypeToPbDictionary = map[surveys.ContentType]answervariantv1.AnswerVariant_ContentType{
	surveys.ContentTypeDigit: answervariantv1.AnswerVariant_CONTENT_TYPE_DIGIT,
	surveys.ContentTypeText:  answervariantv1.AnswerVariant_CONTENT_TYPE_TEXT,
}

var questionTypeEntityToAddRequestDictionary = map[surveys.QuestionType]surveyv1.AddRequest_Survey_Question_QuestionType{
	surveys.QuestionTypeText:        surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_TEXT,
	surveys.QuestionTypeNumber:      surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_INPUT_NUMBER,
	surveys.QuestionTypeRadioImg:    surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_RADIO_IMG,
	surveys.QuestionTypeRadio:       surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_RADIO,
	surveys.QuestionTypeCheckboxImg: surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_CHECKBOX_IMG,
	surveys.QuestionTypeCheckbox:    surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_CHECKBOX,
}

var contentTypeAddToEntityDictionary = map[surveys.ContentType]surveyv1.
	AddRequest_Survey_Question_AnswerVariant_ContentType{
	surveys.ContentTypeDigit: surveyv1.AddRequest_Survey_Question_AnswerVariant_CONTENT_TYPE_DIGIT,
	surveys.ContentTypeText:  surveyv1.AddRequest_Survey_Question_AnswerVariant_CONTENT_TYPE_TEXT,
}

var questionTypePbToEntityDictionary = map[questionv1.Question_QuestionType]surveys.QuestionType{
	questionv1.Question_QUESTION_TYPE_TEXT:         surveys.QuestionTypeText,
	questionv1.Question_QUESTION_TYPE_INPUT_NUMBER: surveys.QuestionTypeNumber,
	questionv1.Question_QUESTION_TYPE_RADIO_IMG:    surveys.QuestionTypeRadioImg,
	questionv1.Question_QUESTION_TYPE_RADIO:        surveys.QuestionTypeRadio,
	questionv1.Question_QUESTION_TYPE_CHECKBOX_IMG: surveys.QuestionTypeCheckboxImg,
	questionv1.Question_QUESTION_TYPE_CHECKBOX:     surveys.QuestionTypeCheckbox,
}

var contentTypeDictionaryToEntity = map[answervariantv1.AnswerVariant_ContentType]surveys.ContentType{
	answervariantv1.AnswerVariant_CONTENT_TYPE_DIGIT: surveys.ContentTypeDigit,
	answervariantv1.AnswerVariant_CONTENT_TYPE_TEXT:  surveys.ContentTypeText,
}

type surveyMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewSurveyMapper(timeUtils timeUtils.TimeUtils) *surveyMapper {
	return &surveyMapper{timeUtils: timeUtils}
}

func (sm surveyMapper) PaginationToEntity(response *sharedv1.PaginationResponse) (*surveys.Pagination, error) {
	pagination := &surveys.Pagination{
		Limit:    response.GetLimit(),
		LastDate: sm.timeUtils.TimestampToTime(response.GetLastCreatedTime()),
		Total:    response.GetTotal(),
	}

	if response.GetLastId() != nil {
		parsedId, err := uuid.Parse(response.GetLastId().GetValue())
		if err != nil {
			return nil, fmt.Errorf("can't parse last survey id in pagination options: %w", err)
		}
		id := surveys.SurveyID(parsedId)
		pagination.LastId = &id
	}

	return pagination, nil
}

func (sm surveyMapper) PaginationToPb(pagination *surveys.Pagination) *sharedv1.PaginationRequest {
	if pagination == nil {
		return nil
	}
	req := &sharedv1.PaginationRequest{
		Limit:           pagination.Limit,
		LastCreatedTime: sm.timeUtils.TimeToTimestamp(pagination.LastDate),
	}

	if pagination.LastId != nil {
		req.LastId = &wrapperspb.StringValue{
			Value: pagination.LastId.String(),
		}
	}

	return req
}

func (sm surveyMapper) RespondentIDsToEntity(respondentIDs []string) (surveys.RespondentIDs, error) {
	idsArr := make([]surveys.RespondentID, 0, len(respondentIDs))
	for _, respondentID := range respondentIDs {
		id, err := uuid.Parse(respondentID)
		if err != nil {
			return nil, fmt.Errorf("can't parse uuid: %w", err)
		}
		idsArr = append(idsArr, surveys.RespondentID(id))
	}

	return idsArr, nil
}

func (sm surveyMapper) RespondentToPb(respondent *surveys.SurveyRespondent) *respondentv1.Respondent {
	if respondent == nil {
		return nil
	}
	result := &respondentv1.Respondent{}
	switch respondent.Type {
	case surveys.RespondentTypeUnknown:
		result.Respondent = &respondentv1.Respondent_Unknown_{}
	case surveys.RespondentTypeAnonymous:
		result.Respondent = &respondentv1.Respondent_Anonymous_{}
	case surveys.RespondentTypeUser:
		result.Respondent = &respondentv1.Respondent_User_{
			User: &respondentv1.Respondent_User{
				Ids: respondent.Ids.ToStringSlice(),
			},
		}
	}

	return result
}

func (sm surveyMapper) RespondentToEntity(respondent *respondentv1.Respondent) (*surveys.SurveyRespondent, error) {
	if respondent.GetRespondent() == nil {
		return nil, nil
	}
	result := &surveys.SurveyRespondent{}
	switch respondent.GetRespondent().(type) {
	case *respondentv1.Respondent_Unknown_:
		result.Type = surveys.RespondentTypeUnknown
	case *respondentv1.Respondent_Anonymous_:
		result.Type = surveys.RespondentTypeAnonymous
	case *respondentv1.Respondent_User_:
		result.Type = surveys.RespondentTypeUser
		rIds, err := sm.RespondentIDsToEntity(respondent.GetUser().GetIds())
		if err != nil {
			return nil, err
		}
		result.Ids = rIds
	}

	return result, nil
}

func (sm surveyMapper) OptionsToPb(options *surveys.SurveyFilterOptions) *sharedv1.Options {
	if options == nil {
		return nil
	}
	result := &sharedv1.Options{
		WithQuestions:         options.WithQuestions,
		WithInactiveQuestions: options.WithInactiveQuestions,
		WithAnswers:           options.WithAnswers,
		WithDeleted:           options.WithDeleted,
	}

	return result
}

func (sm surveyMapper) SurveyToPb(survey *surveys.Survey) *surveyv1.Survey {
	if survey == nil {
		return nil
	}
	newSurvey := &surveyv1.Survey{
		Id:                    survey.GetID().String(),
		Title:                 survey.Title,
		Description:           survey.Description,
		ActivePeriodStartTime: sm.timeUtils.TimeToTimestamp(&survey.ActivePeriodStart),
		ActivePeriodEndTime:   sm.timeUtils.TimeToTimestamp(&survey.ActivePeriodEnd),
		Author:                survey.Author,
		IsPublished:           survey.IsPublished,
		CreatedTime:           sm.timeUtils.TimeToTimestamp(survey.GetCreatedAt()),
		UpdatedTime:           sm.timeUtils.TimeToTimestamp(survey.GetUpdatedAt()),
		DeletedTime:           sm.timeUtils.TimeToTimestamp(survey.GetDeletedAt()),
		Respondent:            sm.RespondentToPb(survey.GetRespondent()),
		Questions:             sm.questionsToPb(survey.Questions),
	}

	return newSurvey
}

func (sm surveyMapper) questionsToPb(questions []*surveys.Question) []*questionv1.Question {
	mappedQuestions := make([]*questionv1.Question, 0, len(questions))
	for _, question := range questions {
		newQuestion := &questionv1.Question{
			Id:             question.GetID().String(),
			SurveyId:       question.GetSurveyID().String(),
			Text:           question.Text,
			QuestionType:   sm.questionTypeToPb(question.Type),
			IsRequired:     question.IsRequired,
			IsActive:       question.IsActive,
			Weight:         int32(question.Weight),
			DeletedTime:    sm.timeUtils.TimeToTimestamp(question.GetDeletedAt()),
			AnswerVariants: sm.answerVariantsToPb(question.Answers),
		}
		if rules := question.GetRules(); rules != nil {
			newRules := &questionv1.Question_Rules{}
			if pickMinCount := rules.GetPickMinCount(); pickMinCount != nil {
				newRules.PickMinCount = &wrapperspb.Int64Value{Value: int64(*pickMinCount)}
			}
			if pickMaxCount := rules.GetPickMaxCount(); pickMaxCount != nil {
				newRules.PickMaxCount = &wrapperspb.Int64Value{Value: int64(*pickMaxCount)}
			}
			newQuestion.Rules = newRules
		}
		mappedQuestions = append(mappedQuestions, newQuestion)
	}

	return mappedQuestions
}

func (sm surveyMapper) questionTypeToPb(questionType surveys.QuestionType) questionv1.Question_QuestionType {
	value, ok := questionTypeEntityToPbDictionary[questionType]
	if !ok {
		return questionv1.Question_QUESTION_TYPE_INVALID
	}

	return value
}

func (sm surveyMapper) answerVariantsToPb(answerVariants []*surveys.AnswerVariant) []*answervariantv1.AnswerVariant {
	mappedVariants := make([]*answervariantv1.AnswerVariant, 0, len(answerVariants))
	for _, variant := range answerVariants {
		newVariant := &answervariantv1.AnswerVariant{
			Id:          variant.GetID().String(),
			QuestionId:  variant.GetQuestionID().String(),
			Text:        variant.Text,
			WithContent: variant.WithContent,
			ContentType: sm.contentTypeToPb(variant.ContentType),
			DeletedTime: sm.timeUtils.TimeToTimestamp(variant.GetDeletedAt()),
			Weight:      int32(variant.Weight),
		}

		if rules := variant.GetRules(); variant.WithContent && rules != nil {
			newRules := &answervariantv1.AnswerVariant_Rules{}
			if rules.ContentMinLength != nil {
				newRules.ContentMinLength = &wrapperspb.Int64Value{Value: int64(*rules.ContentMinLength)}
			}
			if rules.ContentMaxLength != nil {
				newRules.ContentMaxLength = &wrapperspb.Int64Value{Value: int64(*rules.ContentMaxLength)}
			}
			if rules.ContentMinDigit != nil {
				newRules.ContentMinDigit = &wrapperspb.Int64Value{Value: int64(*rules.ContentMinDigit)}
			}
			if rules.ContentMaxDigit != nil {
				newRules.ContentMaxDigit = &wrapperspb.Int64Value{Value: int64(*rules.ContentMaxDigit)}
			}
			newVariant.Rules = newRules
		}

		if image := variant.GetImage(); image != nil {
			newImage := &answervariantv1.AnswerVariant_Image{
				Id: string(image.ID),
			}
			if extImage := image.GetExternalImageInfo(); extImage != nil {
				newImage.ExternalImageInfo = &answervariantv1.AnswerVariant_Image_ExternalInfo{
					Id:       extImage.ID.String(),
					Filename: extImage.FileName,
					Url:      extImage.URL,
					Size:     extImage.Size,
				}
			}
			newVariant.Image = newImage
		}
		mappedVariants = append(mappedVariants, newVariant)
	}

	return mappedVariants
}

func (sm surveyMapper) contentTypeToPb(contentType surveys.ContentType) answervariantv1.
	AnswerVariant_ContentType {
	value, ok := contentTypeToPbDictionary[contentType]
	if !ok {
		return answervariantv1.AnswerVariant_CONTENT_TYPE_INVALID
	}

	return value
}

func (sm surveyMapper) NewSurveyToPb(survey *surveys.Survey) *surveyv1.AddRequest_Survey {
	if survey == nil {
		return nil
	}
	newSurvey := &surveyv1.AddRequest_Survey{
		Title:                 survey.Title,
		Description:           survey.Description,
		ActivePeriodStartTime: sm.timeUtils.TimeToTimestamp(&survey.ActivePeriodStart),
		ActivePeriodEndTime:   sm.timeUtils.TimeToTimestamp(&survey.ActivePeriodEnd),
		Author:                survey.Author,
		IsPublished:           survey.IsPublished,
		Respondent:            sm.RespondentToPb(survey.Respondent),
		Questions:             sm.newQuestionsToPb(survey.Questions),
	}

	return newSurvey
}

func (sm surveyMapper) newQuestionsToPb(questions []*surveys.Question) []*surveyv1.AddRequest_Survey_Question {
	mappedQuestions := make([]*surveyv1.AddRequest_Survey_Question, 0, len(questions))
	for _, question := range questions {
		newQuestion := &surveyv1.AddRequest_Survey_Question{
			Text:           question.Text,
			QuestionType:   sm.newQuestionTypeToPb(question.Type),
			IsRequired:     question.IsRequired,
			IsActive:       question.IsActive,
			Weight:         int32(question.Weight),
			AnswerVariants: sm.newAnswerVariantsToPb(question.Answers),
		}
		if rules := question.GetRules(); rules != nil {
			newRules := &surveyv1.AddRequest_Survey_Question_Rules{}
			if pickMinCount := rules.GetPickMinCount(); pickMinCount != nil {
				newRules.PickMinCount = &wrapperspb.Int64Value{Value: int64(*pickMinCount)}
			}
			if pickMaxCount := rules.GetPickMaxCount(); pickMaxCount != nil {
				newRules.PickMaxCount = &wrapperspb.Int64Value{Value: int64(*pickMaxCount)}
			}
			newQuestion.Rules = newRules
		}
		mappedQuestions = append(mappedQuestions, newQuestion)
	}

	return mappedQuestions
}

func (sm surveyMapper) newQuestionTypeToPb(questionType surveys.
	QuestionType) surveyv1.
	AddRequest_Survey_Question_QuestionType {
	value, ok := questionTypeEntityToAddRequestDictionary[questionType]
	if !ok {
		return surveyv1.AddRequest_Survey_Question_QUESTION_TYPE_INVALID
	}

	return value
}

func (sm surveyMapper) newAnswerVariantsToPb(answers []*surveys.AnswerVariant) []*surveyv1.
	AddRequest_Survey_Question_AnswerVariant {
	mappedVariants := make([]*surveyv1.AddRequest_Survey_Question_AnswerVariant, 0, len(answers))
	for _, variant := range answers {
		newVariant := &surveyv1.AddRequest_Survey_Question_AnswerVariant{
			Text:        variant.Text,
			ContentType: sm.newContentTypeToPb(variant.ContentType),
			WithContent: variant.WithContent,
			Weight:      int32(variant.Weight),
		}

		if rules := variant.GetRules(); rules != nil && variant.WithContent {
			newRules := &surveyv1.AddRequest_Survey_Question_AnswerVariant_Rules{}
			if value := rules.GetContentMinLength(); value != nil {
				newRules.ContentMinLength = &wrapperspb.Int64Value{
					Value: int64(*value),
				}
			}
			if value := rules.GetContentMaxLength(); value != nil {
				newRules.ContentMaxLength = &wrapperspb.Int64Value{
					Value: int64(*value),
				}
			}
			if value := rules.GetContentMinDigit(); value != nil {
				newRules.ContentMinDigit = &wrapperspb.Int64Value{
					Value: int64(*value),
				}
			}
			if value := rules.GetContentMaxDigit(); value != nil {
				newRules.ContentMaxDigit = &wrapperspb.Int64Value{
					Value: int64(*value),
				}
			}
			newVariant.Rules = newRules
		}

		if image := variant.GetImage(); image != nil {
			newImage := &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image{
				Id: string(image.ID),
			}

			if extImage := image.GetExternalImageInfo(); extImage != nil {
				newImage.ExternalImageInfo = &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image_ExternalInfo{
					Id:       extImage.ID.String(),
					Filename: extImage.FileName,
					Url:      extImage.URL,
					Size:     extImage.Size,
				}
			}
			newVariant.Image = newImage
		}
		mappedVariants = append(mappedVariants, newVariant)
	}

	return mappedVariants
}

func (sm surveyMapper) newContentTypeToPb(contentType surveys.ContentType) surveyv1.
	AddRequest_Survey_Question_AnswerVariant_ContentType {
	value, ok := contentTypeAddToEntityDictionary[contentType]
	if !ok {
		return surveyv1.AddRequest_Survey_Question_AnswerVariant_CONTENT_TYPE_INVALID
	}

	return value
}

func (sm surveyMapper) SurveyToEntity(surveyPb *surveyv1.Survey) (*surveys.Survey, error) {
	surveyId, err := uuid.Parse(surveyPb.GetId())
	if err != nil {
		return nil, fmt.Errorf("can't parse survey id: %w", err)
	}

	sID := surveys.SurveyID(surveyId)

	newSurvey := &surveys.Survey{
		ID:                &sID,
		Title:             surveyPb.GetTitle(),
		Description:       surveyPb.GetDescription(),
		ActivePeriodStart: surveyPb.GetActivePeriodStartTime().AsTime(),
		ActivePeriodEnd:   surveyPb.GetActivePeriodEndTime().AsTime(),
		Author:            surveyPb.GetAuthor(),
		IsPublished:       surveyPb.GetIsPublished(),
		CreatedAt:         sm.timeUtils.TimestampToTime(surveyPb.GetCreatedTime()),
		UpdatedAt:         sm.timeUtils.TimestampToTime(surveyPb.GetUpdatedTime()),
		DeletedAt:         sm.timeUtils.TimestampToTime(surveyPb.GetDeletedTime()),
	}

	respondent, err := sm.RespondentToEntity(surveyPb.GetRespondent())
	if err != nil {
		return nil, fmt.Errorf("can't convert respondent to entity: %w", err)
	}
	newSurvey.Respondent = respondent

	questions, err := sm.questionsToEntity(surveyPb.GetQuestions())
	if err != nil {
		return nil, fmt.Errorf("can't convert questions to entity: %w", err)
	}
	newSurvey.Questions = questions

	return newSurvey, nil
}

func (sm surveyMapper) questionsToEntity(questions []*questionv1.Question) ([]*surveys.Question, error) {
	mappedQuestions := make([]*surveys.Question, 0, len(questions))
	for _, question := range questions {
		questionId, err := uuid.Parse(question.GetId())
		if err != nil {
			return nil, fmt.Errorf("can't parse question id: %w", err)
		}
		questionSurveyId, err := uuid.Parse(question.GetSurveyId())
		if err != nil {
			return nil, fmt.Errorf("can't parse survey id: %w", err)
		}

		qID := surveys.QuestionID(questionId)
		sID := surveys.SurveyID(questionSurveyId)
		newQuestion := &surveys.Question{
			ID:         &qID,
			SurveyID:   &sID,
			Text:       question.GetText(),
			Type:       sm.questionTypeToEntity(question.GetQuestionType()),
			IsRequired: question.GetIsRequired(),
			IsActive:   question.GetIsActive(),
			Weight:     int(question.GetWeight()),
			DeletedAt:  sm.timeUtils.TimestampToTime(question.GetDeletedTime()),
		}
		if rules := question.GetRules(); rules != nil {
			newRules := &surveys.QuestionRules{}
			if pickMinCount := rules.GetPickMinCount(); pickMinCount != nil {
				value := int(pickMinCount.GetValue())
				newRules.PickMinCount = &value
			}
			if pickMaxCount := rules.GetPickMaxCount(); pickMaxCount != nil {
				value := int(pickMaxCount.GetValue())
				newRules.PickMaxCount = &value
			}
			newQuestion.Rules = newRules
		}

		variants, err := sm.answerVariantsToEntity(question.GetAnswerVariants())
		if err != nil {
			return nil, fmt.Errorf("can't convert answers variants to entity: %w", err)
		}
		newQuestion.Answers = variants
		mappedQuestions = append(mappedQuestions, newQuestion)
	}

	return mappedQuestions, nil
}

func (sm surveyMapper) questionTypeToEntity(questionType questionv1.Question_QuestionType) surveys.
	QuestionType {
	value, ok := questionTypePbToEntityDictionary[questionType]
	if !ok {
		return surveys.QuestionTypeInvalid
	}

	return value
}

func (sm surveyMapper) answerVariantsToEntity(variants []*answervariantv1.AnswerVariant) ([]*surveys.
	AnswerVariant, error) {
	mappedVariants := make([]*surveys.AnswerVariant, 0, len(variants))
	for _, variant := range variants {
		variantId, err := uuid.Parse(variant.GetId())
		if err != nil {
			return nil, fmt.Errorf("can't parse variant id: %w", err)
		}
		variantQuestionId, err := uuid.Parse(variant.GetQuestionId())
		if err != nil {
			return nil, fmt.Errorf("can't parse question id: %w", err)
		}

		avID := surveys.AnswerID(variantId)
		qID := surveys.QuestionID(variantQuestionId)
		newVariant := &surveys.AnswerVariant{
			ID:          &avID,
			QuestionID:  &qID,
			Text:        variant.GetText(),
			WithContent: variant.WithContent,
			ContentType: sm.contentTypeToEntity(variant.GetContentType()),
			DeletedAt:   nil,
			Weight:      int(variant.GetWeight()),
		}

		if variant.GetRules() != nil {
			newVariant.Rules = &surveys.AnswerRules{}
			if variant.GetRules().GetContentMinLength() != nil {
				minLength := int(variant.GetRules().GetContentMinLength().GetValue())
				newVariant.Rules.ContentMinLength = &minLength
			}
			if variant.GetRules().GetContentMaxLength() != nil {
				maxLength := int(variant.GetRules().GetContentMaxLength().GetValue())
				newVariant.Rules.ContentMaxLength = &maxLength
			}
			if variant.GetRules().GetContentMinDigit() != nil {
				minDigit := int(variant.GetRules().GetContentMinDigit().GetValue())
				newVariant.Rules.ContentMinDigit = &minDigit
			}
			if variant.GetRules().GetContentMaxDigit() != nil {
				maxDigit := int(variant.GetRules().GetContentMaxDigit().GetValue())
				newVariant.Rules.ContentMaxDigit = &maxDigit
			}
		}

		if variant.GetImage() != nil {
			newImage := &surveys.Image{
				ID: surveys.ImageID(variant.GetImage().GetId()),
			}
			if variant.GetImage().GetExternalImageInfo() != nil {
				externalInfoId, err := uuid.Parse(variant.GetImage().GetExternalImageInfo().GetId())
				if err != nil {
					return nil, fmt.Errorf("can't parse image id: %w", err)
				}
				newImage.ExternalImageInfo = &surveys.ExternalProperties{
					ID:       externalInfoId,
					FileName: variant.GetImage().GetExternalImageInfo().GetFilename(),
					URL:      variant.GetImage().GetExternalImageInfo().GetUrl(),
					Size:     variant.GetImage().GetExternalImageInfo().GetSize(),
				}
			}
			newVariant.Image = newImage
		}

		mappedVariants = append(mappedVariants, newVariant)
	}

	return mappedVariants, nil
}

func (sm surveyMapper) contentTypeToEntity(contentType answervariantv1.AnswerVariant_ContentType) surveys.ContentType {
	value, ok := contentTypeDictionaryToEntity[contentType]
	if !ok {
		return surveys.ContentTypeInvalid
	}

	return value
}

func (sm surveyMapper) SurveysToEntities(surveysArr []*surveyv1.Survey) ([]*surveys.Survey, error) {
	newSurveys := make([]*surveys.Survey, 0, len(surveysArr))
	for _, surveyPb := range surveysArr {
		surveyEntity, err := sm.SurveyToEntity(surveyPb)
		if err != nil {
			return nil, fmt.Errorf("can't convert survey to entity: %w", err)
		}
		newSurveys = append(newSurveys, surveyEntity)
	}

	return newSurveys, nil
}
