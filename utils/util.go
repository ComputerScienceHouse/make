package utils

import (
	"errors"
	"makedotcsh/database"
	"makedotcsh/models"
	"makedotcsh/worker"
	"time"

	csh_auth "github.com/computersciencehouse/csh-auth/v2"
	"github.com/gin-gonic/gin"
)

func GetCSHAuth(c *gin.Context) (csh_auth.UserInfo, error) {
	authAny, exists := c.Get("cshauth")
	if !exists {
		err := errors.New("cshauth does not exist in context")
		return csh_auth.UserInfo{}, err
	}

	auth, ok := authAny.(*csh_auth.Claims)
	if !ok {
		err := errors.New("authentication data in gin context does not match structure")
		return csh_auth.UserInfo{}, err
	}

	user := csh_auth.UserInfo{
		Uuid:     auth.Uuid,
		Email:    auth.Email,
		Username: auth.Username,
		FullName: auth.FullName,
		Groups:   auth.Groups,
	}

	return user, nil
}

func GradeTraining(trainingId int, userUUID string, submission models.Submission) (models.SubmissionResponse, error) {
	training, err := database.Helper.GetTraining(trainingId, true)
	if err != nil {
		return models.SubmissionResponse{}, err
	}

	var answers map[int]models.Question = map[int]models.Question{}

	var correct, incorrect int
	graded := map[int]bool{}
	totalQuestions := 0

	for _, q := range training.Questions {
		answers[q.ID] = q

		if !q.Required {
			continue
		}

		_, exists := submission[q.ID]
		if !exists {
			return models.SubmissionResponse{}, errors.New("invalid submission")
		}

		totalQuestions++
	}

	if len(answers) == 0 {
		// no submission allowed
		return models.SubmissionResponse{}, errors.New("invalid submission")
	}

	// grade answers
	for id, r := range submission {
		if !answers[id].Required {
			continue
		}

		var isCorrect bool

		isCorrect = answers[id].Answer == r

		graded[id] = isCorrect

		if isCorrect {
			correct++
		} else {
			incorrect++
		}
	}

	passed := correct >= training.RequiredCorrect

	res := models.SubmissionResponse{
		Passed:       passed,
		NumCorrect:   correct,
		NumIncorrect: incorrect,
		Grade:        int((float64(correct) / float64(totalQuestions)) * 100),
	}

	if training.ShowAnswers {
		res.GradedResponse = graded
	} else {
		res.GradedResponse = map[int]bool{}
	}

	if passed {
		SubmitTraining(userUUID, training)

		// trigger worker to update ldap groups
		worker.TriggerWorker()
	}

	return res, nil
}

func SubmitTraining(uuid string, training models.Training) error {
	//TODO: this needs to be refined when
	// we formally define when a saftey seminar should
	// expire after completion
	now := time.Now()
	future := now.AddDate(1, 0, 0)

	newUserTraining := models.UserTraining{
		UserUUID:    uuid,
		TrainingID:  training.ID,
		CompletedAt: now,
		ExpiresAt:   future,
	}

	err := database.Helper.CreateUserTraining(newUserTraining)
	if err != nil {
		return err
	}

	return nil
}
