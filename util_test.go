package mas

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(md5Lower("政企分公司测试demo0123qwe13800138000移动改变生活。DWItALe3A"))
}

type testReq struct {
	EcName    string `json:"ecName"`
	ApId      string `json:"apId"`
	Mobiles   string `json:"mobiles"`
	Content   string `json:"content"`
	Sign      string `json:"sign"`
	AddSerial string `json:"addSerial"`
	Mac       string `json:"mac"`
}

func TestBase64Encoding(t *testing.T) {
	// 构造测试数据
	req := testReq{
		EcName:    "政企分公司测试",
		ApId:      "demo0",
		Mobiles:   "13800138000",
		Content:   "移动改变生活。",
		Sign:      "DWItALe3A",
		AddSerial: "",
		Mac:       "7997ddb079db2155b517b21b2a812370",
	}

	// 1. JSON 序列化 (使用新的 Unicode 转义逻辑)
	jsonBytes, err := jsonMarshalUnicode(req)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	fmt.Printf("JSON Unicode Escaped: %s\n", string(jsonBytes))

	// 2. Base64 编码
	b64 := base64.StdEncoding.EncodeToString(jsonBytes)
	fmt.Printf("Base64 Encoded: %s\n", b64)

	// 3. 验证是否包含中文 (不应包含)
	for _, b := range jsonBytes {
		if b > 127 {
			t.Errorf("Found non-ASCII character in JSON: %v", b)
		}
	}
}
