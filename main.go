package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	c, cerr := LoadConfig(
		"/app/config/config.yml",
		"/app/config/version.yml",
		"/app/config/secrets.yml",
	)

	if cerr != nil {
		logrus.WithError(cerr).Fatalf("Failed to load config: %v", cerr)
	}

	s := NewService(c)

	if s == nil {
		logrus.Fatalf("Failed to instantiate service object")
	}

	logrus.Info("Service object instantiated")

	// so we have a healthy service object, let's start it
	serr := s.Start()

	if serr != nil {
		logrus.WithError(serr).Fatalf("Failed to start service: %v", serr)
	}

	logrus.Info("Service started")

	tideMessage, err := s.fetchTideData()
	if err != nil {
		logrus.WithError(err).Warn("Error fetching tide data")
		return
	}

	err = s.sendTidesToSlack(tideMessage)
	if err != nil {
		logrus.Warnf("Error sending message to Slack: %v\n", err)
	} else {
		logrus.Infof("Tides successfully posted to Slack!")
	}

	// TODO: does this mean we just fall out the bottom?
	logrus.Fatal("Exiting main.go")
}
