package usecase

import (
	"context"
	"errors"
	"fmt"
	"kami/domain"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"gorm.io/gorm"
)

type kamiOrderUsecase struct {
	kamiOrderRepo domain.KamiOrderRepository
}

func NewKamiOrderUsecase(kamiOrderRepo domain.KamiOrderRepository) domain.KamiOrderUsecase {
	return &kamiOrderUsecase{
		kamiOrderRepo: kamiOrderRepo,
	}
}

func (o *kamiOrderUsecase) OrderCrawler(ctx context.Context, nodes []*cdp.Node) ([]*domain.KamiOrder, error) {
	var orders []*domain.KamiOrder
	loc, _ := time.LoadLocation("Asia/Taipei")

	for _, node := range nodes {

		if err := chromedp.Run(ctx, chromedp.Nodes(`a`, &node.Children, chromedp.ByQueryAll, chromedp.FromNode(node))); err != nil {
			return orders, err
		}

		for _, children := range node.Children {
			order := struct {
				orderId       string
				restaurant    string
				status        string
				billingStatus string
				orderPlacedAt string
			}{}

			err := chromedp.Run(ctx,
				chromedp.TextContent(`div:nth-child(1) > p`, &order.orderId, chromedp.ByQuery, chromedp.FromNode(children)),
				chromedp.TextContent(`div:nth-child(2) > p`, &order.restaurant, chromedp.ByQuery, chromedp.FromNode(children)),
				chromedp.TextContent(`div:nth-child(3) > div > p`, &order.status, chromedp.ByQuery, chromedp.FromNode(children)),
				chromedp.TextContent(`div:nth-child(4) > p`, &order.billingStatus, chromedp.ByQuery, chromedp.FromNode(children)),
				chromedp.TextContent(`div:nth-child(5) > p`, &order.orderPlacedAt, chromedp.ByQuery, chromedp.FromNode(children)),
			)
			if err != nil {
				return orders, err
			}

			orderPlacedAt, _ := time.ParseInLocation("01/02/2006, 3:04 PM", order.orderPlacedAt, loc)

			orders = append(orders, &domain.KamiOrder{
				OrderId:       order.orderId,
				Restaurant:    order.restaurant,
				Status:        order.status,
				BillingStatus: order.billingStatus,
				OrderPlacedAt: orderPlacedAt,
			})
		}
	}
	return orders, nil
}

func (o *kamiOrderUsecase) GetMoreOrderDetail(ctx context.Context, orders []*domain.KamiOrder) ([]*domain.KamiOrder, error) {
	for k, order := range orders {
		if order.Status == "Delivered" {
			orderElement := fmt.Sprintf(`a[data-testid="master-order-desktop-item-%s"]`, order.OrderId)

			chromedp.Run(ctx, chromedp.Tasks{clickDetailUntilDivOpenFunc(orderElement, order)})
			orders[k] = order

			if err := chromedp.Run(ctx, chromedp.Tasks{clickBackBtnUntilDivCloseFunc()}); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	return orders, nil
}

func clickDetailUntilDivOpenFunc(oe string, o *domain.KamiOrder) chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var nodes []*cdp.Node
		var orderDeliveredAt string
		loc, _ := time.LoadLocation("Asia/Taipei")

		for {
			chromedp.Click(oe, chromedp.NodeVisible).Do(ctx)
			chromedp.WaitReady(`div[role="presentation"]`).Do(ctx)
			chromedp.WaitReady(`div[data-testid="master-order-details-status-delivered"]`).Do(ctx)
			chromedp.Nodes(`div[data-testid="master-order-details-status-delivered"]`, &nodes, chromedp.ByQuery, chromedp.AtLeast(0)).Do(ctx)
			if len(nodes) != 0 {
				break
			}
		}
		chromedp.TextContent(`div > p`, &orderDeliveredAt, chromedp.ByQuery, chromedp.FromNode(nodes[0])).Do(ctx)

		if orderDeliveredAt != "" {
			orderPlacedAt, _ := time.ParseInLocation("01/02/2006, 3:04 PM", time.Now().Format("01/02/2006, ")+orderDeliveredAt, loc)
			o.OrderDeliveredAt = orderPlacedAt
		}
		return
	}
}

func clickBackBtnUntilDivCloseFunc() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {
		var nodes []*cdp.Node
		for {
			chromedp.Click(`button[data-testid="master-order-detail-back-btn"]`, chromedp.AtLeast(0)).Do(ctx)
			err = chromedp.Nodes(`div[role="presentation"]`, &nodes, chromedp.ByQuery, chromedp.AtLeast(0)).Do(ctx)
			if len(nodes) == 0 || err != nil {
				break
			}
		}
		return
	}
}

func (o *kamiOrderUsecase) BatchStore(ctx context.Context, platform string, orders []*domain.KamiOrder) error {
	for _, order := range orders {
		order.Platform = platform

		if err := o.kamiOrderRepo.Store(ctx, order); err != nil {
			return err
		}
	}

	return nil
}

func (o *kamiOrderUsecase) CheckOrderInput(ctx context.Context, orderInput *domain.OrderInput) (*domain.KamiOrder, error) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	orderDeliveredAt, err := time.ParseInLocation("2006-01-02 15:04", orderInput.OrderDeliveredAt, loc)
	if err != nil {
		return &domain.KamiOrder{}, err
	}

	order := &domain.KamiOrder{
		OrderId:          fmt.Sprintf("%s-%s", orderInput.Prefix, orderInput.Suffix),
		OrderDeliveredAt: orderDeliveredAt,
	}

	if order, err = o.kamiOrderRepo.Get(ctx, order); err != nil {
		if err == gorm.ErrRecordNotFound {
			return order, errors.New("order not found")
		}
	}

	return order, nil
}

func (o *kamiOrderUsecase) RegisterOrder(ctx context.Context, order *domain.KamiOrder, user *domain.KamiUser) error {
	order.OwnerPhone = user.Phone
	if err := o.kamiOrderRepo.Update(ctx, order); err != nil {
		return err
	}

	user.Points += 1

	return nil
}

func (o *kamiOrderUsecase) GetOrderList(ctx context.Context, order *domain.KamiOrder) ([]*domain.KamiOrder, error) {
	return o.kamiOrderRepo.Gets(ctx, order)
}

func (o *kamiOrderUsecase) GetOrderListByDate(ctx context.Context, order *domain.KamiOrder, startDate string, endDate string) ([]*domain.KamiOrder, error) {
	sd, err := time.Parse("2006-01", startDate)
	if err != nil {
		return []*domain.KamiOrder{}, err
	}

	ed, err := time.Parse("2006-01", endDate)
	if err != nil {
		return []*domain.KamiOrder{}, err
	}
	ed = ed.AddDate(0, 1, 0)

	return o.kamiOrderRepo.Gets(ctx, order, map[string]interface{}{
		"order_delivered_at >= ?": sd,
		"order_delivered_at <= ?": ed,
	})
}
