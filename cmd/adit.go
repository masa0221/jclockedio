/*
Copyright © 2022 Masashi Tsuru

*/
package cmd

import (
	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/client/jobcan/browser"
	"github.com/masa0221/jclockedio/pkg/logger/chatwork_logger"
	"github.com/masa0221/jclockedio/pkg/logger/stdout_logger"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
	"github.com/masa0221/jclockedio/pkg/service/logging"
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
			log.Fatalf("Can't read no-adit flag: %v", err)
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Fatalf("Can't read verbose flag: %v", err)
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		// Clocked in/out
		browser, err := browser.NewAgoutiBrowser()
		if err != nil {
			log.Fatalf("Can't launch a browser: %v", err)
		}
		defer browser.Close()

		// Jobcan client
		credentials := &jobcan.JobcanCredentials{
			Email:    config.Jobcan.Email,
			Password: config.Jobcan.Password,
		}
		jobcanClient := jobcan.NewJobcanClient(browser, credentials)
		jobcanClient.NoAdit = noAdit

		// logger
		chatworkLogger := chatwork_logger.NewChatworkLogger(
			config.Chatwork.ApiToken,
			&chatwork_logger.Config{
				ToRoomId: config.Chatwork.RoomId,
				Unread:   false,
			})
		stdoutLogger := stdout_logger.NewStdoutLogger()

		// logging service
		loggingService := logging.NewLoggingService(
			chatworkLogger,
			stdoutLogger,
		)

		// clocked in / out
		clockIOConfig := &clockio.Config{
			LoggingEnabled:        config.Chatwork.Send,
			ClockedIOResultFormat: config.Output.Format,
		}
		clockIOService := clockio.NewClockIOService(jobcanClient, loggingService, clockIOConfig)
		result, err := clockIOService.Adit()
		if err != nil {
			log.Errorf("Failed to adit. reason: %v", err)
		}

		// output
		log.Debugf("[%s] %s -> %s", result.Clock, result.BeforeWorkingStatus, result.AfterWorkingStatus)
	},
}

func init() {
	rootCmd.AddCommand(aditCmd)
	aditCmd.Flags().Bool("no-adit", false, "It login to Jobcan using by configure, but no adit.(The adit means to push button of clocked in/out)")
}
