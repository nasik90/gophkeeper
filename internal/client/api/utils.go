package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nasik90/gophkeeper/internal/common/constants"
)

// postJSON выполняет HTTP POST запрос с JSON телом и возвращает тело ответа как []byte.
func postJSON(baseURL, method string, data interface{}) ([]byte, error) {
	url := baseURL + method
	// Кодируем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	setToken(client, baseURL)

	// Создаем POST запрос с телом
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform POST request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем код ответа
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Читаем тело ответа
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == constants.CookieName {
			saveToken(cookie.Value)
			break
		}
	}

	return respBody, nil
}
