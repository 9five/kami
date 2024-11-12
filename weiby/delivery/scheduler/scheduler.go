package scheduler

import (
	"context"
	"fmt"
	"kami/domain"
	"time"
)

type WeibyHandler struct {
	WeibyUsecase     domain.WeibyUsecase
	KamiOrderUsecase domain.KamiOrderUsecase
}

func NewWeibyHandler(weibyUsecase domain.WeibyUsecase, kamiOrderUsecase domain.KamiOrderUsecase) *WeibyHandler {
	return &WeibyHandler{
		WeibyUsecase:     weibyUsecase,
		KamiOrderUsecase: kamiOrderUsecase,
	}
}

func (w *WeibyHandler) GetAllStoreOrderList(ctx context.Context, startTime string, endTime string) (map[string]domain.OrderList, error) {
	storeList, err := w.WeibyUsecase.GetStoreList(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]domain.OrderList)
	for _, v := range storeList {
		orderList, err := w.WeibyUsecase.GetOrderList(ctx, v.PartnerId, startTime, endTime)
		if err != nil {
			return nil, err
		}

		result[v.Name] = *orderList
	}
	return result, nil
}

func (w *WeibyHandler) StoreOrder(ctx context.Context, store string, ol domain.OrderList) (err error) {
	var orders []*domain.KamiOrder
	if len(ol.UberEats) != 0 {
		orders, err = uberEatsProcess(ol.UberEats)
		if err != nil {
			return err
		}
		fmt.Println(orders)
		// err = w.KamiOrderUsecase.BatchStore(ctx, "UBEREATS", orders)
		// if err != nil {
		// 	return err
		// }
	} else if len(ol.Foodpanda) != 0 {
		orders, err := foodpandaProcess(store, ol.Foodpanda)
		if err != nil {
			return err
		}
		fmt.Println(orders)
		// err = w.KamiOrderUsecase.BatchStore(ctx, "FOODPANDA", orders)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func uberEatsProcess(ol []domain.UberEatsOrder) ([]*domain.KamiOrder, error) {
	var result []*domain.KamiOrder

	for _, v := range ol {
		placeAt, err := time.Parse("RFC3339", v.PlacedAt)
		if err != nil {
			return result, err
		}

		deliveredAt, err := time.Parse("RFC3339", v.EstimatedReadyForPickupAt)
		if err != nil {
			return result, err
		}

		var res domain.KamiOrder
		res.OrderId = v.Id
		res.Restaurant = v.Store.Name
		res.OrderPlacedAt = placeAt
		res.OrderDeliveredAt = deliveredAt

		result = append(result, &res)
	}
	return result, nil
}

func foodpandaProcess(store string, ol []domain.FoodpandaOrder) ([]*domain.KamiOrder, error) {
	var result []*domain.KamiOrder

	for _, v := range ol {
		var res domain.KamiOrder
		res.OrderId = v.Code
		res.Restaurant = store
		res.OrderPlacedAt = v.CreatedAt
		res.OrderDeliveredAt = v.Delivery.ExpectedDeliveryTime

		result = append(result, &res)
	}
	return result, nil
}
