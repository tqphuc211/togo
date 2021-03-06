package sqllite

import (
	"context"
	"database/sql"
	"time"

	"github.com/manabie-com/togo/internal/storages"
)

// LiteDB for working with sqllite
type LiteDB struct {
	DB        *sql.DB
	DriveName string
}

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = ? AND created_date = ?`
	if l.DriveName == "postgres" {
		stmt = `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = $1 AND created_date = $2`
	}
	rows, err := l.DB.QueryContext(ctx, stmt, userID, createdDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask adds a new task to DB
func (l *LiteDB) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES (?, ?, ?, ?)`
	if l.DriveName == "postgres" {
		stmt = `INSERT INTO tasks (id, content, user_id, created_date) VALUES ($1, $2, $3, $4)`
	}
	_, err := l.DB.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUser returns tasks if match userID AND password
func (l *LiteDB) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := `SELECT id FROM users WHERE id = ? AND password = ?`
	if l.DriveName == "postgres" {
		stmt = `SELECT id FROM users WHERE id = $1 AND password = $2`
	}
	row := l.DB.QueryRowContext(ctx, stmt, userID, pwd)
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		//log.Println(err)
		return false
	}

	return true
}

// Get user limit task
func (l *LiteDB) GetUserLimitTask(ctx context.Context, userID string) (int, error) {
	stmt := `SELECT max_todo FROM users WHERE id = ?`
	if l.DriveName == "postgres" {
		stmt = `SELECT max_todo FROM users WHERE id = $1`
	}
	row := l.DB.QueryRowContext(ctx, stmt, userID)

	var maxTodo int
	err := row.Scan(&maxTodo)

	if err != nil {
		return 0, err
	}
	return maxTodo, nil
}

// Validate Limit task
func (l *LiteDB) GetUserTaskToday(ctx context.Context, userID string) (int, error) {

	today := time.Now().Format("2006-01-02")
	stmt := `SELECT count(*) FROM tasks WHERE user_id = ? AND created_date = ?`
	if l.DriveName == "postgres" {
		stmt = `SELECT count(*) FROM tasks WHERE user_id = $1 AND created_date = $2`
	}
	row := l.DB.QueryRowContext(ctx, stmt, userID, today)

	var countTaskToday int
	err := row.Scan(&countTaskToday)

	if err != nil {
		return 0, err
	}
	return countTaskToday, nil
}
