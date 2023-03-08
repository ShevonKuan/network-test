/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"network-test/pkg"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if the proxy is working",
	Long: `Check if the proxy is working, you can pass multiple proxy addresses
	written in the following format:
		- socks5://
		- http://
	If you don't specify the test server, http://222.201.187.66:10098 will be used by default
`,
	Example: `network-test check socks5://127.0.0.1:6001 socks5://127.0.0.1:6002 -t`,
	Run: func(cmd *cobra.Command, args []string) {
		t, _ := cmd.Flags().GetBool("multi-thread")
		server, _ := cmd.Flags().GetString("server")
		ch := make(chan *pkg.CheckResult)
		if len(args) == 0 {
			log.Fatal("😟 Error: no proxy address provided")
		} else {
			if t {
				log.Println("❗ Warning: multi-thread mode has been enabled.")

				for _, address := range args {
					go pkg.CheckProxy(ch, address, server)
				}
				for _, address := range args {
					result := <-ch
					if result.Err != nil {
						log.Printf("😟 Error <%s>:\t%v", address, result.Err)
					} else {
						log.Printf("😀 Passed <%s>:\tstatus code: %d\tduration: %v", address, result.StatusCode, result.Duration)
					}
				}
			} else {
				for _, address := range args {

					go pkg.CheckProxy(ch, address, server)
					result := <-ch
					if result.Err != nil {
						log.Printf("😟 Error <%s>:\t%v", address, result.Err)
					} else {
						log.Printf("😀 Passed <%s>:\tstatus code: %d\tduration: %v", address, result.StatusCode, result.Duration)
					}
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	checkCmd.Flags().BoolP("multi-thread", "t", false, "Enable multi-thread mode")
	checkCmd.Flags().StringP("server", "s", "http://222.201.187.66:10098", "specify the test server")
}
