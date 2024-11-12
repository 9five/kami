package scheduler

import (
	"context"
	// "encoding/json"
	// "fmt"
	"github.com/chromedp/cdproto/cdp"
	"kami/domain"
)

type KamiOrderHandler struct {
	KamiOrderUsecase domain.KamiOrderUsecase
}

func NewKamiOrderHandler(orderInfoUsecase domain.KamiOrderUsecase) *KamiOrderHandler {
	return &KamiOrderHandler{
		KamiOrderUsecase: orderInfoUsecase,
	}
}

func (o *KamiOrderHandler) FetchOrdersFromFoodPanda(ctx context.Context, nodes []*cdp.Node) error {
	orders, err := o.KamiOrderUsecase.OrderCrawler(ctx, nodes)
	if err != nil {
		return err
	}
	orders, err = o.KamiOrderUsecase.GetMoreOrderDetail(ctx, orders)
	if err != nil {
		return err
	}
	err = o.KamiOrderUsecase.BatchStore(ctx, "FOODPANDA", orders)
	if err != nil {
		return err
	}
	// ordersJson, err := json.Marshal(orders)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(ordersJson))
	return nil
}
