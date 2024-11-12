package postgresql

import (
	"context"
	"kami/domain"

	"gorm.io/gorm"
)

type postgresqlPrizeCardRepository struct {
	db *gorm.DB
}

func NewPostgresqlPrizeCardRepository(db *gorm.DB) domain.PrizeCardRepository {
	return &postgresqlPrizeCardRepository{db}
}

func (pc *postgresqlPrizeCardRepository) New(ctx context.Context, prizeCard *domain.PrizeCard) (*domain.PrizeCard, error) {
	result := pc.db.Create(&prizeCard)
	return prizeCard, result.Error
}

func (pc *postgresqlPrizeCardRepository) Get(ctx context.Context, prizeCard *domain.PrizeCard, optsWhere ...map[string]interface{}) (result *domain.PrizeCard, err error) {
	query := pc.db.Model(&domain.PrizeCard{}).Where(prizeCard)
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

func (pc *postgresqlPrizeCardRepository) Gets(ctx context.Context, prizeCard *domain.PrizeCard, optsWhere ...map[string]interface{}) (result []*domain.PrizeCard, err error) {
	query := pc.db.Model(&domain.PrizeCard{}).Where(prizeCard)
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

func (pc *postgresqlPrizeCardRepository) NewUserPrizeCard(ctx context.Context, userPrizeCard *domain.UserPrizeCard) (*domain.UserPrizeCard, error) {
	result := pc.db.Create(&userPrizeCard)
	return userPrizeCard, result.Error
}

func (pc *postgresqlPrizeCardRepository) GetUserPrizeCard(ctx context.Context, userPrizeCard *domain.UserPrizeCard, optsWhere ...map[string]interface{}) (result *domain.UserPrizeCard, err error) {
	query := pc.db.Model(&domain.UserPrizeCard{}).Where(userPrizeCard)
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

func (pc *postgresqlPrizeCardRepository) GetUserPrizeCardList(ctx context.Context, userPrizeCard *domain.UserPrizeCard, optsWhere ...map[string]interface{}) (result []*domain.UserPrizeCard, err error) {
	query := pc.db.Model(&domain.UserPrizeCard{}).Where(userPrizeCard)
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

func (pc *postgresqlPrizeCardRepository) GetCoupon(ctx context.Context, coupon *domain.Coupon, optsWhere ...map[string]interface{}) (result *domain.Coupon, err error) {
	query := pc.db.Model(&domain.Coupon{}).Where(coupon)
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

func (pc *postgresqlPrizeCardRepository) UpdateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	return pc.db.Save(&coupon).Error
}
