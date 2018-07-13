package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/chamber-s3/store"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <service>",
	Short: "List the secrets set for a service",
	Args:  cobra.ExactArgs(1),
	RunE:  list,
}

var (
	withValues bool
)

func init() {
	listCmd.Flags().BoolVarP(&withValues, "expand", "e", false, "Expand parameter list with values")
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	if bucket == "" {
		return errors.New("bucket not set")
	}
	service := strings.ToLower(args[0])
	if err := validateService(service); err != nil {
		return errors.Wrap(err, "Failed to validate service")
	}

	secretStore := store.NewS3Store(numRetries, bucket, s3PathPrefix)
	secrets, err := secretStore.ReadAll(service, version)
	if err != nil {
		return errors.Wrap(err, "Failed to read")
	}

	fmt.Printf("Version: %s\n", secrets.Meta.Version)
	fmt.Printf("LastModified: %s\n", secrets.Meta.LastModified.Local().Format(ShortTimeFormat))
	fmt.Println()

	for k, v := range secrets.Secrets {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
