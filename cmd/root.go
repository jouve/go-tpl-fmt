package cmd

import (
	"log"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-tpl-fmt",
	Short: "The uncompromising go-template formatter.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		run := formatFile
		if c, _ := cmd.Flags().GetBool("check"); c {
			run = checkFile
		}
		var result *multierror.Error
		for i, arg := range args {
			log.Printf("[%d/%d] %s", i+1, len(args), arg)
			if err := run(arg); err != nil {
				result = multierror.Append(result, err)
			}
		}
		return result.ErrorOrNil()
	},
	SilenceUsage: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("check", "c", false, "Don't write the files back, just return the status.")
}
