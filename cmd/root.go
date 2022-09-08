package cmd

import (
	"github.com/bschaatsbergen/cidr/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Flags struct {
	Debug bool
}

var (
	version string

	flags Flags

	rootCmd = &cobra.Command{
		Use:     "cidr",
		Short:   "cidr - cross platform cli to perform various operations on a cidr range",
		Version: version, // The version is set during the build by making using of `go build -ldflags`
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&flags.Debug, "debug", "d", false, "set log level to debug")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := utils.ConfigureLogLevel(flags.Debug); err != nil {
			return err
		}
		return nil
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
