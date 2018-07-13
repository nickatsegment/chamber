package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/segmentio/chamber/store"
	"github.com/spf13/cobra"
)

var (
	version string
	bucket  string
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
	readCmd.Flags().StringVarP(&bucket, "bucket", "b", os.Getenv("CHAMBERS3_BUCKET"), "s3 bucket. Default: $CHAMBERS3_BUCKET")
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

	val, meta, err := secretStore.Read(service, key, version)
	if err != nil {
		return errors.Wrap(err, "Failed to read")
	}

	if quiet {
		fmt.Fprintf(os.Stdout, "%s\n", val)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Key\tValue\tVersion\tLastModified")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		key,
		val,
		meta.Version,
		meta.LastModified.Local().Format(ShortTimeFormat),
	)
	w.Flush()
	return nil
}
