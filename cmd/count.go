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
		if len(args) < 1 {
			logrus.Fatal("cidr range argument is required")
		}
		main(args[0])
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}

func main(arg string) {
	network, err := core.ParseCIDR(arg)
	if err != nil {
		logrus.Error(err)
	}

	count := core.AddressCount(network)

	logrus.Infof("%s contains %d distinct host addresses", network, count)
}
