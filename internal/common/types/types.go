package types

import "time"

type SecretData struct {
	Id           int
	Key          string
	Value        string
	BinaryValue  bool
	VersionID    int
	CreationDate time.Time
	UpdatingDate time.Time
	DeletionMark bool
}
