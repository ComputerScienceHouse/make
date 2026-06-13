package database

import (
	"database/sql"
	"log"
	"makedotcsh/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type DatabaseHelper struct {
	DB *sql.DB
}

var Helper *DatabaseHelper

func Init() {
	var err error

	DB, err = sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = Migrate()
	if err != nil {
		log.Fatal(err)
	}

	Helper = &DatabaseHelper{DB: DB}
}

func (database *DatabaseHelper) GetArea(areaID int) (models.Area, error) {
	row := database.DB.QueryRow("SELECT id, name, description, ldapgroup, photourl FROM areas WHERE id = ?", areaID)

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
		"INSERT INTO areas (name, description, ldapgroup, photourl) VALUES (?, ?, ?, ?)",
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
		"UPDATE areas SET name = ?, description = ?, ldapgroup = ?, photourl = ? WHERE id = ?",
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

func (database *DatabaseHelper) GetTraining(trainingID int) (models.Training, error) {
	row := database.DB.QueryRow("SELECT id, title, description, questions FROM trainings WHERE id = ?", trainingID)

	var training models.Training

	err := row.Scan(
		&training.ID,
		&training.Title,
		&training.Description,
		&training.Questions,
	)

	if err != nil {
		return models.Training{}, err
	}

	return training, nil
}

func (database *DatabaseHelper) GetAllTrainings() ([]models.Training, error) {
	rows, err := database.DB.Query("SELECT id, title, description, questions FROM trainings")
	if err != nil {
		return []models.Training{}, nil
	}

	var trainings []models.Training

	for rows.Next() {
		var training models.Training

		err := rows.Scan(
			&training.ID,
			&training.Title,
			&training.Description,
			&training.Questions,
		)

		if err != nil {
			return []models.Training{}, nil
		}

		trainings = append(trainings, training)
	}
	if err != nil {
		return []models.Training{}, err
	}

	return trainings, nil
}

func (database *DatabaseHelper) CreateTraining(training models.CreateTrainingRequest) error {
	_, err := database.DB.Exec(
		"INSERT INTO trainings (title, description, questions) VALUES (?, ?, ?)",
		training.Title,
		training.Description,
		training.Questions,
	)

	if err != nil {
		return err
	}

	return nil
}

func (database *DatabaseHelper) UpdateTraining(training models.CreateTrainingRequest, trainingId int) error {
	_, err := database.DB.Exec(
		"UPDATE trainings SET title = ?, description = ?, questions = ? WHERE id = ?",
		training.Title,
		training.Description,
		training.Questions,
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
	rows, err := database.DB.Query("SELECT training_id FROM area_trainings WHERE area_id = ?", area)
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
		training, err := database.GetTraining(trainingID)
		if err != nil {
			return []models.Training{}, err
		}

		trainings = append(trainings, training)
	}

	return trainings, nil
}

func (database *DatabaseHelper) CreateUserTraining(training models.UserTraining) error {
	_, err := database.DB.Exec(
		"INSERT INTO user_trainings (user_uuid, training_id, completed_at, expires_at) VALUES (?, ?, ?, ?)",
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

func (database *DatabaseHelper) GetUserTrainings(uuid string) ([]int, error) {
	rows, err := database.DB.Query("SELECT user_uuid, training_id, completed_at, expires_at FROM user_trainings WHERE user_uuid = ?", uuid)
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

func (database *DatabaseHelper) AddTrainingsToArea(areaID int, trainingIDs []int) error {
	for _, trainingID := range trainingIDs {
		_, err := database.DB.Exec(
			"INSERT INTO area_trainings (area_id, training_id) VALUES (?, ?)",
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
		"DELETE FROM area_trainings WHERE area_id = ? AND training_id = ?",
		areaID, trainingID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (database *DatabaseHelper) GetAllAreasWithUserAccess(uuid string) ([]int, error) {
	// insane query god bless sql (this was 100 lines of logic beforehand)
	query := "SELECT a.id FROM areas a JOIN area_trainings at ON at.area_id = a.id LEFT JOIN user_trainings ut ON ut.training_id = at.training_id AND ut.user_uuid = ? GROUP BY a.id, a.name HAVING COUNT(DISTINCT at.training_id) = COUNT(DISTINCT ut.training_id)"
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
		"DELETE FROM trainings WHERE id = ?",
		trainingID,
	)

	if err != nil {
		return err
	}

	return nil
}
