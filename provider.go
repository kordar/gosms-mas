package mas

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/kordar/gosms"
)

type Provider struct {
	cfg      *gosms.SMSConfig
	endpoint string
	client   *http.Client
}

func New(cfg *gosms.SMSConfig) (gosms.SMSProvider, error) {
	endpoint := cfg.ExtraParams["endpoint"]
	if endpoint == "" {
		endpoint = "https://112.33.46.17:37892/sms/tmpsubmit"
	}

	return &Provider{
		cfg:      cfg,
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}, nil
}

func init() {
	gosms.RegisterProvider("mas", New)
}
