# DMXSmart 客户端

一个用于与DMXSmart仓储管理系统交互的Go客户端库，提供自动化的认证、订单管理和OCR验证码识别功能。

## 功能特性

-   🔐 **自动认证**: 支持账号密码登录，自动处理RSA加密和验证码识别
-   🤖 **OCR集成**: 集成验证码自动识别服务，实现无人值守登录
-   📦 **订单管理**: 支持拣货波次订单查询和管理
-   📝 **日志记录**: 完整的日志记录和调试支持
-   ⚙️ **配置管理**: 灵活的YAML配置文件支持

## 项目结构

```
dmxsmart-client/
├── client.go               # 主客户端实现
├── client_test.go          # 测试文件
├── config.yaml             # 配置文件
├── config.yaml.example     # 配置文件示例
├── go.mod                  # Go模块文件
├── config/                 # 配置管理
├── service/                # 服务层
│   ├── auth.go             # 认证服务
│   ├── pickup.wave.go      # 拣货波次订单服务
│   ├── client.go           # HTTP客户端
│   ├── manager.go          # 服务管理器
│   └── encrypt.go          # 加密服务
├── model/                  # 数据模型
│   ├── auth.go             # 认证相关模型
│   ├── pickup.wave.go      # 拣货订单相关模型
│   └── response.go         # 响应模型
├── ocr/                    # OCR功能
│   └── captcha.ocr.go      # 验证码识别
├── logger/                 # 日志管理
└── logs/                   # 日志文件目录
```

## 安装

确保你已经安装了Go 1.23.3或更高版本。

```bash
go mod tidy
```

## 配置

1. 复制配置文件示例：

```bash
cp config.yaml.example config.yaml
```

2. 编辑 `config.yaml` 文件：

```yaml
account: your_account # DMXSmart账号
password: your_password # DMXSmart密码
access_token: your_token # 访问令牌（可选）
warehouse_id: 745 # 仓库ID
customer_ids: # 客户ID列表
    - 37046
    - 4040
ocr_endpoint: https://ddddocr.xxxx.com/ocr_base64 # OCR服务端点
timeout: 30 # 请求超时时间（秒）
debug: false # 是否开启调试模式
```

## 使用方法

### 基本使用

```go
package main

import (
    "fmt"
    "path/filepath"

    c "github.com/bestk/dmxsmart-client/client"
)

func main() {
    // 创建客户端
    configPath := filepath.Join(".", "config.yaml")
    client, err := c.NewDMXSmartClient(configPath)
    if err != nil {
        panic(err)
    }

    // 执行自动登录（包含OCR验证码识别）
    resp, err := client.services.Auth.LoginWithAutoOCR(
        client.config.Account,
        client.config.Password,
    )
    if err != nil {
        panic(err)
    }

    if resp.Success {
        fmt.Println("登录成功!")
    } else {
        fmt.Printf("登录失败: %s\n", resp.ErrorMessage)
    }
}
```

### 订单管理

```go
// 获取待拣货订单列表
orders, err := client.services.PickupWave.GetWaitingPickOrders(
    1,    // 页码
    10,   // 每页大小
    client.config.CustomerIDs, // 客户ID列表
)
if err != nil {
    panic(err)
}

fmt.Printf("找到 %d 个待拣货订单\n", len(orders.Data.Records))
```

### 会话验证

```go
// 验证当前会话是否有效
err := client.services.Auth.ValidateSession()
if err != nil {
    fmt.Println("会话已失效，需要重新登录")
} else {
    fmt.Println("会话有效")
}
```

## 测试

运行测试：

```bash
go test -v
```

运行特定测试：

```bash
# 测试自动登录功能
go test -run TestLoginWithAutoOCR -v

# 测试订单查询功能
go test -run TestGetWaitingPickOrders -v
```

## 依赖项

-   [resty](https://github.com/go-resty/resty) - HTTP客户端
-   [slog](https://github.com/gookit/slog) - 日志库
-   [yaml.v3](https://gopkg.in/yaml.v3) - YAML解析

## 核心功能说明

### 自动认证流程

1. 获取验证码图片
2. 使用OCR服务自动识别验证码
3. RSA加密密码
4. 提交登录请求
5. 保存认证令牌

### OCR验证码识别

项目集成了外部OCR服务来自动识别登录验证码，支持：

-   Base64图片编码
-   自动重试机制
-   可配置的OCR服务端点

### 安全特性

-   使用RSA公钥加密敏感信息
-   支持访问令牌认证
-   安全的会话管理

## 调试

启用调试模式：

```yaml
debug: true
```

调试模式下会输出详细的HTTP请求和响应信息。

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

本项目采用MIT许可证。详见LICENSE文件。

## 联系方式

如有问题或建议，请通过GitHub Issues联系我们。
