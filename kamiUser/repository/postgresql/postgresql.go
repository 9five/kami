package postgresql

import (
	"context"
	_ "time/tzdata"

	"gorm.io/gorm"

	"kami/domain"
	// "gorm.io/gorm/clause"
)

type postgresqlKamiUserRepository struct {
	db *gorm.DB
}

func NewPostgresqlKamiUserRepository(db *gorm.DB) domain.KamiUserRepository {
	return &postgresqlKamiUserRepository{db}
}

func (p *postgresqlKamiUserRepository) New(ctx context.Context, user *domain.KamiUser) (*domain.KamiUser, error) {
	result := p.db.Create(&user)
	return user, result.Error
}

func (p *postgresqlKamiUserRepository) Get(ctx context.Context, user *domain.KamiUser, optsWhere ...map[string]interface{}) (result *domain.KamiUser, err error) {
	query := p.db.Model(&domain.KamiUser{}).Where(user)
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

func (p *postgresqlKamiUserRepository) Update(ctx context.Context, user *domain.KamiUser) error {
	return p.db.Save(&user).Error
}

func (p *postgresqlKamiUserRepository) NewLog(ctx context.Context, log *domain.KamiUserLog) (*domain.KamiUserLog, error) {
	result := p.db.Create(&log)
	return log, result.Error
}

func (p *postgresqlKamiUserRepository) GetLog(ctx context.Context, log *domain.KamiUserLog) (result *domain.KamiUserLog, err error) {
	err = p.db.Model(&domain.KamiUserLog{}).Where(log).First(&result).Error
	return
}

func (p *postgresqlKamiUserRepository) UpdateLog(ctx context.Context, log *domain.KamiUserLog) error {
	return p.db.Save(&log).Error
}
