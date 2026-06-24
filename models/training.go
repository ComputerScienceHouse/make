package models

import (
	"time"
)

type Question struct {
	ID       int      `json:"id"`
	Label    string   `json:"label"`
	Required bool     `json:"required,omitempty"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
}

type QuestionWithAnswer struct {
	Question
	Answer any `json:"answer"`
}
type Training struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type CreateTrainingRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Questions   []QuestionWithAnswer `json:"questions"`
}

type UserTraining struct {
	UserUUID    string    `json:"userUuid"`
	TrainingID  int       `json:"trainingId"`
	CompletedAt time.Time `json:"completedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
