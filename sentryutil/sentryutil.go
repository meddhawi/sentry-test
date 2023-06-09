package sentryutil

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://568664b21fe7412ebac4a5158e8cea47@sentry.teamdev.id/3",
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

