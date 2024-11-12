package repository

import (
	"kami/domain"
)

type twilioServiceRepository struct {
}

func NewTwilioServiceRepository() domain.TwilioServiceRepository {
	return &twilioServiceRepository{}
}

func (p *twilioServiceRepository) GetTwilioConfig() (*domain.TwilioService, error) {
	return &domain.TwilioService{}, nil
}
