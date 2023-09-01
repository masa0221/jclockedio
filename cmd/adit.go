/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"github.com/masa0221/jclockedio/pkg/client/chatwork"
	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/client/jobcan/browser"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
	"github.com/masa0221/jclockedio/pkg/service/notification"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// aditCmd represents the adit command
var aditCmd = &cobra.Command{
	Use:   "adit",
	Short: "Clocked in/out with Jobcan",
	Long:  `Clocked in/out with Jobcan, then send message to Chatwork.(if you can the setting true)`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})

		noAdit, err := cmd.Flags().GetBool("no-adit")
		if err != nil {
			log.Errorf("Can't read no-adit flag: ", err)
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Errorf("Can't read verbose flag: ", err)
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		// Clocked in/out
		browser, err := browser.NewAgoutiBrowser()
		if err != nil {
			log.Errorf("Can't launch a browser: ", err)
		}
		defer browser.Close()

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
		jobcanClient.NoAdit = noAdit

		chatworkClient := chatwork.NewChatworkClient(chatworkApiToken, chatworkSendMessageConfig)
		notificationService := notification.NewNotificationService(notificationConfig, chatworkClient)
		clockIOService := clockio.NewClockIOService(jobcanClient, notificationService)

		result, err := clockIOService.Adit()
		if err != nil {
			log.Errorf("Failed to adit. reason: %v", err)
		}

		log.Info(result)
	},
}

func init() {
	rootCmd.AddCommand(aditCmd)
	aditCmd.Flags().Bool("no-adit", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}
