package models

import (
	"time"
)

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

type UserTraining struct {
	UserUUID    string    `json:"userUuid"`
	TrainingID  int       `json:"trainingId"`
	CompletedAt time.Time `json:"completedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
