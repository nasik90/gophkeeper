package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/Thiht/transactor"
	//"github.com/nasik90/gophkeeper/internal/server/storage"
)

type Store interface {
	SaveNewUser(ctx context.Context, user, password string) error
	UserIsValid(ctx context.Context, login, password string) (bool, error)
	SaveNewSecret(ctx context.Context, key, value, login string, creationDate time.Time) (int, error)
	GetUserSecretsVersion(ctx context.Context, login string) (int, error)
	UpdateUserSecretsVersion(ctx context.Context, login string, newVersion int, updatingDate time.Time) error
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

func (s *Service) LoadSecret(ctx context.Context, key, value, login string) (int, error) {
	var recordID int
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		var err error
		creationDate := time.Now()
		// Сохраним новый секрет
		recordID, err = s.store.SaveNewSecret(ctx, key, value, login, creationDate)
		if err != nil {
			return err
		}
		// Получим текущую версию данных пользователя
		userSecretsVersion, err := s.store.GetUserSecretsVersion(ctx, login)
		if err != sql.ErrNoRows && err != nil {
			return err
		}
		// Поднимем версию и сохраним ее
		userSecretsVersion++
		err = s.store.UpdateUserSecretsVersion(ctx, login, userSecretsVersion, creationDate)
		return err
	})

	return recordID, err

}
