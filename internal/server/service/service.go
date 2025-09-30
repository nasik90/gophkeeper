package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Thiht/transactor"
	types "github.com/nasik90/gophkeeper/internal/common/types"
)

type Store interface {
	SaveNewUser(ctx context.Context, user, password string) error
	UserIsValid(ctx context.Context, login, password string) (bool, error)
	GetUserID(ctx context.Context, login string) (int, error)

	LoadSecret(ctx context.Context, secretData *types.SecretData, userID int) error
	GetUserSecretList(ctx context.Context, userID int, fromDate time.Time) (*[]types.SecretData, error)
	GetUserSecretsVersion(ctx context.Context, userID int) (int, error)
	GetSecretVersion(ctx context.Context, userID int, guid string) (int, error)
	UpdateUserSecretsVersion(ctx context.Context, userID int, newVersion int, updatingDate time.Time) error

	Close() error
}

type Service struct {
	store      Store
	transactor transactor.Transactor
}

// NewService создает экземпляр объекта типа Service.
func NewService(store Store, transactor transactor.Transactor) *Service {
	return &Service{store: store, transactor: transactor}
}

func (s *Service) RegisterNewUser(ctx context.Context, login, password string) error {
	return s.store.SaveNewUser(ctx, login, password)
}

func (s *Service) UserIsValid(ctx context.Context, login, password string) (bool, error) {
	return s.store.UserIsValid(ctx, login, password)
}

func (s *Service) LoadSecret(ctx context.Context, secretData *types.SecretData, login string) error {
	userID, err := s.store.GetUserID(ctx, login)
	if err != nil {
		return err
	}
	err = s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		updatingDate := time.Now()
		err = s.validateSecretData(ctx, secretData, userID)
		if err != nil {
			return err
		}
		s.prepareSecretData(ctx, secretData, updatingDate)
		// Обновим секрет
		err = s.store.LoadSecret(ctx, secretData, userID)
		if err != nil {
			return err
		}
		return s.saveUserDataVersion(ctx, userID, updatingDate)
	})

	return err

}

func (s *Service) validateSecretData(ctx context.Context, secretData *types.SecretData, userID int) error {
	// Проверим guid
	if secretData.Guid == "" {
		return errors.New("SecredData ID can not be empty")
	}
	// Сравним версии, если не равны, то вернем ошибку
	versionID, err := s.store.GetSecretVersion(ctx, userID, secretData.Guid)
	if err != nil {
		return err
	}
	if versionID != secretData.VersionID {
		return errors.New("version ID doesn`t match. data was changed or deleted")
	}
	return nil
}

func (s *Service) prepareSecretData(ctx context.Context, secretData *types.SecretData, updatingDate time.Time) {
	secretData.VersionID++
	secretData.UpdatingDate = updatingDate
	if secretData.CreationDate.IsZero() {
		secretData.CreationDate = updatingDate
	}
}

// Сохраним версию данных пользователя, сейчас для упрощения в логике обмена не используется.
func (s *Service) saveUserDataVersion(ctx context.Context, userID int, updatingDate time.Time) error {
	// Получим текущую версию данных пользователя
	userSecretsVersion, err := s.store.GetUserSecretsVersion(ctx, userID)
	if err != sql.ErrNoRows && err != nil {
		return err
	}
	// Поднимем версию и сохраним ее
	userSecretsVersion++
	err = s.store.UpdateUserSecretsVersion(ctx, userID, userSecretsVersion, updatingDate)
	return err
}

func (s *Service) GetSecrets(ctx context.Context, login string, fromDate time.Time) (*[]types.SecretData, error) {
	userID, err := s.store.GetUserID(ctx, login)
	if err != nil {
		return nil, err
	}
	secrets, err := s.store.GetUserSecretList(ctx, userID, fromDate)
	return secrets, err
}
