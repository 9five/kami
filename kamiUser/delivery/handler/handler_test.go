package handler_test

import (
	// "context"
	"encoding/json"
	"fmt"
	"kami/domain"
	"kami/domain/mocks"
	kamiUserHandler "kami/kamiUser/delivery/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type KamiUserKamiUserHandlerSuite struct {
	suite.Suite
	kamiUserUsecase *mocks.KamiUserUsercase
	engine          *gin.Engine
	router          *gin.RouterGroup
	w               *httptest.ResponseRecorder
}

type KamiUserLoginHandlerSuite struct {
	suite.Suite
	kamiUserUsecase      *mocks.KamiUserUsercase
	twilioServiceUsecase *mocks.TwilioServiceUsecase
	encryptionUsecase    *mocks.EncryptionUsecase
	engine               *gin.Engine
	router               *gin.RouterGroup
	w                    *httptest.ResponseRecorder
}

func TestStart(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	suite.Run(t, &KamiUserKamiUserHandlerSuite{})
	suite.Run(t, &KamiUserLoginHandlerSuite{})
}

func mockmiddleware(ctx *gin.Context) {
	ctx.Set("id", uint(1))
	ctx.Set("phone", "0900123123")
}

func (s *KamiUserKamiUserHandlerSuite) SetupTest() {
	s.kamiUserUsecase = &mocks.KamiUserUsercase{}
	s.engine = gin.Default()
	s.router = s.engine.Group("/api/user", mockmiddleware)
	s.w = httptest.NewRecorder()
}

func (s *KamiUserLoginHandlerSuite) SetupTest() {
	s.kamiUserUsecase = new(mocks.KamiUserUsercase)
	s.twilioServiceUsecase = new(mocks.TwilioServiceUsecase)
	s.encryptionUsecase = new(mocks.EncryptionUsecase)
	s.engine = gin.Default()
	s.router = s.engine.Group("/api/login")
	s.w = httptest.NewRecorder()
}

func (s *KamiUserKamiUserHandlerSuite) TestGetCurrentUser() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockKamiUserOutput := &domain.KamiUserOutput{
		Phone:    "0900123123",
		Email:    "test@test.com",
		Birthday: "2000-01-01",
	}
	mockKamiUserOutputMarshal, _ := json.Marshal(mockKamiUserOutput)

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Model.ID == mockKamiUser.Model.ID })).
		Return(mockKamiUser, nil).Once()

	handler := kamiUserHandler.KamiUserHandler{
		KamiUserUsecase: s.kamiUserUsecase,
	}
	s.router.GET("/status", mockmiddleware, handler.GetCurrentUser)
	req, _ := http.NewRequest("GET", "/api/user/status", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockKamiUserOutputMarshal), s.w.Body.String())
}

func (s *KamiUserKamiUserHandlerSuite) TestUpdateCurrentUserInfo() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
	}

	mockKamiUserInput := &domain.KamiUserInput{
		Email:    "test@test.com",
		Gender:   "M",
		Birthday: "2000-01-01",
		Name:     "Test",
		Career:   "Test",
	}
	mockKamiUserInputMarshal, _ := json.Marshal(mockKamiUserInput)
	body := strings.NewReader(string(mockKamiUserInputMarshal))

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Model.ID == mockKamiUser.Model.ID })).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsecase.On("UpdateUserInfo", mock.Anything, mockKamiUser, mockKamiUserInput).
		Return(nil).Once()

	handler := kamiUserHandler.KamiUserHandler{
		KamiUserUsecase: s.kamiUserUsecase,
	}
	s.router.PUT("/updateInfo", mockmiddleware, handler.UpdateCurrentUserInfo)
	req, _ := http.NewRequest("PUT", "/api/user/updateInfo", body)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
}

func (s *KamiUserLoginHandlerSuite) TestEnterPhone() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockKamiUserLog := &domain.KamiUserLog{
		Model: gorm.Model{ID: 1},
		Phone: "0900123123",
	}

	mockEnterPhoneOutput := gin.H{
		"token":  "testingToken",
		"status": false,
	}
	mockEnterPhoneOutputMarshal, _ := json.Marshal(mockEnterPhoneOutput)

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsecase.On("CheckKamiUserLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log.Phone == mockKamiUserLog.Phone })).
		Return(nil).Once()
	s.kamiUserUsecase.On("GetKamiUserLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log.Phone == mockKamiUserLog.Phone })).
		Return(mockKamiUserLog, nil).Once()
	s.twilioServiceUsecase.On("SendVerificationSMS", mock.MatchedBy(func(phone string) bool { return phone == mockKamiUser.Phone })).
		Return(nil).Once()
	s.encryptionUsecase.On("DesEncrypt", mock.Anything, fmt.Sprintf("%s-%d", mockKamiUserLog.Phone, mockKamiUserLog.Model.ID)).
		Return("testingToken", nil).Once()

	handler := kamiUserHandler.LoginHandler{
		KamiUserUsecase:      s.kamiUserUsecase,
		TwilioServiceUsecase: s.twilioServiceUsecase,
		EncryptionUsecase:    s.encryptionUsecase,
	}
	s.router.POST("/enterPhone", handler.EnterPhone)
	req, _ := http.NewRequest("POST", "/api/login/enterPhone?phone=0900123123", nil)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.Suite.T(), string(mockEnterPhoneOutputMarshal), s.w.Body.String())
}

