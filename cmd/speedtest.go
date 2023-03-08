/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"network-test/pkg"

	"github.com/spf13/cobra"
)

// speedtestCmd represents the speedtest command
var speedtestCmd = &cobra.Command{
	Use:   "speedtest",
	Short: "speedtest function",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Download()
	},
}

func init() {
	rootCmd.AddCommand(speedtestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// speedtestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// speedtestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
