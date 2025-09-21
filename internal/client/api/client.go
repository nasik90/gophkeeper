package api

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/common/types"
	"go.uber.org/zap"
)

// Client - структура для хранения базового пути и встроенной структуры http.Client.
type Client struct {
	*http.Client
	baseURL string
	//jar     http.CookieJar
}

// NewClient создает структуру типа Client.
func NewClient(baseURL string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Log.Fatal("initialize cookie", zap.Error(err))
	}
	c := &Client{}
	c.Client = &http.Client{}
	c.baseURL = baseURL
	c.Jar = jar
	return c
}

// SendSecret выполянет http вызов для отправки секрета.
func (c *Client) SendSecret(secretData types.SecretData) error {
	return nil
}

// uploadSecrets выполянет http вызов для получения измененных секретов.
func (c *Client) UploadSecrets() ([]types.SecretData, error) {
	return nil, nil
}
