package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	init    bool
	results []string
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
			m.results = m.results[:len(m.results)-1]
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
	log.Println("view")
	if !m.init {
		return "Initialising..." + fmt.Sprint(m.init, m.results)
	}
	return strings.Join(m.results, "\n")
}

func (m *model) streamResult() {
	for result := range tokeniseStream(clipStream(Freq)) {
		if result.Error != nil {
			continue
		}
		tr := result.Value
		line := fmt.Sprintf("Token: %4d, Words: %4d, Chars: %5d, Signature: %q", tr.Tokens, tr.Words, tr.Chars, tr.Sign)
		m.results = append(m.results, line)
	}
}

func main() {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	model := new(model)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
