package main

import (
	"fmt"
)

const (
	slackToken = "xoxb-your-slack-token" // Put your Slack Bot Token here
	channelID  = "your-channel-id"       // ID of the Slack channel
	noaaAPIURL = "https://api.tidesandcurrents.noaa.gov/api/prod/datagetter"
	stationID  = "9414290" // Replace with your station ID
)

func main() {
	tideMessage, err := fetchTideData()
	if err != nil {
		fmt.Printf("Error fetching tide data: %v\n", err)
		return
	}

	err = sendTidesToSlack(tideMessage)
	if err != nil {
		fmt.Printf("Error sending message to Slack: %v\n", err)
	} else {
		fmt.Println("Tides successfully posted to Slack!")
	}
}
