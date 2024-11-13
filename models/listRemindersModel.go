package models

import (
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	db "github.com/gabriel-kimutai/remi/internal/database"
)

type ListRemindersModel struct {
	reminders        []db.Reminder
	cursor           int
	selectedReminder map[int]db.Reminder
	table            table.Model
}

var remindersTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialListRemindersModel(rp db.Repository) ListRemindersModel {
	reminders, err := rp.SelectAllReminders()
	if err != nil {
		log.Printf("failed to fetch reminders")
	}
	return ListRemindersModel{
		reminders: reminders,
	}
}

func (m ListRemindersModel) Init() tea.Cmd {
	return nil
}

func (m ListRemindersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case ctrl_c, q:
			return m, tea.Quit
		case up, k:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, cmd
		case down, j:
			if m.cursor < len(m.reminders)-1 {
				m.cursor++
			}
			return m, cmd
		case space:
			tea.Printf("current %d\n", m.table.Cursor())
		}
	}

	m.table, cmd = m.table.Update(cmd)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ListRemindersModel) View() string {
	var s strings.Builder
	lenReminders := len(m.reminders)
	if lenReminders < 1 {
		_, err := s.WriteString("no reminders")
		if err != nil {
			return err.Error()
		}
	}

	columns := []table.Column{
		{Title: "id", Width: 4},
		{Title: "content", Width: 20},
		{Title: "remind_at", Width: 20},
		{Title: "created_at", Width: 20},
	}

	rows := []table.Row{}
	for _, r := range m.reminders {
		row := table.Row{
			strconv.Itoa(int(r.ID)),
			r.Content,
			r.RemindAt.Format(db.HMN_LAYOUT),
			r.CreatedAt.Format(db.HMN_LAYOUT),
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	ts := table.DefaultStyles()
	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	ts.Selected = ts.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(ts)
	m.table = t

	return remindersTableStyle.Render(m.table.View()) + "\n " + m.table.HelpView()

}
