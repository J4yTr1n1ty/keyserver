package main

import (
	"log"
	"net/http"

	"github.com/dreamsofcode-io/terminal-ui/spinner"

	"github.com/J4yTr1n1ty/keyserver/pkg/api"
	"github.com/J4yTr1n1ty/keyserver/pkg/boot"
	"github.com/J4yTr1n1ty/keyserver/pkg/middleware"
	"github.com/J4yTr1n1ty/keyserver/pkg/models"
)

func main() {
	log.Println("Starting Keyserver...")
	spinner := spinner.New(spinner.Config{})
	spinner.Start()

	boot.LoadEnv()
	boot.ConnectToDatabase()
	err := boot.DB.AutoMigrate(&models.Key{}, &models.Identity{})
	if err != nil {
		log.Fatal(err)
	}

	router := api.SetupRoutes()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
		w.WriteHeader(http.StatusOK)
	})

	// Add middleware
	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    ":" + boot.Environment.GetEnv("PORT"),
		Handler: stack(router), // Wrap the router with the middleware stack created above
	}

	spinner.Stop()
	log.Println("Listening on port :" + boot.Environment.GetEnv("PORT"))
	if err := server.ListenAndServe(); err != nil { // Blocking call to run the server
		log.Fatal("Error starting server: " + err.Error())
	}
}
