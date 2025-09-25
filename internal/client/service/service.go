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
	apiCleint *api.Client
	store     Store
	// recordsNew chan types.SecretData
	// recordsUpd chan types.SecretData
}

// NewService создает экземпляр объекта типа Service.
func NewService(apiCleint *api.Client, store Store) *Service {
	//return &Service{apiCleint: apiCleint, store: store, recordsNew: make(chan types.SecretData), recordsUpd: make(chan types.SecretData)}
	return &Service{apiCleint: apiCleint, store: store}
}

// Login логиниться.
func (s *Service) Login(ctx context.Context, login, password string) error {

	return s.apiCleint.Login(login, password)

}

// RegisterNewUser для регистрации нового пользователя.
func (s *Service) RegisterNewUser(ctx context.Context, login, password string) error {

	return s.apiCleint.RegisterNewUser(login, password)

}

// CreateNewSecret создает секрет в локальной БД.
func (s *Service) CreateNewSecret(ctx context.Context, secretData *types.SecretData) error {

	// Зашифруем данные

	// Поместим в БД

	//err := s.apiCleint.SendSecret(secretData)

	return nil

}

// EditSecret редактирует секрет в локальной БД.
func (s *Service) EditSecret(ctx context.Context, ID int, key, value, comment string) error {
	return nil
}

func (s *Service) GetSecretsToken(ctx context.Context) (string, error) {
	return s.store.GetSecretsToken(ctx)
}

func (s *Service) SaveSecretsToken(ctx context.Context, token string) error {
	return s.store.SaveSecretsToken(ctx, token)
}
