package client

import (
	"path/filepath"
	"testing"

	"github.com/bestk/dmxsmart-client/config"
)

func TestLoginWithAutoOCR(t *testing.T) {
	configPath := filepath.Join(".", "config.yaml")

	client, err := NewDMXSmartClient(configPath)
	if err != nil {
		t.Fatalf("Failed to create DMXSmartClient: %v", err)
	}

	// 执行登录测试
	resp, err := client.services.Auth.LoginWithAutoOCR(client.config.Account, client.config.Password)
	if err != nil {
		t.Errorf("LoginWithAutoOCR() error = %v", err)
		return
	}

	// 验证响应
	if !resp.Success {
		t.Errorf("登录失败: %s", resp.ErrorMessage)
		return
	}

	// 验证token
	if resp.Data.Token == "" {
		t.Error("登录成功但未获取到token")
		return
	}

	t.Logf("登录成功，token: %s", resp.Data.Token)
}

func TestValidateSession(t *testing.T) {
	configPath := filepath.Join(".", "config.yaml")

	client, err := NewDMXSmartClient(configPath)
	if err != nil {
		t.Fatalf("Failed to create DMXSmartClient: %v", err)
	}

	err = client.services.Auth.ValidateSession()
	if err != nil {
		t.Errorf("ValidateSession() error = %v", err)
		return
	}

	t.Log("ValidateSession() success")
}

func TestGetWaitingPickOrders(t *testing.T) {
	configPath := filepath.Join(".", "config.yaml")

	client, err := NewDMXSmartClient(configPath)
	if err != nil {
		t.Fatalf("Failed to create DMXSmartClient: %v", err)
	}

	resp, err := client.services.PickupWave.GetWaitingPickOrders(1, 20, config.GlobalConfig.CustomerIDs, "")
	if err != nil {
		t.Errorf("GetWaitingPickOrders() error = %v", err)
		return
	}

	t.Logf("GetWaitingPickOrders() success, total: %d, data: %+v", resp.Total, resp.Data)
}

func TestCreatePickupWave(t *testing.T) {
	configPath := filepath.Join(".", "config.yaml")

	client, err := NewDMXSmartClient(configPath)
	if err != nil {
		t.Fatalf("Failed to create DMXSmartClient: %v", err)
	}

	resp, err := client.services.PickupWave.CreatePickupWave(true, 1, true, config.GlobalConfig.CustomerIDs, "[BOT]")
	if err != nil {
		t.Errorf("CreatePickupWave() error = %v", err)
		return
	}

	t.Logf("CreatePickupWave() success, id: %d", resp.Data.ID)
}
