package cmd

import (
	"net"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var overlapsCmd = &cobra.Command{
	Use:   "overlaps",
	Short: "Checks if a CIDR range overlaps with another CIDR range",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logrus.Fatal("Provide 2 cidr ranges")
		}
		network1, err := core.ParseCIDR(args[0])
		if err != nil {
			logrus.Error(err)
		}
		network2, err := core.ParseCIDR(args[1])
		if err != nil {
			logrus.Error(err)
		}
		networksOverlap := overlaps(network1, network2)
		logrus.Info(networksOverlap)
	},
}

func init() {
	rootCmd.AddCommand(overlapsCmd)
}

func overlaps(network1, network2 *net.IPNet) bool {
	overlaps := core.Overlaps(network1, network2)
	return overlaps
}
