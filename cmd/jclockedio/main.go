package main

import (
	"fmt"
	"log"

	"github.com/masa0221/jclockedio/internal/chatwork"
	"github.com/masa0221/jclockedio/internal/jobcan"
)

func main() {
	result := jobcan.ClockedInOut()

	message := fmt.Sprintf("打刻 %s (%s → %s)", result.ClockTime, result.BeforeStatus, result.AfterStatus)
	log.Printf(message)

	chatwork.Send(message)
}
