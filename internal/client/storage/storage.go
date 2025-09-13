package storage

import "time"

type SecretData struct {
	Id           int
	Key          string
	Value        string
	VersionID    int
	CreationDate time.Time
	UpdatingDate time.Time
	DeletionMark bool
}
