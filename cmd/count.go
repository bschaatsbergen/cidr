package cmd

import (
	"net"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	countCmd = &cobra.Command{
		Use:   "count",
		Short: "Return the count of distinct host addresses in a given CIDR range",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				logrus.Fatal("Provide a cidr range")
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				logrus.Fatal("Provide a valid cidr range")
			}
			hostAddressCount := count(network)
			logrus.Infof("contains %d host addresses", hostAddressCount)

		},
	}
)

func init() {
	rootCmd.AddCommand(countCmd)
}

func count(network *net.IPNet) uint64 {
	count := core.AddressCount(network)
	return count
}
