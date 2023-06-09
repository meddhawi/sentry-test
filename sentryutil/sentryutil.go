package sentryutil

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		// TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Modify the event here
			// event.User.Email = "" // Don't send user's email address
			// test, _ := json.Marshal(event)
			// log.Println(string(test))
			return event
		},
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

// CaptureError captures an error and returns the corresponding event ID.
func CaptureError(err error) *sentry.EventID {
	eventID := sentry.CaptureException(err)
	return eventID
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

