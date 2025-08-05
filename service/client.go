package service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/bestk/dmxsmart-client/config"
	"github.com/bestk/dmxsmart-client/model"
	"github.com/go-resty/resty/v2"
)

const (
	BaseURL = "https://wms.dmxsmart.com"
)

// Client represents a DMXSmart API client
type Client struct {
	httpClient *resty.Client
	config     *config.ConfigStruct
}

// NewClient creates a new DMXSmart client
func NewClient(config *config.ConfigStruct) *Client {
	client := &Client{
		httpClient: resty.New().
			SetDebug(config.Debug).
			SetBaseURL(BaseURL).
			// 设置超时
			SetTimeout(time.Duration(config.Timeout) * time.Second).
			// 设置重试
			SetRetryCount(3).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(20 * time.Second).
			// 设置TLS配置
			SetTLSClientConfig(&tls.Config{
				MinVersion: tls.VersionTLS12,
			}),
		config: config,
	}

	// Set default headers
	client.httpClient.SetHeaders(map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"Authorization":   "Bearer " + config.AccessToken,
		"Content-Type":    "application/json",
	})

	return client
}

// SetLogger sets the logger for the HTTP client
func (c *Client) SetLogger(logger resty.Logger) {
	c.httpClient.SetLogger(logger)
}

// UpdateToken updates the authorization token
func (c *Client) UpdateToken(token string) {
	c.httpClient.SetHeader("Authorization", "Bearer "+token)
	c.config.AccessToken = token
}

// makeRequest 统一的请求处理方法
func (c *Client) makeRequest(method, url string, body interface{}, result interface{}) error {
	req := c.httpClient.R()

	if body != nil {
		req.SetBody(body)
	}

	if result != nil {
		req.SetResult(result)
	}

	var resp *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = req.Get(url)
	case http.MethodPost:
		resp, err = req.Post(url)
	case http.MethodPut:
		resp, err = req.Put(url)
	case http.MethodDelete:
		resp, err = req.Delete(url)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

// checkResponse 检查API响应是否成功
func (c *Client) checkResponse(resp interface{}) error {
	if respData, ok := resp.(*model.Response); ok {
		if !respData.Success {
			return fmt.Errorf("API request failed: %s", respData.ErrorMessage)
		}
	}
	return nil
}
