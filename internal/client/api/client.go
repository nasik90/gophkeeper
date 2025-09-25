package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	c.Client = &http.Client{Jar: jar}
	c.baseURL = baseURL
	//c.Jar = jar
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
		fmt.Println("Ошибка: статус-код", response.StatusCode)
		return err
	}

	var responseBody struct {
		RecordID int `json:"recordID"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return err
	}

	fmt.Println(responseBody)
	return nil
}

// GetSecrets выполянет http вызов для получения секрета.
func (c *Client) GetSecret(id int, user string) ([]types.SecretData, error) {
	return nil, nil
}

// GetSecrets выполянет http вызов для получения секретов.
func (c *Client) GetSecrets(user string) ([]types.SecretData, error) {
	return nil, nil
}

// uploadSecrets выполянет http вызов для получения измененных секретов.
func (c *Client) UploadSecrets(fromVesrion int, user string) ([]types.SecretData, error) {
	return nil, nil
}
