package view

import "time"

type Video struct {
	ID          string    `json:"id"`
	Platform    string    `json:"platform"`
	BloggerURL  string    `json:"blogger_url"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}
