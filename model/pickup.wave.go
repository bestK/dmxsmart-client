package model

// WaitingPickOrder 等待拣货订单
type WaitingPickOrder struct {
	ID int32 `json:"id"`
}

// WaitingPickOrderResponse 等待拣货订单响应
type WaitingPickOrderResponse struct {
	Response
	Data []WaitingPickOrder `json:"data"`
}

// CreatePickupWaveResponse 创建拣货波次响应
type CreatePickupWaveResponse struct {
	Response
	Data struct {
		ID int32 `json:"id"`
	} `json:"data"`
} 