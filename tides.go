package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// NOAA Response structure
type TidePrediction struct {
	Predictions []struct {
		Time  string `json:"t"`
		Type  string `json:"type"`
		Value string `json:"v"`
	} `json:"predictions"`
}

func (s *Service) fetchTideData() (string, error) {
	url := fmt.Sprintf(
		"%s?product=predictions&datum=MLLW&station=%s&time_zone=lst_ldt&units=english&interval=hilo&format=json",
		s.config.NOAA.apiURL,
		s.config.NOAA.stationID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}

	var tideData TidePrediction
	err = json.Unmarshal(buf.Bytes(), &tideData)
	if err != nil {
		return "", err
	}

	// Let's format the tide data
	message := "Tide Predictions:\n"
	for _, prediction := range tideData.Predictions {
		message += fmt.Sprintf("Time: %s - Type: %s - Height: %s feet\n", prediction.Time, prediction.Type, prediction.Value)
	}
	return message, nil
}

func (s *Service) sendTidesToSlack(message string) error {
	//_, _, err := api.PostMessage(
	//	channelID,
	//	slack.MsgOptionText(message, false),
	//)
	return errors.New("Not implemented")
}
