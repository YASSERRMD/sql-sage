package llm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrBlockedHost  = errors.New("provider host not allowed")
	ErrInvalidURL   = errors.New("invalid provider url")
	ErrUnauthorized = errors.New("provider unauthorized")
)

type Client struct {
	hc *http.Client
}

func NewClient() *Client {
	return &Client{hc: &http.Client{Timeout: 30 * time.Second}}
}

type TestRequest struct {
	BaseURL       string
	APIKey        string
	Model         string
	HostAllowlist []string
}

type TestResult struct {
	OK        bool
	Message   string
	LatencyMs int64
}

func (c *Client) ValidateURL(raw string, allowlist []string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	host := strings.ToLower(u.Hostname())
	if len(allowlist) == 0 {
		return nil
	}
	for _, h := range allowlist {
		if strings.EqualFold(h, host) {
			return nil
		}
	}
	return ErrBlockedHost
}

func (c *Client) TestConnection(ctx context.Context, req TestRequest) (*TestResult, error) {
	if err := c.ValidateURL(req.BaseURL, req.HostAllowlist); err != nil {
		return &TestResult{OK: false, Message: err.Error()}, err
	}
	endpoint := strings.TrimRight(req.BaseURL, "/") + "/models"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return &TestResult{OK: false, Message: "build request"}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
	start := time.Now()
	resp, err := c.hc.Do(httpReq)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &TestResult{OK: false, Message: err.Error(), LatencyMs: latency}, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return &TestResult{OK: true, Message: "ok", LatencyMs: latency}, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return &TestResult{OK: false, Message: "unauthorized", LatencyMs: latency}, ErrUnauthorized
	default:
		return &TestResult{OK: false, Message: fmt.Sprintf("status %d", resp.StatusCode), LatencyMs: latency}, nil
	}
}
