package main

import (
	"fmt"
	"net/http"
)

type Service struct {
	config *Config
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) Start() error {
	s.aboutHandler()
	s.slackHandler()

	http.ListenAndServe(fmt.Sprintf(":%d", s.config.Web.Port), nil)

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
				<h3>ATC</h3>
				<p>ATC is a web application that helps athletes track their performance and progress in swimming, biking, and running.</p>
				<p>author: Jane Arc</p>
				<p>Build Version: %s</p>
				<p>Build Date: %s</p>
				<p>source: <a href="http://github.com/janearc/atc">http://github.com/janearc/atc</a></p>
			</div>
		</body>
		</html>
		`, s.Config.Build.Build, s.Config.Build.BuildDate)

		w.Write([]byte(html))
	})

	return
}
