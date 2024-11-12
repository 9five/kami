package usecase

import (
	"context"
	"errors"
	"kami/domain"
)

type lotteryPrizePoolUsecase struct {
	awsBucket     string
	prizePoolRepo domain.PrizePoolRepository
}

func NewLotteryPrizePoolUsecase(awsBucket string, prizePoolRepo domain.PrizePoolRepository) domain.PrizePoolUsecase {
	return &lotteryPrizePoolUsecase{
		awsBucket:     awsBucket,
		prizePoolRepo: prizePoolRepo,
	}
}

func (pp *lotteryPrizePoolUsecase) GetPrizePool(ctx context.Context, prizePool *domain.PrizePool) (result *domain.PrizePool, err error) {
	return pp.prizePoolRepo.Get(ctx, prizePool)
}

func (pp *lotteryPrizePoolUsecase) GetPrizePoolList(ctx context.Context, prizePool *domain.PrizePool) ([]*domain.PrizePool, error) {
	return pp.prizePoolRepo.Gets(ctx, prizePool)
}

func (pp *lotteryPrizePoolUsecase) SubtractUserPoints(ctx context.Context, user *domain.KamiUser, prizePool *domain.PrizePool) error {
	if user.Points < prizePool.Points {
		return errors.New("not enough points")
	}

	user.Points = user.Points - prizePool.Points
	return nil
}
