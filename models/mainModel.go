package models

import (
	tea "github.com/charmbracelet/bubbletea"
	db "github.com/gabriel-kimutai/remi/internal/database"
)

type sessionState int

const (
	AddReminder sessionState = iota
	ListReminders
)

type MainModel struct {
	state         sessionState
	AddReminder   tea.Model
	ListReminders tea.Model
}

func NewMainModel(rp db.Repository) MainModel {
	return MainModel{
		state:         ListReminders,
		AddReminder:   nil,
		ListReminders: initialListRemindersModel(rp),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrl_c, q:
			return m, tea.Quit
		}
	}

	switch m.state {
	case ListReminders:
		return m.ListReminders.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case AddReminder:
		return m.AddReminder.View()
	case ListReminders:
		return m.ListReminders.View()
	default:
		return m.ListReminders.View()
	}
}
