package service

import "context"

// SendSecrets отправляет локальные секреты на сервер.
// На сервер будут отправлены добавленные и измененные на клиенте секреты.
func (s *Service) SendSecrets(ctx context.Context) error {
	// Читает из БД записи, необходимые отправить
	// Отправляет их
	// Фиксирует в БД серверный айдишник секрета
	return nil
}

// UploadSecrets загружает с сервера секреты.
// С сервера будут получены добавленные и измененные на сервере секреты.
func (s *Service) UploadSecrets(ctx context.Context) error {
	// В цикле происходит вызов s.apiCleint.UploadSecrets()
	// Полученные данные записать в локальный БД.
	return nil
}
