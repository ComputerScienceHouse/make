package models

type Resource struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
