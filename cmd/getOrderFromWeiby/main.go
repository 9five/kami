package main

import (
	"context"
	"fmt"
	"kami/config"
	_kamiOrderRepo "kami/kamiOrder/repository/postgresql"
	_kamiOrderUsecase "kami/kamiOrder/usecase"
	_weibyScheduler "kami/weiby/delivery/scheduler"
	_weibyRepo "kami/weiby/repository/curl"
	_weibyUsecase "kami/weiby/usecase"
	"time"
)

func main() {
	ctx := context.Background()

	kamiOrderRepo := _kamiOrderRepo.NewPostgresqlKamiOrderRepository(config.DB)
	kamiOrderUsecase := _kamiOrderUsecase.NewKamiOrderUsecase(kamiOrderRepo)
	weibyRepo := _weibyRepo.NewCurlWeibyRepository(config.WeibyVendorId, config.WeibyApiKey)
	weibyUsecase := _weibyUsecase.NewWeibyUsecase(weibyRepo)
	scheduler := _weibyScheduler.NewWeibyHandler(weibyUsecase, kamiOrderUsecase)

	ol, err := scheduler.GetAllStoreOrderList(ctx, time.Now().AddDate(0, 0, -1).Format("2006-01-02"), time.Now().Format("2006-01-02"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range ol {
		scheduler.StoreOrder(ctx, k, v)
	}
}
