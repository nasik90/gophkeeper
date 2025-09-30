package api

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/nasik90/gophkeeper/internal/common/constants"
)

func (c *Client) Login(login, password string) error {

	method := "/api/user/login"
	return c.registerLogin(method, login, password)

}

func (c *Client) RegisterNewUser(login, password string) error {

	method := "/api/user/register"
	return c.registerLogin(method, login, password)

}

func (c *Client) registerLogin(method, login, password string) error {
	// Данные для запроса
	var requestData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	requestData.Login = login
	requestData.Password = password

	_, err := postJSON(c.baseURL, method, requestData)
	if err != nil {
		return err
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

func setToken(c *http.Client, baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	c.Jar, err = cookiejar.New(nil)
	if err != nil {
		return err
	}
	cookies := c.Jar.Cookies(u)
	token := getToken()
	if token != "" {
		cookie := http.Cookie{Name: constants.CookieName, Value: token}
		cookies = append(cookies, &cookie)
		c.Jar.SetCookies(u, cookies)
	}
	return nil
}