func (s *KamiUserLoginHandlerSuite) TestVerificationCheck() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	mockVerificationInput := domain.VerificationInput{
		ForgotPw: false,
		Token:    "testingToken",
		Code:     "123",
	}
	mockVerificationInputMarshal, _ := json.Marshal(mockVerificationInput)
	body := strings.NewReader(string(mockVerificationInputMarshal))

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(mockKamiUser, nil).Once()
	s.twilioServiceUsecase.On("VerificationCheck", mockKamiUser.Phone, mockVerificationInput.Code).
		Return(nil).Once()
	s.encryptionUsecase.On("DesDecrypt", mock.Anything, "testingToken").
		Return("0900123123-1", nil).Once()

	handler := kamiUserHandler.LoginHandler{
		KamiUserUsecase:      s.kamiUserUsecase,
		TwilioServiceUsecase: s.twilioServiceUsecase,
		EncryptionUsecase:    s.encryptionUsecase,
	}
	s.router.POST("/verificationCheck", handler.VerificationCheck)
	req, _ := http.NewRequest("POST", "/api/login/verificationCheck", body)
	s.engine.ServeHTTP(s.w, req)

	assert.Equal(s.Suite.T(), http.StatusOK, s.w.Code)
}

func (s *KamiUserLoginHandlerSuite) TestEnterPassword() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		Password: "testingPassword",
	}

	mockKamiUserLog := &domain.KamiUserLog{
		Model: gorm.Model{ID: 1},
		Phone: "0900123123",
	}

	s.kamiUserUsecase.On("LoginKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool {
		return user.Phone == mockKamiUser.Phone && user.Password == mockKamiUser.Password
	})).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsecase.On("GetKamiUserLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log.Phone == mockKamiUserLog.Phone })).
		Return(mockKamiUserLog, nil).Once()
	s.kamiUserUsecase.On("UpdateKamiUserLog", mock.Anything, mockKamiUserLog).
		Return(nil).Once()
	s.kamiUserUsecase.On("GenerateToken", mock.Anything, mockKamiUser).
		Return("testingJWTToken", nil).Once()
	s.encryptionUsecase.On("DesDecrypt", mock.Anything, "testingToken").
		Return("0900123123-1", nil).Once()
	s.encryptionUsecase.On("HashPassword", mock.Anything, mockKamiUser.Password).
		Return("testingHashPassword", nil).Once()
	s.encryptionUsecase.On("CheckPwHash", mock.Anything, mockKamiUser.Password, "testingHashPassword").
		Return(true).Once()

	handler := kamiUserHandler.LoginHandler{
		KamiUserUsecase:      s.kamiUserUsecase,
		TwilioServiceUsecase: s.twilioServiceUsecase,
		EncryptionUsecase:    s.encryptionUsecase,
	}
	s.router.POST("/enterPassword", handler.EnterPassword)
	req, _ := http.NewRequest("POST", "/api/login/enterPassword?token=testingToken&password=testingPassword", nil)
	s.engine.ServeHTTP(s.w, req)
}

func (s *KamiUserLoginHandlerSuite) TestForgotPassword() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Birthday: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		Password: "testingPassword",
	}

	mockKamiUserLog := &domain.KamiUserLog{
		Model: gorm.Model{ID: 1},
		Phone: "0900123123",
	}

	s.kamiUserUsecase.On("GetKamiUser", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(mockKamiUser, nil).Once()
	s.kamiUserUsecase.On("CheckKamiUserLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log.Phone == mockKamiUserLog.Phone })).
		Return(nil).Once()
	s.twilioServiceUsecase.On("SendVerificationSMS", mock.MatchedBy(func(phone string) bool { return phone == mockKamiUser.Phone })).
		Return(nil).Once()
	s.encryptionUsecase.On("DesDecrypt", mock.Anything, "testingToken").
		Return("0900123123-1", nil).Once()

	handler := kamiUserHandler.LoginHandler{
		KamiUserUsecase:      s.kamiUserUsecase,
		TwilioServiceUsecase: s.twilioServiceUsecase,
		EncryptionUsecase:    s.encryptionUsecase,
	}
	s.router.POST("/forgotPassword", handler.ForgotPassword)
	req, _ := http.NewRequest("POST", "/api/login/forgotPassword?token=testingToken", nil)
	s.engine.ServeHTTP(s.w, req)
}
