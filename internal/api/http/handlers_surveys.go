package http

import (
	"context"
	"errors"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
)

type surveysHandlers struct {
	surveysInteractor SurveysSurveysInteractor
	surveysPresenter  SurveysPresenter
	answersInteractor SurveysAnswersInteractor
	answersPresenter  SurveysAnswersPresenter
	imagesInteractor  SurveysImagesInteractor
	imagesPresenter   SurveysImagesPresenter
	logger            ditzap.Logger
}

func NewSurveysHandlers(
	surveysInteractor SurveysSurveysInteractor,
	surveysPresenter SurveysPresenter,
	answersInteractor SurveysAnswersInteractor,
	answersPresenter SurveysAnswersPresenter,
	imagesInteractor SurveysImagesInteractor,
	imagesPresenter SurveysImagesPresenter,
	logger ditzap.Logger) *surveysHandlers {
	return &surveysHandlers{
		surveysInteractor: surveysInteractor,
		surveysPresenter:  surveysPresenter,
		answersInteractor: answersInteractor,
		answersPresenter:  answersPresenter,
		imagesInteractor:  imagesInteractor,
		imagesPresenter:   imagesPresenter,
		logger:            logger,
	}
}

// @Summary Получение опроса по идентификатору
// @Description Выдаётся опрос по ID.
// @Tags     Опросы
// @Produce  json
// @Param     id path string true "Survey ID"
// @Router   /survey/{id} [get]
// @Success  200 {object} view.SurveyInfo "Опрос"
// @Failure  401,404,500 {object} ErrorResponse
func (sh surveysHandlers) getSurvey(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, RequestTimeOut)
	defer cancelCtx()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sh.logger.Debug("can't parse param id into uuid", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	result, err := sh.surveysInteractor.Get(
		ctx,
		surveys.SurveyID(id),
		surveys.SurveyFilterOptions{
			WithQuestions:         true,
			WithAnswers:           true,
			WithDeleted:           false,
			WithInactiveQuestions: false,
		},
	)
	if err != nil {
		if errors.Is(err, diterrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))

			return
		} else if errors.As(err, new(diterrors.ValidationError)) {
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(err))

			return
		}

		sh.logger.Error("can't get survey", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, sh.surveysPresenter.ToShortView(result))
}

// @Summary Создание ответов на опрос
// @Description Добавляет новые ответы на опрос
// @Tags     Опросы
// @Produce  json
// @Param    answers body view.NewSurveyAnswers true "ответы"
// @Router   /survey/answers [post]
// @Success  201 {array} view.SurveyAnswerInfo "Список идентификаторов ответов"
// @Failure  400,500 {object} ErrorResponse
func (sh surveysHandlers) addAnswers(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, RequestTimeOut)
	defer cancelCtx()

	var answers viewSurveys.NewSurveyAnswers
	if err := c.ShouldBindJSON(&answers); err != nil {
		sh.logger.Debug("can't unbind NewSurveyAnswers json", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	result, err := sh.answersInteractor.Add(ctx, sh.answersPresenter.ToNewEntities(&answers))
	if err != nil {
		if errors.As(err, new(diterrors.ValidationError)) {
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(err))
			return
		}

		sh.logger.Error("can't add answers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusCreated, sh.answersPresenter.ToShortViews(result))
}

// @Summary Получение изображения по идентификатору
// @Description Выдаётся изображение по ID.
// @Tags     Опросы
// @Produce  jpeg
// @Produce  png
// @Param    id path string true "Image ID"
// @Router   /survey/images/{id} [get]
// @Success  200 {file} file "Изображение"
// @Failure  400,500 {object} ErrorResponse
func (sh surveysHandlers) getImage(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, RequestTimeOut)
	defer cancelCtx()

	imageID := c.Param("id")

	result, err := sh.imagesInteractor.Get(ctx, imageID)
	if err != nil {
		sh.logger.Error("can't get image", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))

		return
	}

	c.Data(http.StatusOK, "application/octet-stream", result)
}
