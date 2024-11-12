package postgresql_test

import (
	"context"
	"kami/domain"
	"testing"
	"time"

	lotteryPostgresqlRepo "kami/lottery/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LotteryPrizePoolPostgresqlRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.PrizePoolRepository
}

type LotteryPrizeCardPostgresqlRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.PrizeCardRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &LotteryPrizePoolPostgresqlRepoSuite{})
	suite.Run(t, &LotteryPrizeCardPostgresqlRepoSuite{})
}

func (s *LotteryPrizePoolPostgresqlRepoSuite) SetupTest() {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	s.mock = mock
	s.repo = lotteryPostgresqlRepo.NewPostgresqlPrizePoolRepository(gormDB)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) SetupTest() {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	s.mock = mock
	s.repo = lotteryPostgresqlRepo.NewPostgresqlPrizeCardRepository(gormDB)
}

func (s *LotteryPrizePoolPostgresqlRepoSuite) TestNew() {
	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "owner", "name", "points"}).
		AddRow(mockPrizePool.Model.ID, mockPrizePool.Owner, mockPrizePool.Name, mockPrizePool.Points)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "prize_pools" (.+)`).WillReturnRows(rows)
	s.mock.ExpectCommit()

	res, err := s.repo.New(context.TODO(), mockPrizePool)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizePool, res)
}

func (s *LotteryPrizePoolPostgresqlRepoSuite) TestGet() {
	mockPrizePool := &domain.PrizePool{
		Model:  gorm.Model{ID: 1},
		Owner:  "ABC",
		Name:   "testingPool",
		Points: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "owner", "name", "points"}).
		AddRow(mockPrizePool.Model.ID, mockPrizePool.Owner, mockPrizePool.Name, mockPrizePool.Points)

	s.mock.ExpectQuery(`SELECT \* FROM "prize_pools" WHERE "prize_pools"."id" = \$1 AND "prize_pools"."owner" = \$2 AND "prize_pools"."name" = \$3 AND "prize_pools"."points" = \$4`).
		WithArgs(mockPrizePool.Model.ID, mockPrizePool.Owner, mockPrizePool.Name, mockPrizePool.Points).
		WillReturnRows(rows)

	res, err := s.repo.Get(context.TODO(), mockPrizePool)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizePool, res)
}

func (s *LotteryPrizePoolPostgresqlRepoSuite) TestGets() {
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

	rows := sqlmock.NewRows([]string{"id", "owner", "name", "points"}).
		AddRow(mockPrizePools[0].Model.ID, mockPrizePools[0].Owner, mockPrizePools[0].Name, mockPrizePools[0].Points).
		AddRow(mockPrizePools[1].Model.ID, mockPrizePools[1].Owner, mockPrizePools[1].Name, mockPrizePools[1].Points)

	s.mock.ExpectQuery(`SELECT \* FROM "prize_pools" WHERE "prize_pools"."owner" = \$1`).
		WithArgs(mockPrizePools[0].Owner).
		WillReturnRows(rows)

	res, err := s.repo.Gets(context.TODO(), &domain.PrizePool{Owner: "ABC"})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizePools, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestNew() {
	mockPrizeCard := &domain.PrizeCard{
		Model:       gorm.Model{ID: 1},
		PoolId:      1,
		Title:       "testingCard",
		Description: "testingCardDescription",
		Probability: 100,
	}

	rows := sqlmock.NewRows([]string{"id", "pool_id", "title", "description", "probability"}).
		AddRow(mockPrizeCard.Model.ID, mockPrizeCard.PoolId, mockPrizeCard.Title, mockPrizeCard.Description, mockPrizeCard.Probability)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "prize_cards" (.+)`).WillReturnRows(rows)
	s.mock.ExpectCommit()

	res, err := s.repo.New(context.TODO(), mockPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCard, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestGet() {
	mockPrizeCard := &domain.PrizeCard{
		Model:       gorm.Model{ID: 1},
		PoolId:      1,
		Title:       "testingCard",
		Description: "testingCardDescription",
		Probability: 100,
	}

	rows := sqlmock.NewRows([]string{"id", "pool_id", "title", "description", "probability"}).
		AddRow(mockPrizeCard.Model.ID, mockPrizeCard.PoolId, mockPrizeCard.Title, mockPrizeCard.Description, mockPrizeCard.Probability)

	s.mock.ExpectQuery(`SELECT \* FROM "prize_cards" WHERE "prize_cards"."id" = \$1 AND "prize_cards"."pool_id" = \$2 AND "prize_cards"."title" = \$3 AND "prize_cards"."description" = \$4 AND "prize_cards"."probability" = \$5`).
		WithArgs(mockPrizeCard.Model.ID, mockPrizeCard.PoolId, mockPrizeCard.Title, mockPrizeCard.Description, mockPrizeCard.Probability).
		WillReturnRows(rows)

	res, err := s.repo.Get(context.TODO(), mockPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCard, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestGets() {
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

	rows := sqlmock.NewRows([]string{"id", "pool_id", "title", "description", "probability"}).
		AddRow(mockPrizeCards[0].Model.ID, mockPrizeCards[0].PoolId, mockPrizeCards[0].Title, mockPrizeCards[0].Description, mockPrizeCards[0].Probability).
		AddRow(mockPrizeCards[1].Model.ID, mockPrizeCards[1].PoolId, mockPrizeCards[1].Title, mockPrizeCards[1].Description, mockPrizeCards[1].Probability)

	s.mock.ExpectQuery(`SELECT \* FROM "prize_cards" WHERE "prize_cards"."pool_id" = \$1`).
		WithArgs(mockPrizeCards[0].PoolId).
		WillReturnRows(rows)

	res, err := s.repo.Gets(context.TODO(), &domain.PrizeCard{PoolId: 1})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockPrizeCards, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestNewUserPrizeCard() {
	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           1,
		UserId:       1,
		CardId:       1,
		SerialNumber: "a1b2c3",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "card_id", "serial_number"}).
		AddRow(mockUserPrizeCard.ID, mockUserPrizeCard.UserId, mockUserPrizeCard.CardId, mockUserPrizeCard.SerialNumber)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "user_prize_cards" (.+)`).WillReturnRows(rows)
	s.mock.ExpectCommit()

	res, err := s.repo.NewUserPrizeCard(context.TODO(), mockUserPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockUserPrizeCard, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestGetUserPrizeCard() {
	mockUserPrizeCard := &domain.UserPrizeCard{
		ID:           1,
		UserId:       1,
		CardId:       1,
		SerialNumber: "a1b2c3",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "card_id", "serial_number"}).
		AddRow(mockUserPrizeCard.ID, mockUserPrizeCard.UserId, mockUserPrizeCard.CardId, mockUserPrizeCard.SerialNumber)

	s.mock.ExpectQuery(`SELECT \* FROM "user_prize_cards" WHERE "user_prize_cards"."id" = \$1 AND "user_prize_cards"."user_id" = \$2 AND "user_prize_cards"."card_id" = \$3 AND "user_prize_cards"."serial_number" = \$4`).
		WithArgs(mockUserPrizeCard.ID, mockUserPrizeCard.UserId, mockUserPrizeCard.CardId, mockUserPrizeCard.SerialNumber).
		WillReturnRows(rows)

	res, err := s.repo.GetUserPrizeCard(context.TODO(), mockUserPrizeCard)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockUserPrizeCard, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestGetUserPrizeCardList() {
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

	rows := sqlmock.NewRows([]string{"id", "user_id", "card_id", "serial_number"}).
		AddRow(mockUserPrizeCards[0].ID, mockUserPrizeCards[0].UserId, mockUserPrizeCards[0].CardId, mockUserPrizeCards[0].SerialNumber).
		AddRow(mockUserPrizeCards[1].ID, mockUserPrizeCards[1].UserId, mockUserPrizeCards[1].CardId, mockUserPrizeCards[1].SerialNumber)

	s.mock.ExpectQuery(`SELECT \* FROM "user_prize_cards" WHERE "user_prize_cards"."user_id" = \$1`).
		WithArgs(mockUserPrizeCards[0].UserId).
		WillReturnRows(rows)

	res, err := s.repo.GetUserPrizeCardList(context.TODO(), &domain.UserPrizeCard{UserId: mockUserPrizeCards[0].UserId})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockUserPrizeCards, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestGetCoupon() {
	mockCoupon := &domain.Coupon{
		ID:         1,
		CardId:     1,
		UserId:     1,
		Code:       "123",
		ExpiryDate: time.Now().AddDate(0, 0, 1),
	}

	rows := sqlmock.NewRows([]string{"id", "card_id", "user_id", "code", "expiry_date"}).
		AddRow(mockCoupon.ID, mockCoupon.CardId, mockCoupon.UserId, mockCoupon.Code, mockCoupon.ExpiryDate)

	s.mock.ExpectQuery(`SELECT \* FROM "coupons" WHERE "coupons"."id" = \$1 AND "coupons"."card_id" = \$2 AND "coupons"."user_id" = \$3 AND "coupons"."code" = \$4 AND "coupons"."expiry_date" = \$5`).
		WithArgs(mockCoupon.ID, mockCoupon.CardId, mockCoupon.UserId, mockCoupon.Code, mockCoupon.ExpiryDate).
		WillReturnRows(rows)

	res, err := s.repo.GetCoupon(context.TODO(), mockCoupon)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockCoupon, res)
}

func (s *LotteryPrizeCardPostgresqlRepoSuite) TestUpdateCoupon() {
	mockCoupon := &domain.Coupon{
		ID:         1,
		CardId:     1,
		UserId:     1,
		Code:       "123",
		ExpiryDate: time.Now().AddDate(0, 0, 1),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "coupons" `).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateCoupon(context.TODO(), mockCoupon)
	assert.NoError(s.Suite.T(), err)
}
