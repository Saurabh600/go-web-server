package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Saurabh600/go-web-server/internals/config"
	api_controllers "github.com/Saurabh600/go-web-server/internals/controllers/api"
	pages_controllers "github.com/Saurabh600/go-web-server/internals/controllers/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/joho/godotenv"
)

func main() {
	// reading .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("reading .env file failed, error: ", err)
	}

	// database connection
	config.InitilizeDb()

	// handling abort
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := config.GetDb().Close(); err != nil {
			fmt.Println("error closing db connection, err: ", err)
		}
		os.Exit(1)
	}()

	// fiber app
	app := fiber.New(fiber.Config{
		// Prefork:       true,
		StrictRouting: false,
		CaseSensitive: false,
		Views:         handlebars.New("./template", ".html"),
		ServerHeader:  "Fiber",
		AppName:       "HTMX Web Server 0.1.0",
	})

	// home route
	app.Get("/", pages_controllers.HomePage)

	// api route
	apiRouteV1 := app.Group("/api/v1")
	{
		apiRouteV1.Get("/hi", api_controllers.CheckStatus)
		apiRouteV1.Get("/auth", api_controllers.CheckAuth)
		apiRouteV1.Get("/users", api_controllers.GetAllUsers)
		apiRouteV1.Post("/users/new", api_controllers.CreateNewUser)
	}

	// starting server
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatal("failed to start server, error: ", err)
	}
}
