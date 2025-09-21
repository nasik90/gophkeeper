package service

import (
	"context"

	"github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/common/types"
)

type Store interface {
	SaveNewSecret(ctx context.Context, secretData types.SecretData) error
	UpdateSecret(ctx context.Context, secretData types.SecretData) error
	GetSecret(ctx context.Context, id int) error
	GetSecretsToken(ctx context.Context) (string, error)
	SaveSecretsToken(ctx context.Context, token string) error
}

// Service - структура, которая хранит ссылку на репозиторий, апи клиента и каналы для синхронизации с сервером.
type Service struct {
	apiCleint  *api.Client
	store      Store
	recordsNew chan types.SecretData
	recordsUpd chan types.SecretData
}

// NewService создает экземпляр объекта типа Service.
func NewService(apiCleint *api.Client, store Store) *Service {
	return &Service{apiCleint: apiCleint, store: store, recordsNew: make(chan types.SecretData), recordsUpd: make(chan types.SecretData)}
}

// Login логиниться.
func (s *Service) Login(ctx context.Context, login, password string) error {

	return s.apiCleint.Login(login, password)

}

// Login логиниться.
func (s *Service) RegisterNewUser(ctx context.Context, login, password string) error {

	return s.apiCleint.RegisterNewUser(login, password)

}

// CreateNewSecret создает секрет в локальной БД.
// После помещает в канал recordsNew для дальнейшей отправки на сервер.
// В самом начале данные будут зашифрованы.
func (s *Service) CreateNewSecret(ctx context.Context, key, value, comment string) error {

	// // URL для POST-запроса
	// url := "http://localhost:8080/api/user/loadSecret"

	// var requestData types.SecretData

	// requestData.Key = key
	// requestData.Value = value
	// requestData.Comment = comment

	// // Преобразуем данные в JSON
	// jsonData, err := json.Marshal(requestData)
	// if err != nil {
	// 	return err
	// }

	// // Создаем новый POST-запрос
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	return err
	// }

	// req.Header.Set("Content-Type", "application/json")
	// response, err := s.apiCleint.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer response.Body.Close()

	// if response.StatusCode != http.StatusOK {
	// 	fmt.Println("Ошибка: статус-код", response.StatusCode)
	// 	return err
	// }

	// var responseBody struct {
	// 	RecordID int `json:"recordID"`
	// }
	// if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
	// 	fmt.Println("Ошибка при чтении ответа:", err)
	// 	return err
	// }

	// fmt.Println(responseBody)

	return nil

}

// EditSecret редактирует секрет в локальной БД.
// После помещает в канал recordsUpd для дальнейшей отправки на сервер.
// В самом начале данные будут зашифрованы.
func (s *Service) EditSecret(ctx context.Context, ID int, key, value, comment string) error {
	return nil
}

func (s *Service) GetSecretsToken(ctx context.Context) (string, error) {
	return s.store.GetSecretsToken(ctx)
}

func (s *Service) SaveSecretsToken(ctx context.Context, token string) error {
	return s.store.SaveSecretsToken(ctx, token)
}
