package storage

import (
	"errors"
	"time"
)

type SecretData struct {
	Id           int
	Key          string
	Value        string
	VersionID    int
	CreationDate time.Time
	UpdatingDate time.Time
	DeletionMark bool
}

var (
	ErrUserNotUnique    = errors.New("user is not unique")
	ErrVersionIdNotTrue = errors.New("version ID is not true, record was changed by another client")
)

var zeroTime time.Time
