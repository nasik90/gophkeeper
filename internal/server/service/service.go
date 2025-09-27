package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	types "github.com/nasik90/gophkeeper/internal/common/types"
)

type Store interface {
	SaveNewUser(ctx context.Context, user, password string) error
	UserIsValid(ctx context.Context, login, password string) (bool, error)
	GetUserID(ctx context.Context, login string) (int, error)

	SaveNewSecret(ctx context.Context, secretData *types.SecretData, userID int, creationDate time.Time) (int, error)
	UpdateSecret(ctx context.Context, userID int, SecretData *types.SecretData, updatingDate time.Time) error
	GetUserSecretList(ctx context.Context, userID int, secretID int, fromDate time.Time) (*[]types.SecretData, error)
	GetUserSecretsVersion(ctx context.Context, userID int) (int, error)
	GetSecretVersion(ctx context.Context, userID int, secretID int) (int, error)
	UpdateUserSecretsVersion(ctx context.Context, userID int, newVersion int, updatingDate time.Time) error

	Close() error
}

type Service struct {
	store Store
}

// NewService создает экземпляр объекта типа Service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) RegisterNewUser(ctx context.Context, login, password string) error {
	return s.store.SaveNewUser(ctx, login, password)
}

func (s *Service) UserIsValid(ctx context.Context, login, password string) (bool, error) {
	return s.store.UserIsValid(ctx, login, password)
}

func (s *Service) LoadSecret(ctx context.Context, SecretData *types.SecretData, login string) (int, error) {
	var recordID int
	userID, err := s.store.GetUserID(ctx, login)
	if err != nil {
		return 0, err
	}
	err = s.store.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		creationDate := time.Now()
		// Сохраним новый секрет
		recordID, err = s.store.SaveNewSecret(ctx, SecretData, userID, creationDate)
		if err != nil {
			return err
		}
		// TODO: сохранение двоичных данных при необходимости
		// Получим текущую версию данных пользователя
		userSecretsVersion, err := s.store.GetUserSecretsVersion(ctx, userID)
		if err != sql.ErrNoRows && err != nil {
			return err
		}
		// Поднимем версию и сохраним ее
		userSecretsVersion++
		err = s.store.UpdateUserSecretsVersion(ctx, userID, userSecretsVersion, creationDate)
		return err
	})

	return recordID, err

}

func (s *Service) UpdateSecret(ctx context.Context, SecretData *types.SecretData, login string) error {
	if SecretData.Id == 0 {
		return errors.New("SecredData ID can not be empty")
	}
	userID, err := s.store.GetUserID(ctx, login)
	if err != nil {
		return err
	}
	err = s.store.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		updatingDate := time.Now()
		// Сравним версии, если не равны, то вернем ошибку
		versionID, err := s.store.GetSecretVersion(ctx, userID, SecretData.Id)
		if err != nil {
			return err
		}
		if versionID != SecretData.VersionID {
			return errors.New("version ID doesn`t match. data was changed or deleted")
		}
		SecretData.VersionID++
		// Сохраним новый секрет
		err = s.store.UpdateSecret(ctx, userID, SecretData, updatingDate)
		if err != nil {
			return err
		}
		// TODO: сохранение двоичных данных
		// Получим текущую версию данных пользователя
		userSecretsVersion, err := s.store.GetUserSecretsVersion(ctx, userID)
		if err != sql.ErrNoRows && err != nil {
			return err
		}
		// Поднимем версию и сохраним ее
		userSecretsVersion++
		err = s.store.UpdateUserSecretsVersion(ctx, userID, userSecretsVersion, updatingDate)
		return err
	})

	return err

}

func (s *Service) GetSecret(ctx context.Context, id int, login string) (types.SecretData, error) {
	return types.SecretData{}, nil
}

func (s *Service) GetSecrets(ctx context.Context, login string, fromVersionID int) ([]types.SecretData, error) {
	return []types.SecretData{}, nil
}
