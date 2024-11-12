package usecase_test

import (
	"context"
	"kami/domain"
	"kami/domain/mocks"
	"testing"
	"time"

	ucase "kami/lottery/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type LotteryPrizePoolUcaseSuite struct {
	suite.Suite
	repo *mocks.PrizePoolRepository
}

type LotteryPrizeCardUcaseSuite struct {
	suite.Suite
	repo *mocks.PrizeCardRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &LotteryPrizePoolUcaseSuite{})
	suite.Run(t, &LotteryPrizeCardUcaseSuite{})
}

func (s *LotteryPrizePoolUcaseSuite) SetupTest() {
	s.repo = new(mocks.PrizePoolRepository)
}

func (s *LotteryPrizeCardUcaseSuite) SetupTest() {
	s.repo = new(mocks.PrizeCardRepository)
}

func (s *LotteryPrizePoolUcaseSuite) TestGetPrizePool_Success() {
	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(prizePool *domain.PrizePool) bool { return prizePool == mockPrizePool })).
		Return(mockPrizePool, nil).Once()

	u := ucase.NewLotteryPrizePoolUsecase(mock.Anything, s.repo)
	res, err := u.GetPrizePool(context.TODO(), mockPrizePool)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizePool, res)
}

func (s *LotteryPrizePoolUcaseSuite) TestGetPrizePoolList_Success() {
	mockPrizePools := []*domain.PrizePool{
		{
			Model:  gorm.Model{ID: 1},
			Owner:  "ABC",
			Name:   "testingPool1",
			Points: 1,
		},
		{
			Model:  gorm.Model{ID: 2},
			Owner:  "ABC",
			Name:   "testingPool2",
			Points: 1,
		},
	}

	s.repo.On("Gets", mock.Anything, mock.MatchedBy(func(prizePool *domain.PrizePool) bool { return prizePool.Owner == mockPrizePools[0].Owner })).
		Return(mockPrizePools, nil).Once()

	u := ucase.NewLotteryPrizePoolUsecase(mock.Anything, s.repo)
	res, err := u.GetPrizePoolList(context.TODO(), &domain.PrizePool{Owner: "ABC"})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizePools, res)
}

func (s *LotteryPrizePoolUcaseSuite) TestSubtractUserPoints_Success() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
		Points:   10,
	}

	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	u := ucase.NewLotteryPrizePoolUsecase(mock.Anything, s.repo)
	err := u.SubtractUserPoints(context.TODO(), mockKamiUser, mockPrizePool)
	assert.NoError(s.Suite.T(), err)
}

func (s *LotteryPrizePoolUcaseSuite) TestSubtractUserPoints_Fail() {
	mockKamiUser := &domain.KamiUser{
		Model:    gorm.Model{ID: 1},
		Email:    "test@test.com",
		Phone:    "0900123123",
		Password: "123",
		Points:   0,
	}

	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	u := ucase.NewLotteryPrizePoolUsecase(mock.Anything, s.repo)
	err := u.SubtractUserPoints(context.TODO(), mockKamiUser, mockPrizePool)
	assert.EqualError(s.Suite.T(), err, "not enough points")
}

func (s *LotteryPrizeCardUcaseSuite) TestGetPrizeCard_Success() {
	mockPrizeCard := &domain.PrizeCard{
		Model:       gorm.Model{ID: 1},
		PoolId:      1,
		Title:       "testingCard",
		Description: "testingCardDescription",
		Probability: 100,
	}

	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(prizeCard *domain.PrizeCard) bool { return prizeCard == mockPrizeCard })).
		Return(mockPrizeCard, nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res, err := u.GetPrizeCard(context.TODO(), mockPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCard, res)
}

func (s *LotteryPrizeCardUcaseSuite) TestGetPrizeCardList_Success() {
	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard1",
			Description: "testingCard1Description",
			Probability: 50,
		},
		{
			Model:       gorm.Model{ID: 2},
			PoolId:      1,
			Title:       "testingCard2",
			Description: "testingCard2Description",
			Probability: 50,
		},
	}

	s.repo.On("Gets", mock.Anything, mock.MatchedBy(func(prizeCard *domain.PrizeCard) bool { return prizeCard.PoolId == mockPrizeCards[0].PoolId })).
		Return(mockPrizeCards, nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res, err := u.GetPrizeCardList(context.TODO(), &domain.PrizeCard{PoolId: 1})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCards, res)
}

func (s *LotteryPrizeCardUcaseSuite) TestGetWeightedRandomList_Success() {
	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard1",
			Description: "testingCard1Description",
			Probability: 50,
		},
		{
			Model:       gorm.Model{ID: 2},
			PoolId:      1,
			Title:       "testingCard2",
			Description: "testingCard2Description",
			Probability: 50,
		},
		{
			Model:       gorm.Model{ID: 3},
			PoolId:      1,
			Title:       "testingCard3",
			Description: "testingCard3Description",
			Probability: 50,
		},
		{
			Model:       gorm.Model{ID: 4},
			PoolId:      1,
			Title:       "testingCard4",
			Description: "testingCard4Description",
			Probability: 50,
		},
	}

	mockUserPrizeCards := []*domain.UserPrizeCard{
		{
			ID:           1,
			UserId:       1,
			CardId:       1,
			SerialNumber: "a1b2c3",
		},
		{
			ID:           2,
			UserId:       1,
			CardId:       2,
			SerialNumber: "d1e2f3",
		},
	}

	s.repo.On("GetUserPrizeCardList", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool {
		return userPrizeCard.UserId == mockUserPrizeCards[0].UserId
	}), mock.Anything).
		Return(mockUserPrizeCards, nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res, err := u.GetWeightedRandomList(context.TODO(), 1, mockPrizeCards)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), []*domain.PrizeCard{mockPrizeCards[2], mockPrizeCards[3]}, res)
}

