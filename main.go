package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/slack-go/slack"
)

const (
    slackToken    = "xoxb-your-slack-token" // Put your Slack Bot Token here
    channelID     = "your-channel-id"       // ID of the Slack channel
    noaaAPIURL    = "https://api.tidesandcurrents.noaa.gov/api/prod/datagetter"
    stationID     = "9414290"               // Replace with your station ID
)

// NOAA Response structure
type TidePrediction struct {
    Predictions []struct {
        Time  string `json:"t"`
        Type  string `json:"type"`
        Value string `json:"v"`
    } `json:"predictions"`
}

func fetchTideData() (string, error) {
    url := fmt.Sprintf("%s?product=predictions&datum=MLLW&station=%s&time_zone=lst_ldt&units=english&interval=hilo&format=json", noaaAPIURL, stationID)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var tideData TidePrediction
    err = json.Unmarshal(body, &tideData)
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

func sendTidesToSlack(message string) error {
    api := slack.New(slackToken)
    _, _, err := api.PostMessage(
        channelID,
        slack.MsgOptionText(message, false),
    )
    return err
}

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
