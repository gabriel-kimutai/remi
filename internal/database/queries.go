package database

import (
	"fmt"
	"time"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
	HMN_LAYOUT  = "Mon, 02 Jan 15:04"
)

type Reminder struct {
	ID          uint64     `json:"id" db:"id"`
	Content     string     `json:"content" db:"content"`
	IsRead      bool       `json:"is_read" db:"is_read"`
	RemindAt    *time.Time `json:"remind_at" db:"remind_at"`
	CreatedAt   *time.Time `json:"created_at" db:"created_at"`
	DismissedAt *time.Time `json:"dismissed_at" db:"dismissed_at"`
	ReadAt      *time.Time `json:"read_at" db:"read_at"`
}

func (r *repository) InsertReminder(content string, remindAt time.Duration) error {
	stmt := `insert into reminders(content, remind_at) values(?, ?)`
	now := time.Now()
	t := now.Add(remindAt)
	result, err := r.db.Exec(stmt, content, t)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("%d: %s at %s\n", lastId, content, t.Format(TIME_LAYOUT))

	return nil

}

func (r *repository) SelectAllReminders() ([]Reminder, error) {
	reminders := new([]Reminder)

	err := r.db.Select(reminders, "select * from reminders")
	if err != nil {
		return nil, err
	}
	return *reminders, nil
}
