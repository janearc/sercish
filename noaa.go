package main

// NOAA Response structure
type TidePrediction struct {
	Predictions []struct {
		Time  string `json:"t"`
		Type  string `json:"type"`
		Value string `json:"v"`
	} `json:"predictions"`
}
