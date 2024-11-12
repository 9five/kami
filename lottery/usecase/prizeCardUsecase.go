package usecase

import (
	"context"
	"errors"
	"fmt"
	"kami/domain"
	"math/rand"
	"time"

	"github.com/mroth/weightedrand"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type lotteryPrizeCardUsecase struct {
	awsBucket     string
	prizeCardRepo domain.PrizeCardRepository
}

func NewLotteryPrizeCardUsecase(awsBucket string, prizeCardRepo domain.PrizeCardRepository) domain.PrizeCardUsecase {
	return &lotteryPrizeCardUsecase{
		awsBucket:     awsBucket,
		prizeCardRepo: prizeCardRepo,
	}
}

func (pc *lotteryPrizeCardUsecase) GetPrizeCard(ctx context.Context, prizeCard *domain.PrizeCard) (*domain.PrizeCard, error) {
	return pc.prizeCardRepo.Get(ctx, prizeCard)
}

func (pc *lotteryPrizeCardUsecase) GetPrizeCardList(ctx context.Context, prizeCard *domain.PrizeCard) ([]*domain.PrizeCard, error) {
	return pc.prizeCardRepo.Gets(ctx, prizeCard)
}

func (pc *lotteryPrizeCardUsecase) GetWeightedRandomList(ctx context.Context, userId uint, prizeCardList []*domain.PrizeCard) ([]*domain.PrizeCard, error) {
	var prizeCardSliceString []uint
	for _, v := range prizeCardList {
		prizeCardSliceString = append(prizeCardSliceString, v.Model.ID)
	}

	userPrizeCardList, err := pc.prizeCardRepo.GetUserPrizeCardList(ctx, &domain.UserPrizeCard{UserId: userId}, map[string]interface{}{
		"card_id in (?)": prizeCardSliceString,
	})
	if err != nil {
		return []*domain.PrizeCard{}, err
	}

	var userPrizeCardListIdList []uint
	for _, v := range userPrizeCardList {
		userPrizeCardListIdList = append(userPrizeCardListIdList, v.CardId)
	}

	var result []*domain.PrizeCard
	for _, v := range prizeCardList {
		if !slices.Contains(userPrizeCardListIdList, v.Model.ID) {
			result = append(result, v)
		}
	}

	return result, nil
}

func (pc *lotteryPrizeCardUsecase) Draw(ctx context.Context, userId uint, prizeCardList []*domain.PrizeCard) (*domain.PrizeCardOutput, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	var sliceChoice []weightedrand.Choice
	if len(prizeCardList) != 0 {
		for _, v := range prizeCardList {
			sliceChoice = append(sliceChoice, weightedrand.Choice{
				Item:   v,
				Weight: uint(v.Probability),
			})
		}
	} else {
		return nil, errors.New("redemption completed")
	}

	chooser, err := weightedrand.NewChooser(sliceChoice...)
	if err != nil {
		return nil, err
	}

	prize := chooser.Pick().(*domain.PrizeCard)

	userPrizeCardListByCardId, _ := pc.prizeCardRepo.GetUserPrizeCardList(ctx, &domain.UserPrizeCard{CardId: prize.Model.ID})

	userPrizeCard, err := pc.prizeCardRepo.NewUserPrizeCard(ctx, &domain.UserPrizeCard{
		UserId:       userId,
		CardId:       prize.Model.ID,
		SerialNumber: fmt.Sprintf("%d%d%d", prize.Model.ID, userId, len(userPrizeCardListByCardId)+1),
	})
	if err != nil {
		return nil, err
	}

	if coupon, err := pc.prizeCardRepo.GetCoupon(ctx, &domain.Coupon{CardId: prize.Model.ID}, map[string]interface{}{
		"user_id is null": nil,
	}); err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != gorm.ErrRecordNotFound {
		coupon.UserId = userId
		if err = pc.prizeCardRepo.UpdateCoupon(ctx, coupon); err != nil {
			return nil, err
		}
	}

	return &domain.PrizeCardOutput{
		ID:           prize.Model.ID,
		Picture:      prize.Picture,
		Title:        prize.Title,
		Description:  prize.Description,
		Audio:        prize.Audio,
		SerialNumber: userPrizeCard.SerialNumber,
	}, nil
}

func (pc *lotteryPrizeCardUsecase) GetPrizeCardCollection(ctx context.Context, userId uint, prizePool *domain.PrizePool) (*domain.PrizeCardCollection, error) {
	poolPrizeCard, err := pc.prizeCardRepo.Gets(ctx, &domain.PrizeCard{PoolId: prizePool.Model.ID})
	if err != nil {
		return &domain.PrizeCardCollection{}, err
	}

	var poolPrizeCardIdSlice []uint
	for _, v := range poolPrizeCard {
		poolPrizeCardIdSlice = append(poolPrizeCardIdSlice, v.Model.ID)
	}

	userPrizeCardList, err := pc.prizeCardRepo.GetUserPrizeCardList(ctx, &domain.UserPrizeCard{UserId: userId}, map[string]interface{}{
		`card_id in (?)`: poolPrizeCardIdSlice,
	})
	if err != nil {
		return &domain.PrizeCardCollection{}, err
	}

	userPrizeCardIdMap := make(map[uint]string)
	for _, v := range userPrizeCardList {
		userPrizeCardIdMap[v.CardId] = v.SerialNumber
	}

	var userCollection []domain.PrizeCardOutput
	for _, v := range poolPrizeCard {
		if serialNumber, ok := userPrizeCardIdMap[v.Model.ID]; ok {
			userCollection = append(userCollection, domain.PrizeCardOutput{
				ID:           v.Model.ID,
				Picture:      v.Picture,
				Title:        v.Title,
				Description:  v.Description,
				Audio:        v.Audio,
				SerialNumber: serialNumber,
			})
		}

		if len(userCollection) >= len(userPrizeCardIdMap) {
			break
		}
	}

	return &domain.PrizeCardCollection{
		PoolName: prizePool.Name,
		Cards:    userCollection,
		Total:    int64(len(poolPrizeCard)),
	}, nil
}

func (pc *lotteryPrizeCardUsecase) GetPrizeCardCollectionDetail(ctx context.Context, userPrizeCard *domain.UserPrizeCard) (string, *domain.Coupon, *domain.PrizeCard, error) {
	userPrizeCard, err := pc.prizeCardRepo.GetUserPrizeCard(ctx, userPrizeCard)
	if err != nil {
		return "", nil, nil, err
	}

	coupon, err := pc.prizeCardRepo.GetCoupon(ctx, &domain.Coupon{CardId: userPrizeCard.CardId, UserId: userPrizeCard.UserId})
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", nil, nil, err
	}

	card, err := pc.prizeCardRepo.Get(ctx, &domain.PrizeCard{Model: gorm.Model{ID: userPrizeCard.CardId}})
	if err != nil {
		return "", nil, nil, err
	}

	return userPrizeCard.SerialNumber, coupon, card, nil
}
