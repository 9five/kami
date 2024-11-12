package domain

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"gorm.io/gorm"
)

type KamiOrder struct {
	gorm.Model
	OrderId          string    `json:"order_id"`
	Restaurant       string    `json:"restaurant"`
	Status           string    `json:"status"`
	BillingStatus    string    `json:"billing_status"`
	OrderPlacedAt    time.Time `json:"order_placed_at"`
	OrderDeliveredAt time.Time `json:"order_delivered_at"`
	Platform         string    `json:"platform"`
	OwnerPhone       string    `json:"owner_phone"`
}

type KamiOrderRepository interface {
	Store(ctx context.Context, order *KamiOrder) error
	Get(ctx context.Context, order *KamiOrder, optsWhere ...map[string]interface{}) (*KamiOrder, error)
	Gets(ctx context.Context, order *KamiOrder, optsWhere ...map[string]interface{}) ([]*KamiOrder, error)
	Update(ctx context.Context, order *KamiOrder) error
}

type KamiOrderUsecase interface {
	OrderCrawler(ctx context.Context, nodes []*cdp.Node) ([]*KamiOrder, error)
	GetMoreOrderDetail(ctx context.Context, orders []*KamiOrder) ([]*KamiOrder, error)
	BatchStore(ctx context.Context, platform string, orders []*KamiOrder) error
	CheckOrderInput(ctx context.Context, orderInput *OrderInput) (*KamiOrder, error)
	RegisterOrder(ctx context.Context, order *KamiOrder, user *KamiUser) error
	GetOrderList(ctx context.Context, order *KamiOrder) ([]*KamiOrder, error)
	GetOrderListByDate(ctx context.Context, order *KamiOrder, startDate string, endDate string) ([]*KamiOrder, error)
}

type OrderInput struct {
	Prefix           string `json:"prefix"`
	Suffix           string `json:"suffix"`
	OrderDeliveredAt string `json:"order_delivered_at"`
}
