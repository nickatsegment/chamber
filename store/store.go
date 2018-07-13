package store

import (
	"errors"
	"time"
)

var (
	// ErrSecretsNotFound is returned if the specified Secrets is not found
	ErrSecretsNotFound = errors.New("secrets not found")
)

type Secrets struct {
	Secrets RawSecrets
	Meta    *SecretMetadata
}

// A secret without any metadata
type RawSecrets map[string]string

type SecretMetadata struct {
	Version      string
	LastModified time.Time
}

type Store interface {
	WriteAll(id string, secrets RawSecrets) error
	Write(id, key, value string) error
	ReadAll(id, version string) (*Secrets, error)
	Read(id, key, version string) (string, *SecretMetadata, error)
	Delete(id string) error
}
