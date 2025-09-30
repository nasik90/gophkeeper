package api

import (
	"encoding/json"
	"time"

	"github.com/nasik90/gophkeeper/internal/common/types"
)

// Client - структура для хранения базового пути и встроенной структуры http.Client.
type Client struct {
	baseURL string
}

// NewClient создает структуру типа Client.
func NewClient(baseURL string) *Client {
	c := &Client{}
	c.baseURL = baseURL
	return c
}

// SendSecret выполянет http вызов для отправки секрета.
func (c *Client) SendSecret(secretData *types.SecretData) error {

	method := "/api/secrets/loadSecret"
	_, err := postJSON(c.baseURL, method, secretData)

	return err
}

// uploadSecrets выполянет http вызов для получения измененных секретов.
func (c *Client) UploadSecrets(fromDate time.Time) (*[]types.SecretData, error) {

	method := "/api/secrets/getSecrets"

	var reqData struct {
		FromDate time.Time `json:"fromDate"`
	}
	reqData.FromDate = fromDate

	respBody, err := postJSON(c.baseURL, method, reqData)
	if err != nil {
		return nil, err
	}
	var responseBody []types.SecretData
	if err := json.Unmarshal(respBody, &responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil

}
