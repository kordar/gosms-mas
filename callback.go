package mas

import (
	"encoding/json"
	"github.com/kordar/gosms"
	"time"
)

type reportReq struct {
	ReportStatus string `json:"reportStatus"`
	Mobile       string `json:"mobile"`
	SubmitDate   string `json:"submitDate"`
	ReceiveDate  string `json:"receiveDate"`
	ErrorCode    string `json:"errorCode"`
	MsgGroup     string `json:"msgGroup"`
}

func (p *Provider) HandleReport(body []byte) ([]gosms.SMSReport, error) {
	var r reportReq
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	status := "FAILED"
	if r.ReportStatus == "DELIVRD" {
		status = "DELIVERED"
	}

	t, _ := time.Parse("20060102150405", r.ReceiveDate)

	return []gosms.SMSReport{{
		PhoneNumber: r.Mobile,
		Status:      status,
		MsgID:       r.MsgGroup,
		Timestamp:   t,
	}}, nil
}

type inboundReq struct {
	Mobile     string `json:"mobile"`
	SmsContent string `json:"smsContent"`
	SendTime   string `json:"sendTime"`
	AddSerial  string `json:"addSerial"`
}

func (p *Provider) HandleInbound(body []byte) ([]gosms.SMSInbound, error) {
	var r inboundReq
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	t, _ := time.Parse("2006-01-02 15:04:05", r.SendTime)

	return []gosms.SMSInbound{{
		PhoneNumber: r.Mobile,
		Content:     r.SmsContent,
		MsgID:       "",
		Timestamp:   t,
	}}, nil
}
