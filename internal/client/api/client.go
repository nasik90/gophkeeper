package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/common/types"
	"go.uber.org/zap"
)

// Client - структура для хранения базового пути и встроенной структуры http.Client.
type Client struct {
	*http.Client
	baseURL string
}

// NewClient создает структуру типа Client.
func NewClient(baseURL string) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		logger.Log.Fatal("initialize cookie", zap.Error(err))
	}
	c := &Client{}
	c.Client = &http.Client{Jar: jar}
	c.baseURL = baseURL
	return c
}

// SendSecret выполянет http вызов для отправки секрета.
func (c *Client) SendSecret(secretData *types.SecretData) error {

	c.setToken()

	// URL для POST-запроса
	method := "/api/secrets/loadSecret"
	url := c.baseURL + method

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(secretData)
	if err != nil {
		return err
	}

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := c.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New("Ошибка: статус-код " + strconv.Itoa(response.StatusCode))
		return err
	}

	return nil
}

// uploadSecrets выполянет http вызов для получения измененных секретов.
func (c *Client) UploadSecrets(fromDate time.Time) (*[]types.SecretData, error) {

	c.setToken()

	// URL для POST-запроса
	method := "/api/secrets/getSecrets"
	url := c.baseURL + method

	var responseData struct {
		FromDate time.Time `json:"fromDate"`
	}
	responseData.FromDate = fromDate
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return nil, err
	}

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New("Ошибка: статус-код " + strconv.Itoa(response.StatusCode))
		return nil, err
	}

	var responseBody []types.SecretData
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil

}
