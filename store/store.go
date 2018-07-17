package store

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

var (
	// ErrSecretsNotFound is returned if the specified Secrets is not found
	ErrSecretsNotFound = errors.New("secrets not found")
)

type Secrets struct {
	Secrets RawSecrets
	Meta    *SecretsMetadata
}

var validKeyFormat = regexp.MustCompile(`^[A-Za-z0-9-_]+$`)
var validServiceFormat = regexp.MustCompile(`^[A-Za-z0-9-_]+$`)

func validateService(service string) error {
	if !validServiceFormat.MatchString(service) {
		return fmt.Errorf("Failed to validate service name '%s'.  Only alphanumeric, dashes, and underscores are allowed for service names", service)
	}
	return nil
}

func validateKey(key string) error {
	if !validKeyFormat.MatchString(key) {
		return fmt.Errorf("Failed to validate key name '%s'.  Only alphanumeric, dashes, and underscores are allowed for key names", key)
	}
	return nil
}

// A secret without any metadata
type RawSecrets map[string]string

func (r *RawSecrets) Validate() error {
	var result error
	for k, _ := range *r {
		err := validateKey(k)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
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
