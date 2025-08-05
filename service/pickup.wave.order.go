package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bestk/dmxsmart-client/model"
)

// PickupWaveService 拣货波次服务
type PickupWaveService struct {
	client *Client
}

// NewPickupWaveService 创建拣货波次服务
func NewPickupWaveService(client *Client) *PickupWaveService {
	return &PickupWaveService{client: client}
}

// GetWaitingPickOrders retrieves the list of waiting pick orders
func (s *PickupWaveService) GetWaitingPickOrders(page, pageSize int, customerIds []int, keyword string) (*model.WaitingPickOrderResponse, error) {
	urlStr := "/api/tenant/outbound/pickupwave/listWaitingPickOrder"

	params := url.Values{}
	params.Set("current", fmt.Sprintf("%d", page))
	params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	params.Set("warehouseId", fmt.Sprintf("%d", s.client.config.WarehouseID))
	params.Set("keywordType", "referenceId")
	params.Set("timeType", "createTime")
	params.Set("keyword", keyword)

	for _, customerId := range customerIds {
		params.Add("customerIds[]", fmt.Sprintf("%d", customerId))
	}

	fullURL := fmt.Sprintf("%s?%s", urlStr, params.Encode())

	body := map[string]interface{}{
		"keyword": keyword,
	}

	var result model.WaitingPickOrderResponse

	err := s.client.makeRequest(http.MethodPost, fullURL, body, &result)

	if err != nil {
		return nil, fmt.Errorf("failed to get waiting pick orders: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("get waiting pick orders failed: %s", result.ErrorMessage)
	}

	return &result, nil
}

// CreatePickupWave creates a new pickup wave
func (s *PickupWaveService) CreatePickupWave(isAll bool, pickupType int, isOutbound bool, customerIds []int, remark string) (*model.CreatePickupWaveResponse, error) {
	urlStr := "/api/tenant/outbound/pickupwave/createPickupWave"

	params := url.Values{}
	params.Set("isAll", fmt.Sprintf("%t", isAll))
	params.Set("pickupType", fmt.Sprintf("%d", pickupType))
	params.Set("isOutbound", fmt.Sprintf("%t", isOutbound))
	params.Set("warehouseId", fmt.Sprintf("%d", s.client.config.WarehouseID))
	params.Set("remark", remark)

	for _, customerId := range customerIds {
		params.Add("customerIds[]", fmt.Sprintf("%d", customerId))
	}

	fullURL := fmt.Sprintf("%s?%s", urlStr, params.Encode())

	var result model.CreatePickupWaveResponse

	err := s.client.makeRequest(http.MethodPost, fullURL, nil, &result)

	if err != nil {
		return nil, fmt.Errorf("failed to create pickup wave: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("create pickup wave failed: %s", result.ErrorMessage)
	}

	return &result, nil
}
