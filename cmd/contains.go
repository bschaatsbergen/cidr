package cmd

import (
	"net"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var containsCmd = &cobra.Command{
	Use:   "contains",
	Short: "Checks whether a given ip address is part of a cidr range",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logrus.Fatal("Provide an ip address and a cidr range")
		}
		networkAddress, err := core.ParseCIDR(args[0])
		if err != nil {
			logrus.Fatal(err)
		}
		ipAddress := core.ParseIP(args[1])
		if ipAddress == nil {
			logrus.Fatal("Provided ip is not valid")
		}
		networkContainsIP := contains(networkAddress, ipAddress)
		logrus.Info(networkContainsIP)
	},
}

func init() {
	rootCmd.AddCommand(containsCmd)
}

func contains(networkAddress *net.IPNet, ipAddress net.IP) bool {
	contains := core.ContainsAddress(networkAddress, ipAddress)
	return contains
}
