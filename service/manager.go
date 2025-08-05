package service

import (
	"github.com/bestk/dmxsmart-client/config"
	"github.com/go-resty/resty/v2"
)

// Services 聚合所有API服务
type Services struct {
	client     *Client
	Auth       *AuthService
	PickupWave *PickupWaveService
}

// NewServices 创建服务聚合器
func NewServices(config *config.ConfigStruct) *Services {
	client := NewClient(config)

	return &Services{
		client:     client,
		Auth:       NewAuthService(client),
		PickupWave: NewPickupWaveService(client),
	}
}

// SetLogger 为所有服务设置日志记录器
func (s *Services) SetLogger(logger resty.Logger) {
	s.client.SetLogger(logger)
}

// GetClient 获取底层HTTP客户端
func (s *Services) GetClient() *Client {
	return s.client
}
