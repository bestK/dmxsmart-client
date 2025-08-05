package client

import (
	"fmt"

	"github.com/bestk/dmxsmart-client/config"
	"github.com/bestk/dmxsmart-client/logger"
	"github.com/bestk/dmxsmart-client/service"
)

// DMXSmartClient
type DMXSmartClient struct {
	services *service.Services
	config   *config.ConfigStruct
}

// NewDMXSmartClient 创建新的DMXSmart客户端
func NewDMXSmartClient(configPath string) (*DMXSmartClient, error) {

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	if cfg.Account == "" || cfg.Password == "" {
		return nil, fmt.Errorf("Account or password is empty")
	}

	// 初始化日志
	logger.Init()

	// 创建服务
	services := service.NewServices(cfg)
	services.SetLogger(logger.Logger)

	return &DMXSmartClient{
		services: services,
		config:   cfg,
	}, nil
}
