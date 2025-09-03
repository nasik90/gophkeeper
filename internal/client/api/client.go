package api

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/nasik90/gophkeeper/internal/common/logger"
	"go.uber.org/zap"
)

type Client struct {
	*http.Client
	baseURL string
	jar     http.CookieJar
	// возможно, cookieJar или токены для аутентификации
}

func NewClient(baseURL string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Log.Fatal("initialize cookie", zap.Error(err))
	}
	return &Client{baseURL: baseURL, jar: jar}
}
