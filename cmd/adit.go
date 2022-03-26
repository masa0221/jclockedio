/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"fmt"

	"github.com/masa0221/jclockedio/internal/chatwork"
	"github.com/masa0221/jclockedio/internal/jobcan"
	"github.com/spf13/cobra"
)

// aditCmd represents the adit command
var aditCmd = &cobra.Command{
	Use:   "adit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("adit called")
		noAdit, err := cmd.Flags().GetBool("no-adit")
		if err != nil {
			fmt.Println("Can't read no-adit flag: ", err)
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println("Can't read verbose flag: ", err)
		}

		if noAdit {
			fmt.Println("debug")
		} else {
			jobcan.Adit()
			fmt.Println("prod")
		}

		if config.Chatwork.Send {
			chatworkClient := chatwork.New(config.Chatwork.ApiToken)
			chatworkClient.Verbose = verbose
			messageId, err := chatworkClient.SendMessage("hoge", config.Chatwork.RoomId)
			if err != nil {
				fmt.Println("Failed to send to Chatwork: ", err)
			} else {
				fmt.Println("messageId: ", messageId)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(aditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aditCmd.Flags().BoolP("no-adit", "n", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}
