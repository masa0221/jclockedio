/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/masa0221/jclockedio/pkg/client/chatwork"
	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/client/jobcan/browser"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
	"github.com/masa0221/jclockedio/pkg/service/notification"
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
		browser, err := browser.NewAgoutiBrowser()
		if err != nil {
			fmt.Println("Can't launch a browser: ", err)
			os.Exit(1)
		}
		credentials := &jobcan.JobcanCredentials{
			Email:    config.Jobcan.Email,
			Password: config.Jobcan.Password,
		}
		notificationConfig := &notification.NotificationConfig{
			NotifyEnabled:         config.Chatwork.Send,
			ClockedIOResultFormat: config.Output.Format,
		}
		chatworkApiToken := config.Chatwork.ApiToken
		chatworkSendMessageConfig := &chatwork.ChatworkSendMessageConfig{
			ToRoomId: config.Chatwork.RoomId,
			Unread:   false,
		}

		jobcanClient := jobcan.NewJobcanClient(browser, credentials)
		chatworkClient := chatwork.NewChatworkClient(chatworkApiToken, chatworkSendMessageConfig)
		notificationService := notification.NewNotificationService(notificationConfig, chatworkClient)
		clockIOService := clockio.NewClockIOService(jobcanClient, notificationService)
		result, err := clockIOService.Adit()
		if err != nil {
			fmt.Println("Failed to send to Chatwork")
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(aditCmd)
	aditCmd.Flags().Bool("no-adit", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}
