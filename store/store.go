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
	Meta    *SecretsMetadata
}

// A secret without any metadata
type RawSecrets map[string]string

func (r *RawSecrets) Validate() error {
	// TODO; need to move validateKey here
	panic("not implemented")
}

type SecretsMetadata struct {
	Version      string
	LastModified time.Time
}

type Store interface {
	// Write all secrets for secrets `id`
	WriteAll(id string, secrets RawSecrets) (version string, err error)
	// Write one secret `key` for secrets `id`, returning the new Secrets
	Write(id, key, value string) (version string, err error)
	// Read all secrets for `id` at `version`, returning the new Secrets
	ReadAll(id, version string) (*Secrets, error)
	// Read one secret `key` for `id` at `version`
	Read(id, key, version string) (string, *SecretsMetadata, error)
	// Delete all secrets for `id`
	DeleteAll(id string) error
	// Delete one secret for `id` at `key`
	Delete(id, key string) error
}
