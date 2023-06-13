package main

import (
	// "errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/meddhawi/sentry-test/middlewares"
	"github.com/meddhawi/sentry-test/models"
	"github.com/meddhawi/sentry-test/sentryutil"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	sentryutil.Initialize()
}
func main() {
	var err error
	app := fiber.New()
	setUpRoutes(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	port := ":3000"
	// Serve server
	if err := app.Listen(port); err != nil {
		fmt.Fprintf(os.Stdout, "Couldn't start the server: %v\n", err) // should add sentry here
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		err = app.Shutdown()
	}()
	log.Fatal(err)
	
	// You can also capture warnings or other log messages using the same package. 
	// sentryutil.CaptureMessage("This is a warning.")
	// Flush buffered events before the program terminates.
	defer sentryutil.FlushSentry()
}

func setUpRoutes(app *fiber.App){
	// api := app.Group("api")

	app.Use(middlewares.PanicHandler)
	app.Use(middlewares.ErrorHandler)
	app.Get("/panic", func(c *fiber.Ctx) error {
		log.Println("test")
		sentryutil.SetPanic()
		return c.JSON(map[string]interface{}{
			"message": "Hello",
		})
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		sentry.CurrentHub().AddBreadcrumb(&sentry.Breadcrumb{
			Category: "controllers",
		}, &sentry.BreadcrumbHint{})
		// sentry.CurrentHub().AddBreadcrumb()
		models.CallModels()
		// sentry.NewScope().SetTag("api")
		return c.Status(http.StatusUnprocessableEntity).JSON(models.Response{
			Error: true,
			Message: "Something is wong!",
		})
	})	
}