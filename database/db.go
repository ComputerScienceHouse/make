package database

import (
	"database/sql"
	"fmt"
	"log"
	"makedotcsh/models"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type DatabaseHelper struct {
	DB *sql.DB
}

var Helper *DatabaseHelper

func Init() {
	var DB *sql.DB

	host := os.Getenv("MAKE_DB_HOST")
	name := os.Getenv("MAKE_DB_NAME")
	user := os.Getenv("MAKE_DB_USER")
	pass := os.Getenv("MAKE_DB_PASS")
	dbOptions := fmt.Sprintf("postgres://%s:%s@%s/%s", url.QueryEscape(user), url.QueryEscape(pass), host, url.PathEscape(name))

	DB, err := sql.Open("postgres", dbOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = Migrate(DB)
	if err != nil {
		log.Fatal(err)
	}

	Helper = &DatabaseHelper{DB: DB}
}

func (database *DatabaseHelper) GetArea(areaID int) (models.Area, error) {
	row := database.DB.QueryRow("SELECT id, name, description, ldapgroup, photourl FROM areas WHERE id = $1", areaID)

	var area models.Area

	err := row.Scan(
		&area.ID,
		&area.Name,
		&area.Description,
		&area.LdapGroup,
		&area.PhotoURL,
	)

	if err != nil {
		return models.Area{}, err
	}

	return area, nil
}

func (database *DatabaseHelper) GetAllAreas() ([]models.Area, error) {
	rows, err := database.DB.Query("SELECT id, name, description, ldapgroup, photourl FROM areas")
	if err != nil {
		return []models.Area{}, nil
	}

	var areas []models.Area = []models.Area{}

	for rows.Next() {
		var area models.Area

		err := rows.Scan(
			&area.ID,
			&area.Name,
			&area.Description,
			&area.LdapGroup,
			&area.PhotoURL,
		)

		if err != nil {
			return []models.Area{}, nil
		}

		areas = append(areas, area)
	}

	return areas, nil
}

func (database *DatabaseHelper) CreateArea(area models.CreateAreaRequest) error {
	_, err := database.DB.Exec(
		"INSERT INTO areas (name, description, ldapgroup, photourl) VALUES ($1, $2, $3, $4)",
		area.Name,
		area.Description,
		area.LdapGroup,
		area.PhotoURL,
	)

	if err != nil {
		return err
	}

	return nil
}

func (database *DatabaseHelper) UpdateArea(area models.CreateAreaRequest, id int) error {
	_, err := database.DB.Exec(
		"UPDATE areas SET name = $1, description = $2, ldapgroup = $3, photourl = $4 WHERE id = $5",
		area.Name,
		area.Description,
		area.LdapGroup,
		area.PhotoURL,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (database *DatabaseHelper) GetTraining(trainingID int, includeAnswers bool) (models.Training, error) {
	query := `
	SELECT
		t.id, t.title, t.description,
		q.id, q.label, q.type, q.answer, q.required,
		o.label
	FROM trainings t
	LEFT JOIN training_questions q ON q.training_id = t.id
	LEFT JOIN question_options o ON o.question_id = q.id
	WHERE t.id = $1
	ORDER BY q.id, o.id
	`

	rows, err := database.DB.Query(query, trainingID)
	if err != nil {
		return models.Training{}, err
	}

	var training models.Training
	var currentQuestion *models.Question
	lastQuestionId := 0

	for rows.Next() {
		var question models.Question
		var option sql.NullString

		err = rows.Scan(
			&training.ID,
			&training.Title,
			&training.Description,

			&question.ID,
			&question.Label,
			&question.Type,
			&question.Answer,
			&question.Required,

			&option,
		)

		if err != nil {
			return models.Training{}, err
		}

		if lastQuestionId != question.ID { // new question
			training.Questions = append(training.Questions, question)
			currentQuestion = &training.Questions[len(training.Questions)-1]
			lastQuestionId = question.ID

			if !includeAnswers {
				currentQuestion.Answer = nil
			}
		}

		// radio type questions will have additional options
		if question.Type == "radio" && option.Valid {
			// is an option
			currentQuestion.Options = append(currentQuestion.Options, option.String)
		}
	}

	return training, nil
}

func (database *DatabaseHelper) GetAllTrainings() ([]models.Training, error) {
	rows, err := database.DB.Query("SELECT id FROM trainings")
	if err != nil {
		return []models.Training{}, err
	}

	var trainings []models.Training = []models.Training{}

	for rows.Next() {
		var trainingId int

		err := rows.Scan(
			&trainingId,
		)

		if err != nil {
			return []models.Training{}, err
		}

		training, err := database.GetTraining(trainingId, false)
		if err != nil {
			return []models.Training{}, err
		}

		trainings = append(trainings, training)
	}

	return trainings, nil
}

func (database *DatabaseHelper) CreateTraining(training models.CreateTrainingRequest) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var trainingId int
	err = tx.QueryRow(
		"INSERT INTO trainings (title, description) VALUES ($1, $2) RETURNING id",
		training.Title,
		training.Description,
	).Scan(&trainingId)

	if err != nil {
		return err
	}

	for _, q := range training.Questions {
		var questionId int
		err = tx.QueryRow(
			"INSERT INTO training_questions (training_id, label, type, answer, required) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			trainingId,
			q.Label,
			q.Type,
			q.Answer,
			q.Required,
		).Scan(&questionId)

		if err != nil {
			return err
		}

		if q.Type == "radio" {
			for _, option := range q.Options {
				_, err = tx.Exec(
					"INSERT INTO question_options (question_id, label) VALUES ($1, $2)",
					questionId,
					option,
				)

				if err != nil {
					return err
				}
			}

		}
	}

	return tx.Commit()
}

// TODO: make it actually update..
func (database *DatabaseHelper) UpdateTraining(training models.CreateTrainingRequest, trainingId int) error {
	_, err := database.DB.Exec(
		"UPDATE trainings SET title = $1, description = $2 WHERE id = $3",
		training.Title,
		training.Description,
		trainingId,
	)

	if err != nil {
		return err
	}

	return nil
}

type dbAreaTraining struct {
	areaID     int
	trainingID int
}

func (database *DatabaseHelper) GetAreaTrainings(area int) ([]models.Training, error) {
	rows, err := database.DB.Query("SELECT training_id FROM area_trainings WHERE area_id = $1", area)
	if err != nil {
		return []models.Training{}, nil
	}

	var trainings []models.Training = []models.Training{}

	for rows.Next() {
		var trainingID int

		err := rows.Scan(
			&trainingID,
		)

		if err != nil {
			return []models.Training{}, err
		}

		// fetch this training, add to array
		training, err := database.GetTraining(trainingID, false)
		if err != nil {
			return []models.Training{}, err
		}

		trainings = append(trainings, training)
	}

	return trainings, nil
}

func (database *DatabaseHelper) CreateUserTraining(training models.UserTraining) error {
	_, err := database.DB.Exec(
		"INSERT INTO user_trainings (user_uuid, training_id, completed_at, expires_at) VALUES ($1, $2, $3, $4)",
		training.UserUUID,
		training.TrainingID,
		training.CompletedAt,
		training.ExpiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (database *DatabaseHelper) GetCompletedUserTrainings(uuid string) ([]int, error) {
	rows, err := database.DB.Query("SELECT user_uuid, training_id, completed_at, expires_at FROM user_trainings WHERE user_uuid = $1", uuid)
	if err != nil {
		return []int{}, err
	}

	var trainings []int = []int{}

	for rows.Next() {
		var usertraining models.UserTraining

		err := rows.Scan(
			&usertraining.UserUUID,
			&usertraining.TrainingID,
			&usertraining.CompletedAt,
			&usertraining.ExpiresAt,
		)

		if err != nil {
			return []int{}, err
		}

		trainings = append(trainings, usertraining.TrainingID)
	}

	return trainings, nil
}

func (database *DatabaseHelper) GetUserTraining(trainingID int) (models.UserTraining, error) {
	row := database.DB.QueryRow("SELECT user_uuid, training_id, completed_at, expires_at FROM user_trainings WHERE training_id = $1", trainingID)

	var training models.UserTraining

	err := row.Scan(
		&training.UserUUID,
		&training.TrainingID,
		&training.CompletedAt,
		&training.ExpiresAt,
	)

	if err != nil {
		return models.UserTraining{}, err
	}

	return training, nil
}

func (database *DatabaseHelper) DeleteTrainingFromUser(trainingID int, uuid string) error {
	_, err := database.DB.Exec(
		"DELETE FROM user_trainings WHERE user_uuid = $1 AND training_id = $2",
		uuid,
		trainingID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (database *DatabaseHelper) AddTrainingsToArea(areaID int, trainingIDs []int) error {
	for _, trainingID := range trainingIDs {
		_, err := database.DB.Exec(
			"INSERT INTO area_trainings (area_id, training_id) VALUES ($1, $2)",
			areaID, trainingID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (database *DatabaseHelper) RemoveTrainingFromArea(areaID int, trainingID int) error {
	_, err := database.DB.Exec(
		"DELETE FROM area_trainings WHERE area_id = $1 AND training_id = $2",
		areaID, trainingID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (database *DatabaseHelper) GetAllAreasWithUserAccess(uuid string) ([]int, error) {
	// insane query god bless sql (this was 100 lines of logic beforehand)
	query := "SELECT a.id FROM areas a JOIN area_trainings at ON at.area_id = a.id LEFT JOIN user_trainings ut ON ut.training_id = at.training_id AND ut.user_uuid = $1 GROUP BY a.id, a.name HAVING COUNT(DISTINCT at.training_id) = COUNT(DISTINCT ut.training_id)"
	rows, err := database.DB.Query(query, uuid)
	if err != nil {
		return []int{}, err
	}

	areas := []int{}

	for rows.Next() {
		var area_id int

		err := rows.Scan(&area_id)

		if err != nil {
			return []int{}, err
		}

		areas = append(areas, area_id)
	}

	return areas, nil
}

func (database *DatabaseHelper) DeleteTraining(trainingID int) error {
	_, err := database.DB.Exec(
		"DELETE FROM trainings WHERE id = $1",
		trainingID,
	)

	if err != nil {
		return err
	}

	return nil
}

// [ { "id": 1, "label": "test", "required": true, "type": "radio", "options": [ "test", "tes2", "test3" ], "answer": "test" }, { "id": 2, "label": "answer is 5", "required": true, "type": "number", "answer": 5 } ]
