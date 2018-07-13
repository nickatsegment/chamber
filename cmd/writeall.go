package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/chamber-s3/store"
	"github.com/spf13/cobra"
)

var (
	// writeAllCmd represents the writeall command
	writeAllCmd = &cobra.Command{
		Use:   "writeall <service>",
		Short: "write all secrets for a service",
		Args:  cobra.ExactArgs(1),
		RunE:  writeall,
	}
)

func init() {
	RootCmd.AddCommand(writeAllCmd)
}

func writeall(cmd *cobra.Command, args []string) error {
	service := strings.ToLower(args[0])
	if err := validateService(service); err != nil {
		return errors.Wrap(err, "Failed to validate service")
	}

	dec := json.NewDecoder(os.Stdin)
	rawsec := make(store.RawSecrets)
	if err := dec.Decode(&rawsec); err != nil {
		return err
	}

	secretStore := store.NewS3Store(numRetries, bucket, s3PathPrefix)
	newVersion, err := secretStore.WriteAll(service, rawsec)
	if err != nil {
		return err
	}
	fmt.Printf("Version: %s\n", newVersion)
	return nil
}
