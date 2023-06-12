package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pandodao/tokenizer-go"
)

type (
	validated struct{}
	TickMsg   struct {
		Id int
	}
)

type model struct {
	init     bool
	results  []TokenResult
	paused   atomic.Bool
	pipeline *pipeline
	table    table.Model
	columns  []table.Column
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.validate(), tick(0, Freq))
}

func (m *model) validate() tea.Cmd {
	return func() tea.Msg {
		_, err := tokenizer.CalToken("Init Goja runtime")
		if err != nil {
			log.Println("Cannot initialise tokenizer", err)
			os.Exit(1)
		}
		m.init = true
		return validated{}
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w1 := Min(msg.Width/10, 8)
		w2 := msg.Width - 3*w1
		SignLen = w2 * 8 / 10

		for i, w := range []int{w1, w1, w1, w2} {
			m.columns[i].Width = w
		}
		for i, r := range m.table.Rows() {
			r[3] = GenerateSign(m.results[i].Clip)
		}
		m.table.UpdateViewport()
		log.Println("Window resize", m.columns)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c", "C":
			// Delete all rows
			m.results = m.results[:0]
			rows := m.table.Rows()
			rows = rows[:0]
			m.table.SetRows(rows)
		case "x", "X":
			// Delete row at the current position
			if len(m.table.Rows()) == 0 {
				return m, nil
			}
			current := m.table.Cursor()
			m.results = append(m.results[:current], m.results[current+1:]...)
			rows := m.table.Rows()
			rows = append(rows[:current], rows[current+1:]...)
			m.table.SetRows(rows)
			if current == len(rows) {
				m.table.SetCursor(current - 1)
			}
		case "J":
			// Move current row down
			current := m.table.Cursor()
			rows := m.table.Rows()
			if current == len(rows)-1 {
				return m, nil
			}
			rows[current+1], rows[current] = rows[current], rows[current+1]
			m.results[current+1], m.results[current] = m.results[current], m.results[current+1]
			m.table.SetCursor(current + 1)
		case "K":
			// Move current row up
			current := m.table.Cursor()
			if current == 0 {
				return m, nil
			}
			rows := m.table.Rows()
			rows[current-1], rows[current] = rows[current], rows[current-1]
			m.results[current-1], m.results[current] = m.results[current], m.results[current-1]
			m.table.SetCursor(current - 1)
		case "s", "S":
			// Enable/disable reading the clipStream
			m.pipeline.pause <- struct{}{}
			m.paused.Store(!m.paused.Load())
		}
	case validated:
		go m.streamResult()
	case TickMsg:
		return m, tick(msg.Id+1, Freq)
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func tick(id int, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return TickMsg{Id: id}
	})
}

func (m *model) View() string {
	if !m.init {
		return "Initialising..." + fmt.Sprint(m.init, m.results)
	}

	builder := strings.Builder{}
	builder.Grow(len(m.results) * 100)
	// builder.WriteString("\n\n")
	// for _, res := range m.results {
	// 	builder.WriteKeyValue("Token", "%4d", res.Tokens)
	// 	builder.WriteKeyValue("Words", "%4d", res.Words)
	// 	builder.WriteKeyValue("Chars", "%5d", res.Chars)
	// 	builder.WriteKeyValue("Sign", "%q", res.Sign)
	// 	builder.LineBreak()
	// }
	builder.WriteString(m.table.View())
	if m.paused.Load() {
		builder.WriteString("\nPaused...\n")
	}
	return builder.String()
}

func (m *model) streamResult() {
	for result := range m.pipeline.tokeniseStream(m.pipeline.clipStream()) {
		if result.Error != nil {
			continue
		}
		m.results = append(m.results, result.Value)
		m.table.SetRows(append(m.table.Rows(), result.Value.ToTableRow()))
		m.table.UpdateViewport()
	}
}

func main() {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	initFunc()
	columns := []table.Column{
		{Title: "Tokens", Width: 10},
		{Title: "Words", Width: 10},
		{Title: "Chars", Width: 10},
		{Title: "Signature", Width: 70},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	model := &model{
		pipeline: NewPipeline(Freq),
		table:    t,
		columns:  columns,
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
