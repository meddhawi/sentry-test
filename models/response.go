package models

import (
	"time"

	"github.com/getsentry/sentry-go"
)

type Response struct {
	Error           bool             `json:"error"`
	Message         string           `json:"message"`
	DebugMessage    string           `json:"debug_message"`
	Data            interface{}      `json:"data"`
	RequestTime     time.Duration    `json:"request_time,omitempty"`
}

func CallModels() {
	breadcrumb := &sentry.Breadcrumb{
		Category: "models",
		Data: map[string]interface{}{
			"method": "CallModels",
		},
	}
	
	sentry.CurrentHub().AddBreadcrumb(breadcrumb, &sentry.BreadcrumbHint{})

	breadcrumb.Message = "An error occurred"

	
}