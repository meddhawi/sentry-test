package middlewares

import (
	"encoding/json"
	"errors"
	// "errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/meddhawi/sentry-test/models"
)

func ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()
	response := new(models.Response)
	responseBody := c.Response().Body()
	json.Unmarshal(responseBody, &response)

	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetTags(map[string]string{
		"api": c.OriginalURL(),
		"endpoint":        c.Route().Path,
		"http.method":     c.Method(),
		"http.url":        c.Path(),
		"http.status_code": strconv.Itoa(c.Response().StatusCode()), // Set the actual status code
	})
	startTime := time.Now()

	if response.Error {
		log.Println("Error found!")
		duration := time.Since(startTime)
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			// Message:  "",
			Category: "middlewares",
			Data: map[string]interface{}{
					"method":      c.Method(),
					"url":         c.Path(),
					"status_code": c.Response().StatusCode(),
					"duration":    duration.String(),
				},
		}, &sentry.BreadcrumbHint{})
		var message string
		if response.DebugMessage != "" {
			message = response.DebugMessage
		} else {
			message = response.Message
		}

		exception := sentry.Exception{
			Type:  "API Error",
			Value: message,
			Stacktrace: sentry.ExtractStacktrace(errors.New(message)),
			// You can set the stack trace here using `sentry.ExtractStacktrace(err)`,
			// where `err` is the original error that caused the response error.
			// Stacktrace: sentry.ExtractStacktrace(err),
		}

		event := sentry.NewEvent()
		event.Exception = []sentry.Exception{exception}
		event.Transaction = c.Route().Path
		event.Level = sentry.LevelError

		hub.CaptureEvent(event)

		sentry.Flush(2 * time.Second)
	}
	
	return err
}

func PanicHandler(c *fiber.Ctx) error {
	startTime := time.Now()

	defer func() {
		if r := recover(); r != nil {
			// Recover from panic and capture the exception
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}

			// Capture the exception along with endpoint information
			event := sentry.NewEvent()
			event.Exception = []sentry.Exception{{
				Type:       "error",
				Value:      err.Error(),
				Stacktrace: sentry.ExtractStacktrace(err),
			}}
			event.Transaction = c.Route().Path
			event.Level = sentry.LevelError
			event.Tags = map[string]string{
				"endpoint":        c.Route().Path,
				"http.method":     c.Method(),
				"http.url":        c.Path(),
				"http.status_code": strconv.Itoa(c.Response().StatusCode()), // Set the actual status code
			}
			sentry.CaptureEvent(event)

			// Calculate the request duration
			duration := time.Since(startTime)
			sentry.AddBreadcrumb(&sentry.Breadcrumb{
				Category: "request",
				Type:     "http",
				Data: map[string]interface{}{
					"method":      c.Method(),
					"url":         c.Path(),
					"status_code": c.Response().StatusCode(),
					"duration":    duration.String(),
				},
			})

			sentry.Flush(2 * time.Second)
		}
	}()

	return c.Next()
}