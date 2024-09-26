package main

import (
	"encoding/json"
	"fmt"
	"github.com/janearc/sux/sux"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
)

type Service struct {
	config *Config
	api    *slack.Client
	sux    *sux.Sux
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) GetConfig() *Config {
	return s.config
}

func (s *Service) Start() error {
	logrus.Info("Starting service")

	// there's probably a better way to wrap these up but i don't feel like
	// making bitey that complicated. we don't even have namespaces.
	s.aboutHandler()
	s.slackHandler()
	s.eventsHandler()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Web.Port), nil); err != nil {
		logrus.WithError(err).Fatal("Server failed to start")
	}

	logrus.Infof("Listener started on port %d", s.config.Web.Port)

	api := slack.New(s.GetConfig().OAuthToken())
	if api == nil {
		logrus.Fatal("Failed to instantiate Slack API object")
	}
	if _, err := api.AuthTest(); err != nil {
		logrus.WithError(err).Fatal("Failed to authenticate to Slack")
	}

	logrus.Info("Slack API object instantiated")

	s.api = api

	// TODO: yep this is gross too but that's okay, this is just bitey
	sux := sux.NewSux(s.config.Sources.cfg, s.config.Sources.ver, s.config.Sources.sec)

	if sux == nil {
		logrus.Fatal("Failed to instantiate Sux object")
	} else {
		logrus.Infof("UX is now stateful shades dot gif party parrot")
	}
	
	return nil
}

func (s *Service) slackHandler() {
	http.HandleFunc("/slack", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			logrus.Errorf("Failed to decode request body: %v", err)
			return
		}

		// Check for the "challenge" parameter
		if challenge, ok := body["challenge"].(string); ok {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{"challenge": challenge}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				logrus.Errorf("Failed to encode challenge response: %v", err)
				return
			}
			logrus.Infof("Responded to Slack challenge: %s", challenge)
			return
		}

		logrus.Infof("Received request: %v", r)
	})

	logrus.Info("handler: /slack")
}

// slack will tell us about assorted events, which maybe we care about
func (s *Service) eventsHandler() {
	http.HandleFunc("/slack/event", func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Received event: %v", r)
	})
	logrus.Info("handler: /slack/event")
}

// returns information about the service
func (s *Service) aboutHandler() {
	// handle the "about" request
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Activity Dashboard</title>
			<style>
				/* TODO: */
			</style>
		</head>
		<body>
			<div class="container">
				<h3>bitey</h3>
				<p>bitey is an artificial harbor seal that supports the sercish slack.</p>
				<p>author: Jane Arc</p>
				<p>Build Version: %s</p>
				<p>Build Date: %s</p>
				<p>source: <a href="http://github.com/janearc/sercish">http://github.com/janearc/sercish</a></p>
			</div>
		</body>
		</html>
		`, s.config.Version.Build, s.config.Version.BuildDate)

		if _, err := w.Write([]byte(html)); err != nil {
			logrus.WithError(err).Warn("Error writing response: %v", err)
		}
	})

	logrus.Info("handler: /about")
	return
}
