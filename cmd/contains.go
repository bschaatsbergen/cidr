package cmd

import (
	"net"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	containsCmd = &cobra.Command{
		Use:   "contains",
		Short: "Checks whether an IP address belongs to a CIDR range",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				logrus.Fatal("Provide an ip address and a cidr range")
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				logrus.Fatal(err)
			}
			ip := core.ParseIP(args[1])
			if ip == nil {
				logrus.Fatal("Provided ip is not valid")
			}
			networkContainsIP := contains(network, ip)
			logrus.Info(networkContainsIP)
		},
	}
)

func init() {
	rootCmd.AddCommand(containsCmd)
}

func contains(network *net.IPNet, ip net.IP) bool {
	contains := core.ContainsAddress(network, ip)
	return contains
}
