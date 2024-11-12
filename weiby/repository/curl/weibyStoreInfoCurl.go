package curl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kami/domain"
	"net/http"
)

type curlWeibyRepository struct {
	vendorId string
	apiKey   string
	client   *http.Client
}

func NewCurlWeibyRepository(vendorId string, apiKey string) domain.WeibyRepository {
	return &curlWeibyRepository{
		vendorId: vendorId,
		apiKey:   apiKey,
		client:   &http.Client{},
	}
}

func (c *curlWeibyRepository) GetStoreList(ctx context.Context) (result []*domain.WeibyStoreInfo, err error) {
	weibyGetStoreListAPI := fmt.Sprintf(`https://wis.weibyapps.com/%s/api/v1/stores`, c.vendorId)

	r, err := http.NewRequest("GET", weibyGetStoreListAPI, nil)
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-api-key", c.apiKey)

	response, err := c.client.Do(r)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return
}

func (c *curlWeibyRepository) GetStore(ctx context.Context, pid string) (result *domain.WeibyStoreInfo, err error) {
	weibyGetStoreAPI := fmt.Sprintf(`https://wis.weibyapps.com/%s/api/v1/stores/%s`, c.vendorId, pid)

	r, err := http.NewRequest("GET", weibyGetStoreAPI, nil)
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-api-key", c.apiKey)

	response, err := c.client.Do(r)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return
}

func (c *curlWeibyRepository) GetOrderList(ctx context.Context, pid, startTime, endTime string, ptype int) (result *domain.OrderList, err error) {
	weibyGetOrderListAPI := fmt.Sprintf(`https://wis.weibyapps.com/%s/api/v1/detail/%s/orders?start_time=%s&end_time=%s`, c.vendorId, pid, startTime, endTime)

	r, err := http.NewRequest("GET", weibyGetOrderListAPI, nil)
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("x-api-key", c.apiKey)

	response, err := c.client.Do(r)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var list domain.OrderList
	if ptype == 1 {
		err = json.Unmarshal(body, &list.UberEats)
	} else {
		err = json.Unmarshal(body, &list.Foodpanda)
	}

	result = &list
	return
}
