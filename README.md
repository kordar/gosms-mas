# gosms-mas

`gosms-mas` 是 [gosms](https://github.com/kordar/gosms) 的移动云短信（MAS）提供商实现。

## 安装

```bash
go get github.com/kordar/gosms-mas
```

## 配置

在使用之前，需要初始化 `gosms` 并注册 `mas` 提供商。

### 配置参数说明

| 参数 | 说明 | 备注 |
| --- | --- | --- |
| `AccessKey` | 用户账号 | 对应 MAS 平台的 `apId` |
| `SecretKey` | 用户密码 | 对应 MAS 平台的 `secretKey` |
| `SignName` | 签名编码 | 对应 MAS 平台的 `sign` |
| `ExtraParams["endpoint"]` | 接口地址 | 默认为 `https://112.33.46.17:37892/sms/tmpsubmit` |
| `ExtraParams["ecName"]` | 企业名称 | 对应 MAS 平台的 `ecName` |

## 使用示例

### 1. 初始化

```go
package main

import (
	"fmt"
	"github.com/kordar/gosms"
	_ "github.com/kordar/gosms-mas" // 自动注册 mas 提供商
)

func main() {
	cfg := &gosms.SMSConfig{
		AccessKey: "YOUR_AP_ID",
		SecretKey: "YOUR_SECRET_KEY",
		SignName:  "YOUR_SIGN_ID",
		ExtraParams: map[string]string{
			"ecName":   "YOUR_EC_NAME",
			"endpoint": "https://112.33.46.17:37892/sms/tmpsubmit", // 可选
		},
	}
    
    // 初始化名为 "mas-client" 的实例，使用 "mas" 提供商
	if err := gosms.New("mas-client", "mas", cfg); err != nil {
		panic(err)
	}
}
```

### 2. 发送普通短信

```go
func SendNormal() {
    client := gosms.Get("mas-client")
    if client == nil {
        panic("client not found")
    }

    req := gosms.SMSRequest{
        PhoneNumbers: []string{"13800138000"},
        Content:      "您的验证码是 123456",
        ExtraParams: map[string]string{
            "addSerial": "", // 可选，扩展码
        },
    }

    results, err := client.SendSingle(req)
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
        return
    }

    for _, res := range results {
        fmt.Printf("发送结果:Success=%v, Message=%s\n", res.Success, res.Message)
    }
}
```

### 3. 发送模板短信

```go
func SendTemplate() {
    client := gosms.Get("mas-client")
    
    req := gosms.SMSRequest{
        PhoneNumbers: []string{"13800138000"},
        TemplateID:   "TEMPLATE_ID", // 模板ID
        TemplateVars: map[string]string{
            "param1": "value1", // 模板参数，按顺序填充
            "param2": "value2",
        },
    }

    results, err := client.SendTemplate(req)
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
        return
    }
    
    for _, res := range results {
        fmt.Printf("发送结果:Success=%v, Message=%s\n", res.Success, res.Message)
    }
}
```

### 4. 批量发送

```go
func SendBatch() {
    client := gosms.Get("mas-client")
    
    reqs := []gosms.SMSRequest{
        {
            PhoneNumbers: []string{"13800138000"},
            Content:      "内容1",
        },
        {
            PhoneNumbers: []string{"13900139000"},
            Content:      "内容2",
        },
    }

    results, err := client.SendMultiple(reqs)
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
        return
    }
    
    // MAS 批量发送通常返回一个整体结果
    for _, res := range results {
        fmt.Printf("批量发送结果:Success=%v, Message=%s\n", res.Success, res.Message)
    }
}
```
