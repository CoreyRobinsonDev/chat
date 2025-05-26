package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coreyrobinsondev/search/settings"
	u "github.com/coreyrobinsondev/utils"
)

const gap = "\n\n"
var Chat ChatModel

func RunChat() {
	Chat = initialModel()
	p := tea.NewProgram(Chat)

	u.Unwrap(p.Run())
}

type (
	errMsg error
)

type ChatModel struct {
	Viewport    viewport.Model
	Messages    []string
	Textarea    textarea.Model
	SenderStyle lipgloss.Style
	Err         error
}

func initialModel() ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Ask AI..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("Active model: %s", settings.ConfigFile.Model)))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		Textarea:    ta,
		Messages:    []string{},
		Viewport:    vp,
		SenderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		Err:         nil,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.Textarea, tiCmd = m.Textarea.Update(msg)
	m.Viewport, vpCmd = m.Viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Textarea.SetWidth(msg.Width)
		m.Viewport.Height = msg.Height - m.Textarea.Height() - lipgloss.Height(gap)

		if len(m.Messages) > 0 {
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(strings.Join(m.Messages, "\n")))
		}
		m.Viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.Textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.Messages = append(m.Messages, m.SenderStyle.Render("You: ")+m.Textarea.Value())
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(strings.Join(m.Messages, "\n")))
			m.Textarea.Reset()
			m.Viewport.GotoBottom()
		}

	case errMsg:
		m.Err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.Viewport.View(),
		gap,
		m.Textarea.View(),
	)
}
