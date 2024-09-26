package main

import (
	"fmt"
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

	// so we have a healthy service object, let's start it
	s.Start()

	tideMessage, err := s.fetchTideData()
	if err != nil {
		logrus.WithError(err).Warn("Error fetching tide data")
		return
	}

	err = s.sendTidesToSlack(tideMessage)
	if err != nil {
		fmt.Printf("Error sending message to Slack: %v\n", err)
	} else {
		fmt.Println("Tides successfully posted to Slack!")
	}
}
