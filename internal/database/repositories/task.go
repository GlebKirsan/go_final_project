package repositories

import (
	"database/sql"
	"strings"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/models"
)

const LIMIT = 50

type taskRepo struct {
	db         *sql.DB
	createStmt *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
	getStmt    *sql.Stmt
}

func NewTaskRepo(db *sql.DB) (*taskRepo, error) {
	repo := &taskRepo{db: db}

	var err error

	repo.createStmt, err = db.Prepare(create)
	if err != nil {
		return nil, err
	}

	repo.updateStmt, err = db.Prepare(update)
	if err != nil {
		return nil, err
	}

	repo.getStmt, err = db.Prepare(get)
	if err != nil {
		return nil, err
	}

	repo.deleteStmt, err = db.Prepare(delete)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *taskRepo) Create(task *models.Task) (int64, error) {
	res, err := repo.createStmt.Exec(sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (repo *taskRepo) Update(task *models.Task) error {
	_, err := repo.updateStmt.Exec(sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	return err
}

func (repo *taskRepo) Delete(id int64) error {
	_, err := repo.deleteStmt.Exec(sql.Named("id", id))
	return err
}

func (repo *taskRepo) Get(id int64) (*models.Task, error) {
	row := repo.getStmt.QueryRow(sql.Named("id", id))

	task := &models.Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (repo *taskRepo) GetAll(filter string) ([]models.Task, error) {
	var query strings.Builder
	var args []interface{}

	query.WriteString("SELECT id, date, title, comment, repeat FROM scheduler")
	if filter != "" {
		if date, err := time.Parse("02.01.2006", filter); err == nil {
			query.WriteString(" WHERE date = :date")
			args = append(args, sql.Named("date", date.Format("20060102")))
		} else {
			query.WriteString(" WHERE title LIKE CONCAT('%', :title, '%')")
			args = append(args, sql.Named("title", filter))
		}
	}

	query.WriteString(" LIMIT :limit")
	args = append(args, sql.Named("limit", LIMIT))

	rows, err := repo.db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (repo *taskRepo) Close() error {
	if err := repo.createStmt.Close(); err != nil {
		return err
	}
	if err := repo.updateStmt.Close(); err != nil {
		return err
	}
	if err := repo.getStmt.Close(); err != nil {
		return err
	}
	if err := repo.deleteStmt.Close(); err != nil {
		return err
	}
	return nil
}
