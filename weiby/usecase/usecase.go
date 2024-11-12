package usecase

import (
	"context"
	"kami/domain"
)

type weibyUsecase struct {
	weibyRepo domain.WeibyRepository
}

func NewWeibyUsecase(weibyRepo domain.WeibyRepository) domain.WeibyUsecase {
	return &weibyUsecase{
		weibyRepo: weibyRepo,
	}
}

func (w *weibyUsecase) GetStoreList(ctx context.Context) ([]*domain.WeibyStoreInfo, error) {
	return w.weibyRepo.GetStoreList(ctx)
}

func (w *weibyUsecase) GetStore(ctx context.Context, pid string) (*domain.WeibyStoreInfo, error) {
	return w.weibyRepo.GetStore(ctx, pid)
}

func (w *weibyUsecase) GetOrderList(ctx context.Context, pid, startTime, endTime string) (*domain.OrderList, error) {
	storeInfo, err := w.weibyRepo.GetStore(ctx, pid)
	if err != nil {
		return &domain.OrderList{}, err
	}

	return w.weibyRepo.GetOrderList(ctx, pid, startTime, endTime, storeInfo.PartnerType)
}