func (s *LotteryPrizeCardUcaseSuite) TestDraw_Success() {
	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard1",
			Description: "testingCard1Description",
			Probability: 100,
		},
	}

	mockUserPrizeCards := []*domain.UserPrizeCard{
		{
			ID:           1,
			UserId:       1,
			CardId:       1,
			SerialNumber: "a1b2c3",
		},
	}

	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           2,
		UserId:       1,
		CardId:       1,
		SerialNumber: "112",
	}

	mockCoupon := &domain.Coupon{
		ID:         1,
		CardId:     1,
		UserId:     1,
		Code:       "123",
		ExpiryDate: time.Now().AddDate(0, 0, 1),
	}

	mockPrizeCardOutput := &domain.PrizeCardOutput{
		ID:           mockPrizeCards[0].Model.ID,
		Title:        mockPrizeCards[0].Title,
		Description:  mockPrizeCards[0].Description,
		SerialNumber: mockUserPrizeCard.SerialNumber,
	}
	s.repo.On("GetUserPrizeCardList", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool {
		return userPrizeCard.CardId == mockUserPrizeCards[0].CardId
	})).
		Return(mockUserPrizeCards, nil).Once()
	s.repo.On("NewUserPrizeCard", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool {
		return userPrizeCard.UserId == mockUserPrizeCard.UserId && userPrizeCard.CardId == mockUserPrizeCard.CardId && userPrizeCard.SerialNumber == mockUserPrizeCard.SerialNumber
	})).Return(mockUserPrizeCard, nil).Once()
	s.repo.On("GetCoupon", mock.Anything, mock.MatchedBy(func(coupon *domain.Coupon) bool { return coupon.CardId == mockCoupon.CardId }), mock.Anything).
		Return(mockCoupon, nil).Once()
	s.repo.On("UpdateCoupon", mock.Anything, mock.MatchedBy(func(coupon *domain.Coupon) bool { return coupon == mockCoupon })).
		Return(nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res, err := u.Draw(context.TODO(), 1, mockPrizeCards)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCardOutput, res)
}

func (s *LotteryPrizeCardUcaseSuite) TestGetPrizeCardCollection_Success() {
	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	mockPrizeCards := []*domain.PrizeCard{
		{
			Model:       gorm.Model{ID: 1},
			PoolId:      1,
			Title:       "testingCard",
			Description: "testingCardDescription",
			Probability: 100,
		},
	}

	mockUserPrizeCards := []*domain.UserPrizeCard{
		{
			ID:           1,
			UserId:       1,
			CardId:       1,
			SerialNumber: "a1b2c3",
		},
	}

	mockPrizeCardOutputs := []domain.PrizeCardOutput{
		{
			ID:           mockPrizeCards[0].Model.ID,
			Title:        mockPrizeCards[0].Title,
			Description:  mockPrizeCards[0].Description,
			SerialNumber: mockUserPrizeCards[0].SerialNumber,
		},
	}

	mockPrizeCardCollection := &domain.PrizeCardCollection{
		PoolName: mockPrizePool.Name,
		Cards:    mockPrizeCardOutputs,
		Total:    int64(len(mockPrizeCards)),
	}

	s.repo.On("Gets", mock.Anything, mock.MatchedBy(func(prizeCard *domain.PrizeCard) bool { return prizeCard.PoolId == mockPrizePool.Model.ID })).
		Return(mockPrizeCards, nil).Once()
	s.repo.On("GetUserPrizeCardList", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool {
		return userPrizeCard.UserId == mockUserPrizeCards[0].UserId
	}), mock.Anything).
		Return(mockUserPrizeCards, nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res, err := u.GetPrizeCardCollection(context.TODO(), 1, mockPrizePool)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCardCollection, res)
}

func (s *LotteryPrizeCardUcaseSuite) TestGetPrizeCardCollectionDetail_Success() {
	mockPrizeCard := &domain.PrizeCard{
		Model:       gorm.Model{ID: 1},
		PoolId:      1,
		Title:       "testingCard",
		Description: "testingCardDescription",
		Probability: 100,
	}

	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           2,
		UserId:       1,
		CardId:       1,
		SerialNumber: "112",
	}

	mockCoupon := &domain.Coupon{
		ID:         1,
		CardId:     1,
		UserId:     1,
		Code:       "123",
		ExpiryDate: time.Now().AddDate(0, 0, 1),
	}

	s.repo.On("GetUserPrizeCard", mock.Anything, mock.MatchedBy(func(userPrizeCard *domain.UserPrizeCard) bool { return userPrizeCard == mockUserPrizeCard })).
		Return(mockUserPrizeCard, nil).Once()
	s.repo.On("GetCoupon", mock.Anything, mock.MatchedBy(func(coupon *domain.Coupon) bool {
		return coupon.CardId == mockCoupon.CardId && coupon.UserId == mockCoupon.UserId
	})).
		Return(mockCoupon, nil).Once()
	s.repo.On("Get", mock.Anything, mock.MatchedBy(func(prizeCard *domain.PrizeCard) bool { return prizeCard.Model.ID == mockUserPrizeCard.CardId })).
		Return(mockPrizeCard, nil).Once()

	u := ucase.NewLotteryPrizeCardUsecase(mock.Anything, s.repo)
	res1, res2, res3, err := u.GetPrizeCardCollectionDetail(context.TODO(), mockUserPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockUserPrizeCard.SerialNumber, res1)
	assert.Equal(s.Suite.T(), mockCoupon, res2)
	assert.Equal(s.Suite.T(), mockPrizeCard, res3)
}
