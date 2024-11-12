package domain

import (
	"context"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PrizePool struct {
	gorm.Model
	Owner  string         `json:"owner"`
	Name   string         `json:"name"`
	Banner string         `json:"banner"`
	Points int64          `json:"points"`
	Style  datatypes.JSON `json:"style"`
}

type PrizePoolRepository interface {
	New(ctx context.Context, prizePool *PrizePool) (*PrizePool, error)
	Get(ctx context.Context, prizePool *PrizePool, optsWhere ...map[string]interface{}) (*PrizePool, error)
	Gets(ctx context.Context, prizePool *PrizePool, optsWhere ...map[string]interface{}) ([]*PrizePool, error)
}

type PrizePoolUsecase interface {
	GetPrizePool(ctx context.Context, prizePool *PrizePool) (*PrizePool, error)
	GetPrizePoolList(ctx context.Context, prizePool *PrizePool) ([]*PrizePool, error)
	SubtractUserPoints(ctx context.Context, user *KamiUser, prizePool *PrizePool) error
}

type PrizeCard struct {
	gorm.Model
	PoolId      uint   `json:"pool_id"`
	Picture     string `json:"picture"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Audio       string `json:"audio"`
	Probability int64  `json:"probability"`
	Alert       bool   `json:"alert"`
}

type PrizeCardRepository interface {
	New(ctx context.Context, prizeCard *PrizeCard) (*PrizeCard, error)
	Get(ctx context.Context, prizeCard *PrizeCard, optsWhere ...map[string]interface{}) (*PrizeCard, error)
	Gets(ctx context.Context, prizeCard *PrizeCard, optsWhere ...map[string]interface{}) ([]*PrizeCard, error)
	NewUserPrizeCard(ctx context.Context, userPrizeCard *UserPrizeCard) (*UserPrizeCard, error)
	GetUserPrizeCard(ctx context.Context, userPrizeCard *UserPrizeCard, optsWhere ...map[string]interface{}) (*UserPrizeCard, error)
	GetUserPrizeCardList(ctx context.Context, userPrizeCard *UserPrizeCard, optsWhere ...map[string]interface{}) (result []*UserPrizeCard, err error)
	GetCoupon(ctx context.Context, coupon *Coupon, optsWhere ...map[string]interface{}) (result *Coupon, err error)
	UpdateCoupon(ctx context.Context, coupon *Coupon) error
}

type PrizeCardUsecase interface {
	GetPrizeCard(ctx context.Context, prizeCard *PrizeCard) (result *PrizeCard, err error)
	GetPrizeCardList(ctx context.Context, prizeCard *PrizeCard) ([]*PrizeCard, error)
	GetWeightedRandomList(ctx context.Context, userId uint, prizeCardList []*PrizeCard) ([]*PrizeCard, error)
	Draw(ctx context.Context, userId uint, prizeCardList []*PrizeCard) (*PrizeCardOutput, error)
	GetPrizeCardCollection(ctx context.Context, userId uint, prizePool *PrizePool) (*PrizeCardCollection, error)
	GetPrizeCardCollectionDetail(ctx context.Context, userPrizeCard *UserPrizeCard) (string, *Coupon, *PrizeCard, error)
}

type UserPrizeCard struct {
	ID           uint   `json:"id"`
	UserId       uint   `json:"User_id"`
	CardId       uint   `json:"card_id"`
	SerialNumber string `json:"serial_number"`
}

type Coupon struct {
	ID         uint      `json:"id"`
	CardId     uint      `json:"card_id"`
	UserId     uint      `json:"user_id"`
	Code       string    `json:"code"`
	ExpiryDate time.Time `json:"expiry_date"`
}

type PrizeCardOutput struct {
	ID           uint           `json:"id"`
	Picture      string         `json:"picture"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Audio        string         `json:"audio"`
	SerialNumber string         `json:"serial_number"`
	Style        datatypes.JSON `json:"style"`
}

type PrizeCardCollection struct {
	PoolName string            `json:"pool_name"`
	Cards    []PrizeCardOutput `json:"cards"`
	Total    int64             `json:"total"`
}

type PrizeCardCollectionDetail struct {
	Coupon       Coupon          `json:"coupon"`
	Card         PrizeCardOutput `json:"card"`
	SerialNumber string          `json:"serial_number"`
}
