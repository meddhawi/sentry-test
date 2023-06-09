package main

import (
	// "encoding/json"
	"errors"
	"log"

	// "github.com/meddhawi/sentry-test/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/meddhawi/sentry-test/sentryutil"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}
func main() {
	err := errors.New("It didn't work!")
	hub := sentry.CurrentHub().Clone()
	eventID := hub.CaptureException(err)
	log.Println("Event ID:", eventID)

	// You can also capture warnings or other log messages using the same package.
	// sentryutil.CaptureMessage("This is a warning.")

	// Flush buffered events before the program terminates.
	defer sentryutil.FlushSentry()
}
