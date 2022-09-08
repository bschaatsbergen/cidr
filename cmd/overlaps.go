package cmd

import (
	"net"

	"github.com/bschaatsbergen/cidr/pkg/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var overlapsCmd = &cobra.Command{
	Use:   "overlaps",
	Short: "Checks whether 2 cidr ranges overlap",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			logrus.Fatal("Provide 2 cidr ranges")
		}
		n1, err := core.ParseCIDR(args[0])
		if err != nil {
			logrus.Error(err)
		}
		n2, err := core.ParseCIDR(args[1])
		if err != nil {
			logrus.Error(err)
		}
		doesOverlap := overlaps(n1, n2)
		logrus.Info(doesOverlap)
	},
}

func init() {
	rootCmd.AddCommand(overlapsCmd)
}

func overlaps(n1, n2 *net.IPNet) bool {
	overlaps := core.Overlaps(n1, n2)
	return overlaps
}
