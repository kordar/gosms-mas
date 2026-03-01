package mas

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unicode/utf8"
)

func md5Lower(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// jsonMarshalUnicode 进行 JSON 序列化，并将所有非 ASCII 字符转换为 \uXXXX
func jsonMarshalUnicode(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for i := 0; i < len(raw); {
		r, size := utf8.DecodeRune(raw[i:])
		if r < 128 {
			buf.WriteRune(r)
		} else {
			buf.WriteString(fmt.Sprintf("\\u%04x", r))
		}
		i += size
	}
	return buf.Bytes(), nil
}

// postBase64JSON 序列化 payload 为 JSON (Unicode 转义)，然后进行 Base64 编码，作为请求体发送
// Content-Type 设置为 text/plain，避免服务端尝试直接解析 JSON
func postBase64JSON(client *http.Client, url string, payload any, resp any) error {
	// 使用自定义的 Unicode 序列化
	// raw, err := jsonMarshalUnicode(payload)
	// if err != nil {
	// return err
	// }

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	fmt.Println("Unicode JSON:")
	fmt.Println(string(raw))

	// raw = []byte(`{"secretKey": "123qwe", "sign": "4sEuJxDpC", "apId": "demo0", "mac": "02009a533ee3fb8603062971a53beff0", "ecName": "\u653f\u4f01\u5206\u516c\u53f8\u6d4b\u8bd5", "params": "[\"abcde\"]", "templateId": "38516fabae004eddbfa3ace1d4194696", "addSerial": "", "mobiles": "13800138000"}`)

	// 转 base64
	b64 := base64.StdEncoding.EncodeToString(raw)
	fmt.Println("\nBase64:")
	fmt.Println(b64)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(b64))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %v", err)
	}

	// Debug: 打印原始响应 Body
	fmt.Printf("[MAS Debug] Response Body: %s\n", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, resp); err != nil {
		preview := string(bodyBytes)
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		return fmt.Errorf("decode response failed: %v, status: %s, body: %s", err, r.Status, preview)
	}

	return nil
}
