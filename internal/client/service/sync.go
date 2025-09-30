package service

import (
	"context"
	"time"
)

// SendSecrets отправляет локальные секреты на сервер.
// На сервер будут отправлены добавленные и измененные на клиенте секреты.
func (s *Service) SendSecrets(ctx context.Context) error {
	// Читаем из БД записи, которые необходимые отправить
	toSend := true
	secrets, err := s.store.GetSecrets(ctx, toSend)
	if err != nil {
		return err
	}
	// Отправляем
	for _, secretData := range *secrets {
		err := s.apiCleint.SendSecret(&secretData)
		if err != nil {
			return err
		}
		// Фиксирует в БД серверный айдишник секрета
		secretData.ToSend = false
		err = s.store.UpdateSecret(ctx, &secretData)
		if err != nil {
			return err
		}
	}

	return nil
}

// UploadSecrets загружает с сервера секреты.
// С сервера будут получены добавленные и измененные на сервере секреты.
func (s *Service) UploadSecrets(ctx context.Context) error {

	// Получим dataVersion
	fromDate, err := s.store.GetDataVersion(ctx)
	if err != nil {
		return err
	}

	// Забираем данные с сервера
	secrets, err := s.apiCleint.UploadSecrets(fromDate)
	if err != nil {
		return err
	}

	var dataVersion time.Time
	// Полученные данные обновляем/добавляем в локальный БД
	for _, secret := range *secrets {
		err = s.store.InsertUpdateSecret(ctx, &secret)
		if err != nil {
			return err
		}
		// Определеим новую dataVersion
		if secret.UpdatingDate.After(dataVersion) {
			dataVersion = secret.UpdatingDate
		}
	}

	// Сохраним dataVersion
	if dataVersion != fromDate {
		err = s.store.SaveDataVersion(ctx, dataVersion)
	}

	return err
}
