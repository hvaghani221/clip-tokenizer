package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "c", "C":
			m.results = m.results[:0]
		case "x", "X":
			if len(m.results) > 0 {
				m.results = m.results[:len(m.results)-1]
			}
		case "s", "S":
			m.pipeline.pause <- struct{}{}
			m.paused.Store(!m.paused.Load())
		}
	case validated:
		go m.streamResult()
	case TickMsg:
		return m, tick(msg.Id+1, Freq)
	}
	return m, nil
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
	builder := NewColorBuilder()
	builder.Grow(len(m.results) * 100)
	for _, res := range m.results {
		builder.WriteKeyValue("Token", "%4d", res.Tokens)
		builder.WriteKeyValue("Words", "%4d", res.Words)
		builder.WriteKeyValue("Chars", "%5d", res.Chars)
		builder.WriteKeyValue("Sign", "%q", res.Sign)
		builder.LineBreak()
	}
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
	}
}

func main() {
	initFunc()
	model := &model{
		pipeline: NewPipeline(Freq),
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
