package model

type Book struct {
	ID        string `json:"id,omitempty"`
	AuthorsID string `json:"authors_id,omitempty"`
	Title     string `json:"title"`
	Genre     string `json:"genre"`
	ISBN      string `json:"isbn"`
	Author    Author `json:"author,omitempty"`
}
