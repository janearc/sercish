package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
)

type Service struct {
	config *Config
	api    *slack.Client
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
	s.aboutHandler()
	s.slackHandler()

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

	return nil
}

func (s *Service) slackHandler() {
	http.HandleFunc("/slack", func(w http.ResponseWriter, r *http.Request) {
		// TODO: "challenge" parameter
		fmt.Fprintf(w, "Hello, Slack!")
	})
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

	return
}
