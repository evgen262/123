package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type RouterSuite struct {
	suite.Suite

	db                *sqlx.DB
	middlewareOptions []*MiddlewareOption
	logger            *ditzap.MockLogger
	router            *router
}

func (s *RouterSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())

	db, _, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.db = sqlx.NewDb(db, "sqlmock")
	s.logger = ditzap.NewMockLogger(ctrl)
	s.middlewareOptions = []*MiddlewareOption{}

	gin.SetMode(gin.TestMode)
	s.router = NewRouter(
		EnvironmentTest,
		s.logger,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		s.middlewareOptions...,
	)
}

func (s *RouterSuite) TearDownAllSuite() {
	err := s.db.Close()
	s.NoError(err)
}

func (s *RouterSuite) TestRegisterRoutes() {
	tests := []struct {
		name     string
		prepare  func()
		testFunc func()
	}{
		{
			name: "handlers not init",
			prepare: func() {
				s.logger.EXPECT().Error("handlers not init")
			},
			testFunc: func() {
				ginTest := gintest.NewGinTest()
				s.router.engine = ginTest.GetEngine()
				s.router.registerRoutes()
			},
		},
		{
			name: "userGroup route test",
			prepare: func() {
				s.router.initHandlers()
			},
			testFunc: func() {
				testCases := []*gintest.RouterTestCase{}

				ginTest := gintest.NewGinTest()
				s.router.engine = ginTest.GetEngine()
				s.router.registerRoutes()
				ginTest.SetEngine(s.router.engine).TestRouterRoutes(s.T(), testCases...)
			},
		},
	}

	for _, tt := range tests {
		tt.prepare()
		s.Run(tt.name, tt.testFunc)
	}
}

func (s *RouterSuite) TestRecovery() {
	testRouter := gin.New()
	testWriter := httptest.NewRecorder()
	recovered := struct {
		testField string
	}{
		testField: "recovery testField",
	}
	s.logger.EXPECT().Error("http server recovery", zap.Error(fmt.Errorf("%s", recovered)))

	s.router.recovery(gin.CreateTestContextOnly(testWriter, testRouter), recovered)
	assert.Equal(s.T(), http.StatusInternalServerError, testWriter.Code)
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, &RouterSuite{})
}
