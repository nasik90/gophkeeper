package types

import (
	"time"
)

// TODO теги
type SecretData struct {
	Guid         string    `json:"guid"`
	Key          []byte    `json:"key"`
	Value        []byte    `json:"value"`
	BinaryValue  bool      `json:"binaryValue"`
	VersionID    int       `json:"versionID"`
	CreationDate time.Time `json:"creationDate"`
	UpdatingDate time.Time `json:"updatingDate"`
	DeletionMark bool      `json:"deletionMark"`
	Comment      string    `json:"comment"`
	ToSend       bool      `json:"toSend"`
}
