package main

import (
	"context"
	"fmt"
	"kami/config"
	_kamiOrderScheduler "kami/kamiOrder/delivery/scheduler"
	_kamiOrderRepo "kami/kamiOrder/repository/postgresql"
	_kamiOrderUsecase "kami/kamiOrder/usecase"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	domain = `https://foodpanda.portal.restaurant`
	//ckitchen account
	ac = `ckitchen@capsulecorporation.cc`
	pw = `capsule1111`
)

var (
	configCtx               context.Context
	configCtxCancel         context.CancelFunc
	configOptionalCtx       context.Context
	configOptionalCtxCancel context.CancelFunc
)

func init() {
	// create chrome instance
	options := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	// create context
	configOptionalCtx, configOptionalCtxCancel = chromedp.NewExecAllocator(context.Background(), options...)
	configCtx, configCtxCancel = chromedp.NewContext(configOptionalCtx)
}

func loginPlatform(ctx context.Context) error {
	url := domain + `/login?redirect=/`

	return chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.SendKeys(`#login-email-field`, ac),
		chromedp.SendKeys(`#login-password-field`, pw),
		chromedp.Click(`#button_login`, chromedp.NodeVisible),
		chromedp.Click(`#root > div > div > div > div.MuiBox-root.css-14iaz0u > div.MuiBox-root.css-15jayxm > div > div.MuiBox-root.css-4cxybv > div.MuiBox-root.css-h5fkc8 > button.MuiButton-root.MuiButton-text.MuiButton-textPrimary.MuiButton-sizeMedium.MuiButton-textSizeMedium.MuiButton-disableElevation.MuiButtonBase-root.css-1ss1yo0`, chromedp.NodeVisible),
	)
}

func getOrderList(ctx context.Context) ([]*cdp.Node, error) {
	var nodes []*cdp.Node

	url := domain + `/orders?int_ref=side-nav`
	orderRefreshButton := `button[aria-label="Refresh Page"]`
	orderList := `#sideBarV2 > div.MuiBox-root.css-c31z9o > div.MuiBox-root.css-oilzl3 > main > div.plugin-muiv4-MuiPaper-root.sc-iGrrsa.jorOXY.plugin-muiv4-MuiTableContainer-root.plugin-muiv4-MuiPaper-elevation2.plugin-muiv4-MuiPaper-rounded > div > div.plugin-muiv4-MuiPaper-root.plugin-muiv4-MuiTableBody-root.plugin-muiv4-MuiPaper-elevation1.plugin-muiv4-MuiPaper-rounded`

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// chromedp.Navigate(`https://foodpanda.portal.restaurant/orders?int_ref=side-nav&from=2022-11-13&to=2022-11-13&billableFilterState=ALL`),
		chromedp.Click(orderRefreshButton, chromedp.NodeVisible),
		chromedp.Nodes(orderList, &nodes, chromedp.ByQuery, chromedp.AtLeast(0)),
	)

	return nodes, err
}

func main() {
	fmt.Println("start!")
	defer configOptionalCtxCancel()
	defer configCtxCancel()

	if err := loginPlatform(configCtx); err != nil {
		fmt.Println(err.Error())
		return
	}

	nodes, err := getOrderList(configCtx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(nodes) != 0 {
		kamiOrderRepo := _kamiOrderRepo.NewPostgresqlKamiOrderRepository(config.DB)
		kamiOrderUsecase := _kamiOrderUsecase.NewKamiOrderUsecase(kamiOrderRepo)

		scheduler := _kamiOrderScheduler.NewKamiOrderHandler(kamiOrderUsecase)
		if err := scheduler.FetchOrdersFromFoodPanda(configCtx, nodes); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		fmt.Println("order not found")
	}
}
