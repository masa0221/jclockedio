/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/masa0221/jclockedio/internal/chatwork"
	"github.com/masa0221/jclockedio/internal/jobcan"
	"github.com/spf13/cobra"
)

// aditCmd represents the adit command
var aditCmd = &cobra.Command{
	Use:   "adit",
	Short: "Clocked in/out with Jobcan",
	Long:  `Clocked in/out with Jobcan, then send message to Chatwork.(if you can the setting true)`,
	Run: func(cmd *cobra.Command, args []string) {
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
		if config.Output.Format != "" {
			// stdout
			outputMessage := generateOutputMessage(config.Output.Format, aditResult.Clock, aditResult.BeforeWorkingStatus, aditResult.AfterWorkingStatus)
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
		}
	},
}

func init() {
	rootCmd.AddCommand(aditCmd)
	aditCmd.Flags().Bool("no-adit", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}

func generateOutputMessage(outputFormat string, clock string, beforeStatus string, afterStatus string) string {
	assignData := map[string]interface{}{
		"clock":        clock,
		"beforeStatus": beforeStatus,
		"afterStatus":  afterStatus,
	}

	tpl, err := template.New("").Parse(outputFormat)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	writer := new(strings.Builder)
	if err := tpl.Execute(writer, assignData); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return writer.String()
}
