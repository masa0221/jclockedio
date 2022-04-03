/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"fmt"
	"os"

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
			os.Exit(1)
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println("Can't read verbose flag: ", err)
			os.Exit(1)
		}

		// Clocked in/out
		jobcanClient := jobcan.New(config.Jobcan.Email, config.Jobcan.Password)
		jobcanClient.Verbose = verbose
		jobcanClient.NoAdit = noAdit
		aditResult := jobcanClient.Adit()

		// Output message
		outputMessage := generateOutputMessage(aditResult.Clock, aditResult.BeforeWorkingStatus, aditResult.AfterWorkingStatus)
		fmt.Println(outputMessage)

		// Send to Chatwork
		if config.Chatwork.Send {
			chatworkClient := chatwork.New(config.Chatwork.ApiToken)
			chatworkClient.Verbose = verbose
			_, err := chatworkClient.SendMessage(outputMessage, config.Chatwork.RoomId)
			if err != nil {
				fmt.Println("Failed to send to Chatwork")
				os.Exit(1)
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
	aditCmd.Flags().Bool("no-adit", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}

func generateOutputMessage(clock string, beforeStatus string, afterStatus string) string {
	return fmt.Sprintf("clock: %s, %s -> %s", clock, beforeStatus, afterStatus)
}
