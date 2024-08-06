package repositories

import (
	"database/sql"

	"github.com/GlebKirsan/go-final-project/internal/models"
)

const LIMIT = 50

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (repo *TaskRepo) Create(task *models.Task) (int64, error) {
	res, err := repo.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (:date, :title, :comment, :repeat);`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (repo *TaskRepo) UpdateTask(task *models.Task) error {
	_, err := repo.db.Exec(`UPDATE scheduler
	SET date = :date,
	    title = :title,
		comment = :comment,
		repeat = :repeat
	WHERE id = :id;`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	return err
}

func (repo *TaskRepo) Delete(id int64) error {
	_, err := repo.db.Exec("DELETE FROM scheduler WHERE id = :id;", sql.Named("id", id))
	return err
}

func (repo *TaskRepo) GetTask(id int64) (*models.Task, error) {
	row := repo.db.QueryRow(`SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE id = :id;`, sql.Named("id", id))

	task := &models.Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func parseTasks(rows *sql.Rows) ([]models.Task, error) {
	tasks := make([]models.Task, 0, LIMIT)
	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TaskRepo) GetAll() ([]models.Task, error) {
	rows, err := repo.db.Query("SELECT id, date, title, comment, repeat FROM scheduler LIMIT :limit;", sql.Named("limit", LIMIT))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}

func (repo *TaskRepo) GetAllByDate(date string) ([]models.Task, error) {
	rows, err := repo.db.Query(`SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE date = @date
	LIMIT :limit;`, sql.Named("date", date), sql.Named("limit", LIMIT))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}

func (repo *TaskRepo) GetAllByTitle(title string) ([]models.Task, error) {
	rows, err := repo.db.Query(`SELECT id, date, title, comment, repeat 
	FROM scheduler
	WHERE title LIKE :title
	LIMIT :limit;`, sql.Named("title", "%"+title+"%"), sql.Named("limit", LIMIT))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}
