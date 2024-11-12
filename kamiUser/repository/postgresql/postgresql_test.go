package postgresql_test

import (
	"context"
	"kami/domain"
	"testing"
	"time"

	kamiUserPostgresqlRepo "kami/kamiUser/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type KamiUserPostgresqlRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.KamiUserRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &KamiUserPostgresqlRepoSuite{})
}

func (s *KamiUserPostgresqlRepoSuite) SetupTest() {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	s.mock = mock
	s.repo = kamiUserPostgresqlRepo.NewPostgresqlKamiUserRepository(gormDB)
}

func (s *KamiUserPostgresqlRepoSuite) TestNew() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "kami_users" (.+)`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	res, err := s.repo.New(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserPostgresqlRepoSuite) TestGet() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "phone"}).
		AddRow(mockKamiUser.Model.ID, mockKamiUser.Email, mockKamiUser.Phone)

	s.mock.ExpectQuery(`SELECT (.*) FROM "kami_users" WHERE "kami_users"."id" = \$1 AND "kami_users"."email" = \$2 AND "kami_users"."phone" = \$3`).
		WithArgs(mockKamiUser.Model.ID, mockKamiUser.Email, mockKamiUser.Phone).WillReturnRows(rows)

	res, err := s.repo.Get(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUser, res)
}

func (s *KamiUserPostgresqlRepoSuite) TestUpdate() {
	mockKamiUser := &domain.KamiUser{
		Model: gorm.Model{ID: 1},
		Email: "test@test.com",
		Phone: "0900123123",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "kami_users"`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.Update(context.TODO(), mockKamiUser)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiUserPostgresqlRepoSuite) TestNewLog() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "phone", "auth_time"}).
		AddRow(mockKamiUserLog.Model.ID, mockKamiUserLog.Phone, mockKamiUserLog.AuthTime)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "kami_user_logs" (.+)`).WillReturnRows(rows)
	s.mock.ExpectCommit()

	res, err := s.repo.NewLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUserLog, res)
}

func (s *KamiUserPostgresqlRepoSuite) TestGetLog() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "phone", "auth_time"}).
		AddRow(mockKamiUserLog.Model.ID, mockKamiUserLog.Phone, mockKamiUserLog.AuthTime)

	s.mock.ExpectQuery(`SELECT (.*) FROM "kami_user_logs" WHERE "kami_user_logs"."id" = \$1 AND "kami_user_logs"."phone" = \$2 AND "kami_user_logs"."auth_time" = \$3`).
		WithArgs(mockKamiUserLog.Model.ID, mockKamiUserLog.Phone, mockKamiUserLog.AuthTime).
		WillReturnRows(rows)

	res, err := s.repo.GetLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiUserLog, res)
}

func (s *KamiUserPostgresqlRepoSuite) TestUpdateLog() {
	mockKamiUserLog := &domain.KamiUserLog{
		Model:    gorm.Model{ID: 1},
		Phone:    "0900123123",
		AuthTime: time.Now(),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "kami_user_logs"`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.UpdateLog(context.TODO(), mockKamiUserLog)
	assert.NoError(s.Suite.T(), err)
}
