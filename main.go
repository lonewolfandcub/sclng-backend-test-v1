package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/github"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)
	router.HandleFunc("/repos", reposHandler).Methods(http.MethodGet)
	router.HandleFunc("/stats", statsHandler).Methods(http.MethodGet)

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}

func reposHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())

	ghClient := github.NewClient()
	repos, err := ghClient.ListLatestRepositories()
	if err != nil {
		log.WithError(err).Error("Fail to list repositories")
	}

	body, err := json.Marshal(repos)
	if err != nil {
		log.WithError(err).Error("Fail to encode response")
		body = []byte{}
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())

	ghClient := github.NewClient()
	body, err := ghClient.GatherLatestRepositoriesStats()

	if err != nil {
		log.WithError(err).Error("Fail to gather repositories statistics")
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)

	return nil
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
