package models

type Training struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Questions   any    `json:"questions"`
}

type CreateTrainingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Questions   any    `json:"questions"`
}
