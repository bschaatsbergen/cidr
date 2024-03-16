// Copyright (c) Bruno Schaatsbergen
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version string

	rootCmd = &cobra.Command{
		Use:     "cidr",
		Short:   "cidr - CLI to perform various actions on CIDR ranges",
		Version: version, // The version is set during the build by making using of `go build -ldflags`.
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func setupCobraUsageTemplate() {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgBlue).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func init() {
	setupCobraUsageTemplate()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
