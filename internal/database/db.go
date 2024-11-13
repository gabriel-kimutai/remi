package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	// Initialize the tables in the database
	CreateTables() error

	// Insert reminder
	InsertReminder(content string, remindAt time.Duration) error

	// Get all reminders
	SelectAllReminders() ([]Reminder, error)

	// CustomLogger
	Log(sql.Result)

	// Close DB
	Close() error
}

var repoInstance *repository

type repository struct {
	db *sqlx.DB
}

func New() Repository {
	if repoInstance != nil {
		return repoInstance
	}
	db, err := sqlx.Open("sqlite3", "./remi.db")
	if err != nil {
		log.Panicf("failed to open database: %e\n", err)
		return nil
	}
	repoInstance = &repository{
		db: db,
	}
	return repoInstance
}

func (r *repository) CreateTables() error {
	schema := `
        CREATE TABLE IF NOT EXISTS reminders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            is_read BOOLEAN DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            remind_at TIMESTAMP NOT NULL,
            dismissed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            read_at TIMESTAMP
        );
    `
	result, err := r.db.Exec(schema)
	if err != nil {
		return err
	}

	r.Log(result)

	return nil
}

func (r *repository) Log(result sql.Result) {

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return
	}
	if rowsAffected != 0 {
		slog.Info(fmt.Sprintf("[db]: %d rows affected", rowsAffected))
	}
	if lastInsertId != 0 {
		slog.Info(fmt.Sprintf("[db]: %d last insert ID", lastInsertId))
	}
}

func (r *repository) Close() error {
	if err := r.db.Close(); err != nil {
		return err
	}
	return nil
}
