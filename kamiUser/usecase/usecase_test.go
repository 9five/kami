package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"kami/domain"
	"kami/domain/mocks"
	ucase "kami/kamiUser/usecase"
)

type KamiUserUcaseSuite struct {
	suite.Suite
	repo *mocks.KamiUserRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &KamiUserUcaseSuite{})
}

func (s *KamiUserUcaseSuite) SetupTest() {
	s.repo = new(mocks.KamiUserRepository)
}

func (s *KamiUserUcaseSuite) TestNewKamiUser_Success() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	s.repo.On("New", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(mockKamiUser, nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	res, err := u.NewKamiUser(context.TODO(), mockKamiUser)
	assert.Nil(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserUcaseSuite) TestGetKamiUser_Success() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(mockKamiUser, nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	res, err := u.GetKamiUser(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserUcaseSuite) TestUpdateKamiUser_Success() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	s.repo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.UpdateKamiUser(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserUcaseSuite) TestGetKamiUserLog_Success() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	s.repo.On("GetLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(mockKamiUserLog, nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	res, err := u.GetKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUserLog, res)
}

func (s *KamiUserUcaseSuite) TestUpdateKamiUserLog_Success() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	s.repo.On("UpdateLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.UpdateKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserUcaseSuite) TestLoginKamiUser_Success_NewUser() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(&domain.KamiUser{}, nil).Once()
	s.repo.On("New", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(mockKamiUser, nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	res, err := u.LoginKamiUser(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserUcaseSuite) TestLoginKamiUser_Success_UpdatePassword() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user.Phone == mockKamiUser.Phone })).
		Return(mockKamiUser, nil).Once()
	s.repo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(nil).Once()

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	res, err := u.LoginKamiUser(context.TODO(), &domain.KamiUser{Phone: "0900123123", Password: "123"})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserUcaseSuite) TestCheckKamiUserLog_Success_NewLog() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	s.repo.On("GetLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(&domain.KamiUserLog{}, nil).Once()
	s.repo.On("NewLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log.Phone == mockKamiUserLog.Phone })).
		Return(mockKamiUserLog, nil)
	s.repo.On("UpdateLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(nil)

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.CheckKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserUcaseSuite) TestCheckKamiUserLog_Success_CheckLog() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now().Add(-time.Second * 60),
	}

	s.repo.On("GetLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(mockKamiUserLog, nil).Once()
	s.repo.On("UpdateLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(nil)

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.CheckKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserUcaseSuite) TestCheckKamiUserLog_Fail_TheSendingTimeOfTheTwoVerificationCodesMustExceed60Seconds() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	s.repo.On("GetLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(mockKamiUserLog, nil).Once()
	s.repo.On("UpdateLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(nil)

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.CheckKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.EqualError(s.Suite.T(), err, "the sending time of the two verification codes must exceed 60 seconds")
}

func (s *KamiUserUcaseSuite) TestCheckKamiUserLog_Fail_TheLoginActivityIsAbnormal() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:         gorm.Model{ID: 1},
		Phone:         "0900123123",
		AuthTime:      time.Now(),
		AuthFrequency: 3,
	}

	s.repo.On("GetLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(mockKamiUserLog, nil).Once()
	s.repo.On("UpdateLog", mock.Anything, mock.MatchedBy(func(log *domain.KamiUserLog) bool { return log == mockKamiUserLog })).
		Return(nil)

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.CheckKamiUserLog(context.TODO(), mockKamiUserLog)
	assert.EqualError(s.Suite.T(), err, "the login activity is abnormal, please try again in 24 hours or contact KAMIKAMI customer service")
}

func (s *KamiUserUcaseSuite) TestGenerateToken_Success() {
	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	_, err := u.GenerateToken(context.TODO(), &domain.KamiUser{Phone: "0900123123"})
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserUcaseSuite) TestUpdateUserInfo_Success() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
	}

	s.repo.On("Update", mock.Anything, mock.MatchedBy(func(user *domain.KamiUser) bool { return user == mockKamiUser })).
		Return(nil)

	u := ucase.NewKamiUserUsecase([]byte(mock.Anything), s.repo)
	err := u.UpdateUserInfo(context.TODO(), mockKamiUser, &domain.KamiUserInput{Gender: "M", Name: "Test", Career: "Test", Birthday: "2000-01-01"})
	assert.NoError(s.Suite.T(), err)
}
