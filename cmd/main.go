package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabriel-kimutai/remi/internal/database"
	"github.com/gabriel-kimutai/remi/models"
)

const (
	// Subcommands
	add  = "add"
	list = "list"

	// Durations
	secs = 's'
	mins = 'm'
	hrs  = 'h'
	day  = 'd'
)

type Main struct {
}

type App struct {
	repo database.Repository
}

func New() *App {
	return &App{
		repo: database.New(),
	}

}

func init() {
	app := New()
	if err := app.repo.CreateTables(); err != nil {
		log.Panicf("failed to create tables: %e\n", err)
	}

}

func main() {
	app := New()
	defer app.repo.Close()

	p := tea.NewProgram(models.NewMainModel(app.repo), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	} else {
		return
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println(`
Usage: remi <subcommand> <content> <time>
 Example
 remi add foo 1d`)
		return
	}

	subcommand := args[1]

	switch subcommand {
	case add:
		content := args[2]
		at := args[3]
		t, err := time.ParseDuration(at)
		if err != nil {
			log.Fatalf("failed to parse duration: %e", err)
		}
		if err := app.repo.InsertReminder(content, t); err != nil {
			slog.Warn(fmt.Sprintf("failed to insert reminder: %e", err))
			return
		}
	case list:
		reminders, err := app.repo.SelectAllReminders()
		if err != nil {
			log.Fatalf("failed to fetch reminders: %e", err)
		}
		for _, reminder := range reminders {
			fmt.Printf("%d: %s\n", reminder.ID, reminder.Content)
		}
	}

}
