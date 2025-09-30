package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nasik90/gophkeeper/internal/client/api"
	"github.com/nasik90/gophkeeper/internal/client/crypto"
	"github.com/nasik90/gophkeeper/internal/common/types"
)

type Store interface {
	SaveNewSecret(ctx context.Context, secretData *types.SecretData) error
	UpdateSecret(ctx context.Context, secretData *types.SecretData) error
	InsertUpdateSecret(ctx context.Context, secretData *types.SecretData) error
	GetSecret(ctx context.Context, id int) error
	GetSecrets(ctx context.Context, toSend bool) (*[]types.SecretData, error)
	GetSecretsToken(ctx context.Context) (string, error)
	SaveSecretsToken(ctx context.Context, token string) error
	SaveDataVersion(ctx context.Context, dataVersion time.Time) error
	GetDataVersion(ctx context.Context) (time.Time, error)
}

// Service - структура, которая хранит ссылку на репозиторий, апи клиента и каналы для синхронизации с сервером.
type Service struct {
	apiCleint *api.Client
	store     Store
	key       []byte
}

// NewService создает экземпляр объекта типа Service.
func NewService(apiCleint *api.Client, store Store, masterPassword string) *Service {
	//return &Service{apiCleint: apiCleint, store: store, recordsNew: make(chan types.SecretData), recordsUpd: make(chan types.SecretData)}
	key := crypto.GenerateKey(masterPassword)
	return &Service{apiCleint: apiCleint, store: store, key: key}
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
// Для упрощения бинарные данные будем хранить в БД.
func (s *Service) CreateNewSecret(ctx context.Context, secretData *types.SecretData) error {

	newUUID := uuid.New()
	secretData.Guid = newUUID.String()

	// Зашифруем чувствительные данные
	err := s.encryptSensitiveData(secretData)
	if err != nil {
		return err
	}

	// Поместим в БД
	secretData.ToSend = true
	err = s.store.SaveNewSecret(ctx, secretData)
	if err != nil {
		return err
	}

	return nil

}

// EditSecret редактирует секрет в локальной БД.
func (s *Service) EditSecret(ctx context.Context, secretData *types.SecretData) error {
	err := s.encryptSensitiveData(secretData)
	if err != nil {
		return err
	}
	secretData.ToSend = true
	return s.store.UpdateSecret(ctx, secretData)
}

func (s *Service) encryptSensitiveData(secretData *types.SecretData) error {
	encryptedKey, err := crypto.Encrypt([]byte(secretData.Key), s.key)
	if err != nil {
		return err
	}
	secretData.Key = encryptedKey

	encryptedValue, err := crypto.Encrypt([]byte(secretData.Value), s.key)
	if err != nil {
		return err
	}
	secretData.Value = encryptedValue
	return nil
}

// GetSecrets получает секреты из локальной БД и расшифровывает.
func (s *Service) GetSecrets(ctx context.Context) (*[]types.SecretData, error) {
	toSend := false
	secrets, err := s.store.GetSecrets(ctx, toSend)
	// Расшифруем данные
	for i, secretData := range *secrets {
		decryptedKey, err := crypto.Decrypt([]byte(secretData.Key), s.key)
		if err != nil {
			return nil, err
		}
		(*secrets)[i].Key = decryptedKey
		decryptedValue, err := crypto.Decrypt([]byte(secretData.Value), s.key)
		if err != nil {
			return nil, err
		}
		(*secrets)[i].Value = decryptedValue
	}
	return secrets, err
}

func (s *Service) GetSecretsToken(ctx context.Context) (string, error) {
	return s.store.GetSecretsToken(ctx)
}

func (s *Service) SaveSecretsToken(ctx context.Context, token string) error {
	return s.store.SaveSecretsToken(ctx, token)
}
