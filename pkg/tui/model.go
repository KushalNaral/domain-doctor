package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhanushnehru/domain-doctor/pkg/health"
)

type state int

const (
	stateInput state = iota
	stateChecking
	stateDone
)

type checkDoneMsg struct{ result *health.DomainHealth }
type checkErrMsg struct{ err error }

type Model struct {
	state   state
	input   textinput.Model
	spinner spinner.Model
	result  *health.DomainHealth
	err     error
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "e.g. google.com"
	ti.Focus()
	ti.CharLimit = 253

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))

	return Model{state: stateInput, input: ti, spinner: sp}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == stateInput {
				domain := strings.TrimSpace(m.input.Value())
				if domain == "" {
					return m, nil
				}
				m.state = stateChecking
				return m, tea.Batch(m.spinner.Tick, runCheck(domain))
			}
			if m.state == stateDone {
				return m, tea.Quit
			}
		}

	case checkDoneMsg:
		m.state = stateDone
		m.result = msg.result
		return m, nil

	case checkErrMsg:
		m.err = msg.err
		m.state = stateDone
		return m, nil

	case spinner.TickMsg:
		if m.state == stateChecking {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if m.state == stateInput {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateInput:
		return inputView(m)
	case stateChecking:
		return checkingView(m)
	case stateDone:
		return doneView(m)
	}
	return ""
}

func inputView(m Model) string {
	content := promptTitleStyle.Render("🩺 Domain Doctor") + "\n\n" +
		promptLabelStyle.Render("Enter a domain to analyze:") + "\n" +
		m.input.View()
	return "\n" + promptCardStyle.Render(content) + "\n"
}

func checkingView(m Model) string {
	return "\n  " + m.spinner.View() + " " +
		spinnerLabelStyle.Render("Checking "+m.input.Value()+"...") + "\n"
}

func doneView(m Model) string {
	if m.err != nil {
		return "\n  Error: " + m.err.Error() + "\n"
	}
	return "\n" + m.result.RenderReport(true)
}

func runCheck(domain string) tea.Cmd {
	return func() tea.Msg {
		result := health.Check(domain)
		return checkDoneMsg{result: result}
	}
}
