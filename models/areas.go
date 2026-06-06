package models

type Area struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateAreaRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
