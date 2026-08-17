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
	Answer   any      `json:"answer,omitempty"`
}

type Training struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	RequiredCorrect int        `json:"requiredCorrect"`
	ShowAnswers     bool       `json:"showAnswers"`
	Questions       []Question `json:"questions"`
}

type CreateTrainingRequest struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	RequiredCorrect int        `json:"requiredCorrect"`
	ShowAnswers     bool       `json:"showAnswers"`
	Questions       []Question `json:"questions"`
}

type UserTraining struct {
	UserUUID    string    `json:"userUuid"`
	TrainingID  int       `json:"trainingId"`
	CompletedAt time.Time `json:"completedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type Submission = map[int]string

type SubmissionResponse struct {
	Passed         bool         `json:"passed"`
	NumCorrect     int          `json:"numCorrect"`
	NumIncorrect   int          `json:"numIncorrect"`
	Grade          int          `json:"grade"`
	GradedResponse map[int]bool `json:"gradedResponse"`
}
