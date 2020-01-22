package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startTicker() chan struct{} {
	ticker := time.NewTicker(time.Second)
	timer := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Print("Tick...")
			case <-timer:
				ticker.Stop()
				return
			}
		}
	}()

	return timer
}

func stopTicker(channel chan struct{}) {
	close(channel)
}

func initializeInterruptChannel() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	return interrupt
}

func waitForInterrupt(channel chan os.Signal) {
	killSignal := <-channel
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}
}

func main() {
	var name = os.Getenv("NAME")
	var portString = "8080"
	var portEnv = os.Getenv("PORT")

	if name == "" {
		name = "Service"
	}

	if portEnv != "" {
		portString = portEnv
	}

	log.Printf("Starting server '%s' on port %s...", name, portString)

	http.HandleFunc("/liveness", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Liveness?")
		_, _ = fmt.Fprint(w, "live")
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Readiness?")
		_, _ = fmt.Fprint(w, "ready")
	})

	http.HandleFunc("/die", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Dying...")
		os.Exit(1)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "A fancy 'Hello World' from server '%s'", name)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", portString)}

	ticker := startTicker()
	interrupt := initializeInterruptChannel()

	// starting server in the background
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	log.Print("Service is up and running")

	waitForInterrupt(interrupt)

	log.Print("The service is shutting down...")

	stopTicker(ticker)
	_ = server.Shutdown(context.Background())

	log.Print("Done")
}
