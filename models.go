package main

type Click struct {
	ID             int    `json:"id"`
	Code           string `json:"code"`
	DestinationURL string `json:"destination_url"`
	IP             string `json:"ip"`
	UserAgent      string `json:"user_agent"`
	Referer        string `json:"referer"`
	CreatedAt      string `json:"created_at"`
}

type ClickCount struct {
	Code           string `json:"code"`
	DestinationURL string `json:"destination_url"`
	Clicks         int    `json:"clicks"`
}
