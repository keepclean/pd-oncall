package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string
	Token     string

	httpClient *http.Client
}

func NewPDApiClient(u *url.URL, v, t string, timeout time.Duration) *Client {
	return &Client{
		BaseURL:   u,
		UserAgent: fmt.Sprintf("pd-oncall-%s", v),
		Token:     t,

		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}
