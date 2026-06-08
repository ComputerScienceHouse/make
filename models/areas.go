package models

type Area struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoURL    string `json:"photourl"`
}

type CreateAreaRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoURL    string `json:"photourl"`
}
