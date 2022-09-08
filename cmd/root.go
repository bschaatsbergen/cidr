package cmd

import (
	"github.com/bschaatsbergen/cidr/model"
	"github.com/bschaatsbergen/cidr/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version string

	FlagStore model.FlagStore

	rootCmd = &cobra.Command{
		Use:     "cidr",
		Short:   "cidr - cross platform cli to perform various operations on a cidr range",
		Version: version, // The version is set during the build by making using of `go build -ldflags`
		Run: func(cmd *cobra.Command, args []string) {
			utils.ConfigureLogLevel(FlagStore.Debug)
			cmd.Help()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&FlagStore.Debug, "debug", "d", false, "set log level to debug")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
