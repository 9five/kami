package postgresql_test

import (
	"context"
	"kami/domain"
	"testing"

	kamiOrderPostgresqlRepo "kami/kamiOrder/repository/postgresql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type KamiOrderPostgresqlRepoSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	repo domain.KamiOrderRepository
}

func TestStart(t *testing.T) {
	suite.Run(t, &KamiOrderPostgresqlRepoSuite{})
}

func (s *KamiOrderPostgresqlRepoSuite) SetupTest() {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}))
	if err != nil {
		s.FailNowf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	s.mock = mock
	s.repo = kamiOrderPostgresqlRepo.NewPostgresqlKamiOrderRepository(gormDB)
}

func (s *KamiOrderPostgresqlRepoSuite) TestStore() {
	mockKamiOrder := &domain.KamiOrder{
		Model:         gorm.Model{ID: 1},
		OrderId:       "123-123",
		Restaurant:    "ABC",
		Status:        "Delivered",
		BillingStatus: "Payable",
	}

	rows := sqlmock.NewRows([]string{"id", "order_id", "restaurant", "status", "billing_status"}).
		AddRow(mockKamiOrder.Model.ID, mockKamiOrder.OrderId, mockKamiOrder.Restaurant, mockKamiOrder.Status, mockKamiOrder.BillingStatus)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "kami_orders" (.+)`).WillReturnRows(rows)
	s.mock.ExpectCommit()

	err := s.repo.Store(context.TODO(), mockKamiOrder)
	assert.NoError(s.Suite.T(), err)
}

func (s *KamiOrderPostgresqlRepoSuite) TestGet() {
	mockKamiOrder := &domain.KamiOrder{
		Model:         gorm.Model{ID: 1},
		OrderId:       "123-123",
		Restaurant:    "ABC",
		Status:        "Delivered",
		BillingStatus: "Payable",
	}

	rows := sqlmock.NewRows([]string{"id", "order_id", "restaurant", "status", "billing_status"}).
		AddRow(mockKamiOrder.Model.ID, mockKamiOrder.OrderId, mockKamiOrder.Restaurant, mockKamiOrder.Status, mockKamiOrder.BillingStatus)

	s.mock.ExpectQuery(`SELECT \* FROM "kami_orders" WHERE "kami_orders"."id" = \$1 AND "kami_orders"."order_id" = \$2 AND "kami_orders"."restaurant" = \$3 AND "kami_orders"."status" = \$4 AND "kami_orders"."billing_status" = \$5 AND "kami_orders"."deleted_at" IS NULL ORDER BY "kami_orders"."id" LIMIT 1`).
		WithArgs(mockKamiOrder.Model.ID, mockKamiOrder.OrderId, mockKamiOrder.Restaurant, mockKamiOrder.Status, mockKamiOrder.BillingStatus).
		WillReturnRows(rows)

	res, err := s.repo.Get(context.TODO(), mockKamiOrder)
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiOrder, res)
}

func (s *KamiOrderPostgresqlRepoSuite) TestGets() {
	mockKamiOrders := []*domain.KamiOrder{
		{
			Model:         gorm.Model{ID: 1},
			OrderId:       "123-123",
			Restaurant:    "ABC",
			Status:        "Delivered",
			BillingStatus: "Payable",
		},
		{
			Model:         gorm.Model{ID: 2},
			OrderId:       "456-456",
			Restaurant:    "ABC",
			Status:        "Delivered",
			BillingStatus: "Payable",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "order_id", "restaurant", "status", "billing_status"}).
		AddRow(mockKamiOrders[0].Model.ID, mockKamiOrders[0].OrderId, mockKamiOrders[0].Restaurant, mockKamiOrders[0].Status, mockKamiOrders[0].BillingStatus).
		AddRow(mockKamiOrders[1].Model.ID, mockKamiOrders[1].OrderId, mockKamiOrders[1].Restaurant, mockKamiOrders[1].Status, mockKamiOrders[1].BillingStatus)

	s.mock.ExpectQuery(`SELECT \* FROM "kami_orders" WHERE "kami_orders"."restaurant" = \$1`).
		WithArgs("ABC").WillReturnRows(rows)

	res, err := s.repo.Gets(context.TODO(), &domain.KamiOrder{Restaurant: "ABC"})
	assert.NoError(s.Suite.T(), err)
	assert.Equal(s.Suite.T(), mockKamiOrders, res)
}

func (s *KamiOrderPostgresqlRepoSuite) TestUpdate() {
	mockKamiOrder := &domain.KamiOrder{
		Model:         gorm.Model{ID: 1},
		OrderId:       "123-123",
		Restaurant:    "ABC",
		Status:        "Delivered",
		BillingStatus: "Payable",
	}

	// s.mock.ExpectExec(`UPDATE "kami_orders" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"order_id"=\$4,"restaurant"=\$5,"status"=\$6,"billing_status"=\$7,"order_placed_at"=\$8,"order_delivered_at"=\$9,"platform"=\$10,"owner_phone"=\$11 WHERE "kami_orders"."deleted_at" IS NULL AND "id" = \$12`).
	// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), mockKamiOrder.OrderId, mockKamiOrder.Restaurant, mockKamiOrder.Status, mockKamiOrder.BillingStatus, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), mockKamiOrder.Model.ID).
	// 	WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "kami_orders"`).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repo.Update(context.TODO(), mockKamiOrder)
	assert.NoError(s.Suite.T(), err)
}
