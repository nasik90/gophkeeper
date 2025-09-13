package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nasik90/gophkeeper/internal/client/api"
)

type Service struct {
	apiCleint *api.Client
}

func NewService(apiCleint *api.Client) *Service {
	return &Service{apiCleint: apiCleint}
}

func (s *Service) Login(ctx context.Context, login, password string) error {

	return s.apiCleint.Login(login, password)

}

func (s *Service) CreateNewSecret(ctx context.Context, key, value, comment string) error {

	// URL для POST-запроса
	url := "http://localhost:8080/api/user/loadSecret"

	// Данные для запроса
	var requestData struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Comment string `json:"comment"`
	}

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

	// // Добавляем cookie
	// cookie := &http.Cookie{
	// 	Name:  "gophkeeper",
	// 	Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjU4MjQ0NTksIlVzZXJJRCI6Im5hc2lrOTAifQ.hpWOJIvMcba8pU2a1gVTV5I1v4k5LyoPG8girT5ih38",
	// }
	// req.AddCookie(cookie)

	// // Выполняем запрос
	// client := &http.Client{}
	// response, err := client.Do(req)
	response, err := s.apiCleint.Do(req)
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

	// Читаем и выводим ответ
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
