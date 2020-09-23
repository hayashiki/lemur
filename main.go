package main

import (
	"fmt"
	"github.com/hayashiki/lemur/app"
	"github.com/hayashiki/lemur/config"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %#v", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := config.NewReadMustFromEnv()

	if err != nil {
		log.Fatalf("Failed to load config %v", err)
		return err
	}

	s := app.NewServer(config)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: s.Router(),
	}
	log.Print("Listening requests on :" + config.Port + "...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe failed on :" + config.Port)
		return err
	}

	log.Print(fmt.Sprintf("Listening on port %s", config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))

	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello world with hot reload!")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
