package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Login(login, password string) error {

	// URL для POST-запроса
	method := "/api/user/login"
	url := c.baseURL + method

	// Данные для запроса
	var requestData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	requestData.Login = login
	requestData.Password = password

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return err
	}

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return err
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return err
	}
	defer response.Body.Close() // Закрываем тело ответа после завершения работы с ним

	// Проверяем статус-код ответа
	if response.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: статус-код", response.StatusCode)
		return err
	}

	return nil

}
