package mas

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"time"

	"github.com/kordar/gosms"
)

type normalReq struct {
	EcName    string `json:"ecName"`
	ApId      string `json:"apId"`
	Mobiles   string `json:"mobiles"`
	Content   string `json:"content"`
	Sign      string `json:"sign"`
	AddSerial string `json:"addSerial"`
	Mac       string `json:"mac"`
}

type masResp struct {
	RspCod   string `json:"rspcod"`
	MsgGroup string `json:"msgGroup"`
	Success  bool   `json:"success"`
}

func (p *Provider) SendSingle(req gosms.SMSRequest) ([]gosms.SMSResult, error) {
	if len(req.PhoneNumbers) == 0 {
		return nil, fmt.Errorf("mas: empty mobiles")
	}
	if len(req.PhoneNumbers) > 5000 {
		return nil, fmt.Errorf("mas: mobiles > 5000")
	}

	ecName := p.cfg.ExtraParams["ecName"]
	addSerial := req.ExtraParams["addSerial"]
	mobiles := strings.Join(req.PhoneNumbers, ",")

	mac := md5Lower(
		ecName +
			p.cfg.AccessKey +
			p.cfg.SecretKey +
			mobiles +
			req.Content +
			p.cfg.SignName +
			addSerial,
	)

	payload := normalReq{
		EcName:    ecName,
		ApId:      p.cfg.AccessKey,
		Mobiles:   mobiles,
		Content:   req.Content,
		Sign:      p.cfg.SignName,
		AddSerial: addSerial,
		Mac:       mac,
	}

	var resp masResp
	client := &http.Client{Timeout: 10 * time.Second}
	if err := postBase64JSON(client, p.endpoint, payload, &resp); err != nil {
		return nil, err
	}

	code := mapMASError(resp.RspCod)

	var results []gosms.SMSResult
	for _, phone := range req.PhoneNumbers {
		results = append(results, gosms.SMSResult{
			PhoneNumber: phone,
			Success:     code == gosms.ErrSuccess,
			Code:        string(code),
			Message:     resp.MsgGroup,
		})
	}

	if code != gosms.ErrSuccess {
		return results, &gosms.SMSError{
			Code:     code,
			Provider: "mas",
			RawCode:  resp.RspCod,
			Message:  "mas send failed",
		}
	}

	return results, nil
}

func (p *Provider) SendMultiple(reqs []gosms.SMSRequest) ([]gosms.SMSResult, error) {
	if len(reqs) == 0 {
		return nil, nil
	}

	contentMap := map[string]string{}
	for _, r := range reqs {
		for _, m := range r.PhoneNumbers {
			contentMap[m] = r.Content
		}
	}
	if len(contentMap) > 1000 {
		return nil, fmt.Errorf("mas: multi mobiles > 1000")
	}

	contentBytes, _ := json.Marshal(contentMap)
	content := string(contentBytes)

	ecName := p.cfg.ExtraParams["ecName"]
	addSerial := ""
	mac := md5Lower(
		ecName +
			p.cfg.AccessKey +
			p.cfg.SecretKey +
			content +
			p.cfg.SignName +
			addSerial,
	)

	payload := normalReq{
		EcName:    ecName,
		ApId:      p.cfg.AccessKey,
		Content:   content,
		Sign:      p.cfg.SignName,
		AddSerial: addSerial,
		Mac:       mac,
	}

	var resp masResp
	if err := postBase64JSON(p.client, p.endpoint, payload, &resp); err != nil {
		return nil, err
	}

	var results []gosms.SMSResult
	for phone := range contentMap {
		results = append(results, gosms.SMSResult{
			PhoneNumber: phone,
			Success:     resp.Success,
			Code:        resp.RspCod,
			Message:     resp.MsgGroup,
		})
	}

	return results, nil
}
