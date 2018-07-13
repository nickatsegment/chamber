package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/chamber/store"
	"github.com/spf13/cobra"
)

var (
	version string
	quiet   bool

	// readCmd represents the read command
	readCmd = &cobra.Command{
		Use:   "read <service> <key>",
		Short: "Read a specific secret from s3",
		Args:  cobra.ExactArgs(2),
		RunE:  read,
	}
)

func init() {
	readCmd.Flags().StringVarP(&version, "version", "v", "", "The version number of the secret. Defaults to latest.")
	readCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only print the secret")
	RootCmd.AddCommand(readCmd)
}

func read(cmd *cobra.Command, args []string) error {
	if bucket == "" {
		return errors.New("bucket not set")
	}
	service := strings.ToLower(args[0])
	if err := validateService(service); err != nil {
		return errors.Wrap(err, "Failed to validate service")
	}

	key := strings.ToLower(args[1])
	if err := validateKey(key); err != nil {
		return errors.Wrap(err, "Failed to validate key")
	}

	// TODO: pass prefix
	secretStore := store.NewS3Store(numRetries, bucket, "")
	secrets, err := secretStore.ReadAll(service, version)
	if err != nil {
		return errors.Wrap(err, "Failed to read")
	}
	val, ok := secrets.Secrets[key]
	if !ok {
		return errors.New("key not found")
	}

	if quiet {
		fmt.Fprintf(os.Stdout, "%s\n", val)
		return nil
	}

	// filter out everything else
	for k, _ := range secrets.Secrets {
		if k != key {
			delete(secrets.Secrets, k)
		}
	}

	printSecrets(secrets, true)
	return nil
}
