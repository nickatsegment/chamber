package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	numRetries       int
	chamberS3Version string
	bucket           string
	s3PathPrefix     string
)

const DefaultPrefix string = "chamber/"

const (
	// ShortTimeFormat is a short format for printing timestamps
	ShortTimeFormat = "01-02 15:04:05"

	// DefaultNumRetries is the default for the number of retries we'll use for our SSM client
	DefaultNumRetries = 10
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "chamber-s3",
	Short:        "CLI for storing secrets",
	SilenceUsage: true,
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&numRetries, "retries", "r", DefaultNumRetries, "number of retries we'll make before giving up")
	RootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", os.Getenv("CHAMBER_S3_BUCKET"), "s3 bucket. Default from $CHAMBER_S3_BUCKET")
	RootCmd.PersistentFlags().StringVarP(&s3PathPrefix, "prefix", "p", DefaultPrefix, "s3 path prefix")
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(vers string) {
	chamberS3Version = vers
	if cmd, err := RootCmd.ExecuteC(); err != nil {
		if strings.Contains(err.Error(), "arg(s)") || strings.Contains(err.Error(), "usage") {
			cmd.Usage()
		}
		os.Exit(1)
	}
}
