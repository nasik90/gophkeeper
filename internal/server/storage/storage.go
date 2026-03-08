package storage

import (
	"errors"
	"time"
)

var (
	ErrUserNotUnique    = errors.New("user is not unique")
	ErrVersionIdNotTrue = errors.New("version ID is not true, record was changed by another client")
)

var zeroTime time.Time
