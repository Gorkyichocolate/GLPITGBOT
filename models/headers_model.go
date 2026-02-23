package models

type Headers struct {
	ContentType string `json:"Content-Type"`
	Authorization string `json:"Authorization"`
	AppToken string `json:"App-Token"`
	SessionToken string `json:"Session-Token"`
}