package main

import (
	"log"
	"net/http"

	"github.com/dreamsofcode-io/terminal-ui/spinner"

	"github.com/J4yTr1n1ty/keyserver/boot"
	"github.com/J4yTr1n1ty/keyserver/middleware"
)

func main() {
	log.Println("Starting Keyserver...")
	spinner := spinner.New(spinner.Config{})
	spinner.Start()

	boot.LoadEnv()

	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
	server.ListenAndServe() // Blocking call to run the server
}
