package postgresql

import (
	"context"
	"kami/domain"

	"gorm.io/gorm"
)

type postgresqlPrizePoolRepository struct {
	db *gorm.DB
}

func NewPostgresqlPrizePoolRepository(db *gorm.DB) domain.PrizePoolRepository {
	return &postgresqlPrizePoolRepository{db}
}

func (pp *postgresqlPrizePoolRepository) New(ctx context.Context, prizePool *domain.PrizePool) (*domain.PrizePool, error) {
	result := pp.db.Create(&prizePool)
	return prizePool, result.Error
}

func (pp *postgresqlPrizePoolRepository) Get(ctx context.Context, prizePool *domain.PrizePool, optsWhere ...map[string]interface{}) (result *domain.PrizePool, err error) {
	query := pp.db.Model(&domain.PrizePool{}).Where(prizePool)
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

func (pp *postgresqlPrizePoolRepository) Gets(ctx context.Context, prizePool *domain.PrizePool, optsWhere ...map[string]interface{}) (result []*domain.PrizePool, err error) {
	query := pp.db.Model(&domain.PrizePool{}).Where(prizePool)
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
