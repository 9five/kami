package postgresql

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kami/domain"
)

type postgresqlKamiOrderRepository struct {
	db *gorm.DB
}

func NewPostgresqlKamiOrderRepository(db *gorm.DB) domain.KamiOrderRepository {
	return &postgresqlKamiOrderRepository{db}
}

func (p *postgresqlKamiOrderRepository) Store(ctx context.Context, order *domain.KamiOrder) error {
	return p.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "billing_status", "order_delivered_at"}),
	}).Create(&order).Error
}

func (p *postgresqlKamiOrderRepository) Get(ctx context.Context, order *domain.KamiOrder, optsWhere ...map[string]interface{}) (result *domain.KamiOrder, err error) {
	query := p.db.Model(&domain.KamiOrder{}).Where(order)
	if len(optsWhere) != 0 {
		for _, optWhere := range optsWhere {
			for k, v := range optWhere {
				if v == nil {
					query.Where(k)
				} else {
					query.Where(k, v)
				}
			}
		}
	}
	err = query.First(&result).Error
	return
}

func (p *postgresqlKamiOrderRepository) Gets(ctx context.Context, order *domain.KamiOrder, optsWhere ...map[string]interface{}) (result []*domain.KamiOrder, err error) {
	query := p.db.Model(&domain.KamiOrder{}).Where(order)
	if len(optsWhere) != 0 {
		for _, optWhere := range optsWhere {
			for k, v := range optWhere {
				if v == nil {
					query.Where(k)
				} else {
					query.Where(k, v)
				}
			}
		}
	}
	err = query.Find(&result).Error
	return
}

func (p *postgresqlKamiOrderRepository) Update(ctx context.Context, order *domain.KamiOrder) error {
	return p.db.Save(&order).Error
}
