package cmd

import (
	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Returns the number of host addresses within the given cidr range",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("Provide a cidr range")
		}
		count(args[0])
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}

func count(arg string) {
	network, err := core.ParseCIDR(arg)
	if err != nil {
		logrus.Fatal("Provide a valid cidr range")
	}
	count := core.AddressCount(network)
	logrus.Infof("contains %d host addresses", count)
}
