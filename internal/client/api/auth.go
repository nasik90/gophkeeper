package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"

	"github.com/nasik90/gophkeeper/internal/common/constants"
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
		return errors.New("Ошибка при преобразовании в JSON:" + err.Error())
	}

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("Ошибка при создании запроса:" + err.Error())
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	response, err := c.Do(req)
	if err != nil {
		return errors.New("Ошибка при выполнении запроса:" + err.Error())
	}
	defer response.Body.Close()

	// Проверяем статус-код ответа
	if response.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == constants.CookieName {
			saveToken(cookie.Value)
			break
		}
	}

	return nil

}

func (c *Client) RegisterNewUser(login, password string) error {

	// URL для POST-запроса
	method := "/api/user/register"
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
		return errors.New("Ошибка при преобразовании в JSON:" + err.Error())
	}

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("Ошибка при создании запроса:" + err.Error())
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	response, err := c.Do(req)
	if err != nil {
		return errors.New("Ошибка при выполнении запроса:" + err.Error())
	}
	defer response.Body.Close()

	// Проверяем статус-код ответа
	if response.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	return nil

}

func saveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".gophkeeper")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	// Попытка прочитать существующий конфиг (если есть)
	_ = viper.ReadInConfig()

	viper.Set("jwt_token", token)

	return viper.WriteConfigAs(filepath.Join(configDir, "config.json"))
}

func getToken() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configDir := filepath.Join(homeDir, ".gophkeeper")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return ""
	}

	return viper.GetString("jwt_token")
}

func (c *Client) setToken() error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	tokenSet := false
	cookies := c.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == constants.CookieName {
			tokenSet = true
			break
		}
	}
	if !tokenSet {
		token := getToken()
		if token != "" {
			cookie := http.Cookie{Name: constants.CookieName, Value: token}
			cookies = append(cookies, &cookie)
			c.Jar.SetCookies(u, cookies)
		}
	}
	return nil
}
