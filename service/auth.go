package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bestk/dmxsmart-client/model"
	"github.com/bestk/dmxsmart-client/ocr"
)

// AuthService 认证服务
type AuthService struct {
	client *Client
}

// NewAuthService 创建认证服务
func NewAuthService(client *Client) *AuthService {
	return &AuthService{client: client}
}

// ValidateSession validates the session
func (s *AuthService) ValidateSession() error {
	url := "/api/user/getUserInfo"

	var result model.Response
	err := s.client.makeRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to validate session: %w", err)
	}

	return s.client.checkResponse(&result)
}

// GetCaptcha retrieves a captcha image
func (s *AuthService) GetCaptcha() (*model.CaptchaResponse, error) {
	url := "/api/login/captcha"

	// 临时移除Authorization头
	authHeader := s.client.httpClient.Header.Get("Authorization")
	if authHeader != "" {
		s.client.httpClient.Header.Del("Authorization")
		defer s.client.httpClient.SetHeader("Authorization", authHeader)
	}

	var result model.CaptchaResponse

	resp, err := s.client.httpClient.R().
		SetQueryParam("lang", "zh-CN").
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", BaseURL)).
		SetResult(&result).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get captcha: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &result, nil
}

// Login 执行登录操作
func (s *AuthService) Login(username, password, captcha, uuid string) (*model.LoginResponse, error) {
	url := "/api/login/authenticate"

	// 加密密码
	encryptedPassword, err := encryptPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// 构建登录请求
	loginReq := model.LoginRequest{
		Username:    username,
		Password:    encryptedPassword,
		Captcha:     captcha,
		UUID:        uuid,
		LoginType:   "USERNAME",
		DeviceToken: nil,
		Lang:        "zh-CN",
	}

	var result model.LoginResponse
	// 发送登录请求
	resp, err := s.client.httpClient.R().
		SetHeader("Accept", "application/json, text/plain, */*").
		SetHeader("Referer", fmt.Sprintf("%s/user/login", BaseURL)).
		SetHeader("Cookie", "locale=zh-CN").
		SetBody(loginReq).
		SetResult(&result).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send login request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// 如果登录成功，更新客户端的 token
	if result.Success {
		s.client.UpdateToken(result.Data.Token)
	}

	return &result, nil
}

// LoginWithAutoOCR 自动识别验证码并登录，失败时最多重试3次
func (s *AuthService) LoginWithAutoOCR(username, password string) (*model.LoginResponse, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 获取验证码
		captchaResp, err := s.GetCaptcha()
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to get captcha: %w", attempt+1, err)
			continue
		}

		// 从base64图片中提取图片数据
		imgData := captchaResp.Data.Img
		imgData = strings.TrimPrefix(imgData, "data:image/png;base64,")
		imgData = strings.TrimPrefix(imgData, "data:image/jpeg;base64,")

		// 识别验证码
		captchaText, err := ocr.RecognizeBase64Image(imgData)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to recognize captcha: %w", attempt+1, err)
			continue
		}

		// 使用识别出的验证码进行登录
		resp, err := s.Login(username, password, captchaText, captchaResp.Data.UUID)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: failed to login: %w", attempt+1, err)
			continue
		}

		// 如果登录成功，直接返回结果
		if resp.Success {
			return resp, nil
		}

		// 如果登录失败但没有报错，记录失败原因
		lastErr = fmt.Errorf("attempt %d: login failed: %s", attempt+1, resp.ErrorMessage)
	}

	return nil, fmt.Errorf("login failed after %d attempts, last error: %w", maxRetries, lastErr)
}
