package mas

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"time"

	"github.com/kordar/gosms"
)

type templateReq struct {
	// SecretKey  string `json:"secretKey"`
	Sign       string `json:"sign"`
	ApId       string `json:"apId"`
	Mac        string `json:"mac"`
	EcName     string `json:"ecName"`
	Params     string `json:"params"`
	TemplateId string `json:"templateId"`
	AddSerial  string `json:"addSerial"`
	Mobiles    string `json:"mobiles"`
}

func (p *Provider) SendTemplate(req gosms.SMSRequest) ([]gosms.SMSResult, error) {
	if req.TemplateID == "" {
		return nil, fmt.Errorf("mas: missing templateId")
	}

	params := []string{}
	if len(req.TemplateParams) > 0 {
		params = req.TemplateParams
	} else {
		// 兼容旧的 TemplateVars 方式，但顺序不保证
		for _, v := range req.TemplateVars {
			params = append(params, v)
		}
	}
	paramJSON, _ := json.Marshal(params)

	mobiles := strings.Join(req.PhoneNumbers, ",")
	ecName := p.cfg.ExtraParams["ecName"]
	addSerial := req.ExtraParams["addSerial"]

	// ecName = "\u653f\u4f01\u5206\u516c\u53f8\u6d4b\u8bd5"
	// p.cfg.SecretKey = "123qwe"
	// p.cfg.SignName = "4sEuJxDpC"
	// p.cfg.AccessKey = "demo0"
	// mac := "02009a533ee3fb8603062971a53beff0"
	// paramJSONStr := "[\"abcde\"]"
	// req.TemplateID = "38516fabae004eddbfa3ace1d4194696"
	// addSerial = ""
	// mobiles = "13800138000"

	mac := md5Lower(
		ecName +
			p.cfg.AccessKey +
			p.cfg.SecretKey +
			req.TemplateID +
			mobiles +
			string(paramJSON) +
			p.cfg.SignName +
			addSerial,
	)

	payload := templateReq{
		EcName:     ecName,
		ApId:       p.cfg.AccessKey,
		TemplateId: req.TemplateID,
		Mobiles:    mobiles,
		Params:     string(paramJSON),
		Sign:       p.cfg.SignName,
		AddSerial:  addSerial,
		Mac:        mac,
		// SecretKey:  p.cfg.SecretKey,
	}

	var resp masResp
	client := &http.Client{Timeout: 10 * time.Second}
	if err := postBase64JSON(client, p.endpoint, payload, &resp); err != nil {
		return nil, err
	}

	var results []gosms.SMSResult
	for _, phone := range req.PhoneNumbers {
		results = append(results, gosms.SMSResult{
			PhoneNumber: phone,
			Success:     resp.Success,
			Code:        resp.RspCod,
			Message:     resp.MsgGroup,
		})
	}

	return results, nil
}
