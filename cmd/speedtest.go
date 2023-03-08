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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		threads, _ := cmd.Flags().GetInt("multi-threads")
		size, _ := cmd.Flags().GetString("size")
		proxy, _ := cmd.Flags().GetString("proxy")
		if threads <= 0 {
			threads = 1
		}
		pkg.Download(threads, size, proxy)
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
	speedtestCmd.Flags().IntP("multi-threads", "t", 1, "number of threadings")
	speedtestCmd.Flags().StringP("size", "s", "10", "size of the test file to download. Allowed 1/5/10/100 MB")
	speedtestCmd.Flags().StringP("proxy", "p", "", "specify Proy to use. Format: http://ip:port or socks5://ip:port")
}
