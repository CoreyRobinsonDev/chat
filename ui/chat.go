package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/coreyrobinsondev/chat/ai"
	"github.com/coreyrobinsondev/chat/settings"
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
	Input		chan string 
	Sub 		chan struct{}
	AiResponse 	chan string
	Viewport    viewport.Model
	Spinner		spinner.Model
	Messages    []string
	Textarea    textarea.Model
	SenderStyle lipgloss.Style
	Err         error
}

func initialModel() ChatModel {
	in := make(chan string)
	sub := make(chan struct{})
	res := make(chan string)

	sp := spinner.New()
	sp.Spinner.Frames = []string{""}
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("247"))


	ta := textarea.New()
	ta.Placeholder = "Ask AI..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)


	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf("Active model: %s", settings.ConfigFile.Model)))
	vp.KeyMap.PageDown = key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	)
	vp.KeyMap.PageUp = key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	)
	vp.KeyMap.HalfPageDown = key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "½ page down"),
	)
	vp.KeyMap.HalfPageUp = key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "½ page up"),
	)
	vp.KeyMap.Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)
	vp.KeyMap.Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		Input: 		 in,
		Sub:		 sub,
		AiResponse:  res,
		Textarea:    ta,
		Spinner:     sp,
		Messages:    []string{},
		Viewport:    vp,
		SenderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		Err:         nil,
	}
}

type responseMsg struct {}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func (m ChatModel) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink, 
		m.Spinner.Tick, 
		ai.RunGemini(m.Input, m.Sub, m.AiResponse),
		waitForActivity(m.Sub),
	)
		
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)

	m.Textarea, tiCmd = m.Textarea.Update(msg)
	m.Viewport, vpCmd = m.Viewport.Update(msg)
	m.Spinner, spCmd = m.Spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width
		m.Textarea.SetWidth(msg.Width)
		m.Viewport.Height = msg.Height - m.Textarea.Height() - lipgloss.Height("\n")

		if len(m.Messages) > 0 {
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(strings.Join(m.Messages, "\n")))
		}
		m.Viewport.GotoBottom()
	case responseMsg:
			// to prevent (error) from being rendered
			// the length of this array must be > whatever spinner was just running
			m.Spinner.Spinner.Frames = []string{
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
				lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
			}
			m.Messages = append(m.Messages, lipgloss.NewStyle().Foreground(lipgloss.Color("38")).Render("Gemini: ")+<-m.AiResponse)
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(strings.Join(m.Messages, "\n")))
			m.Viewport.GotoBottom()
		return m, waitForActivity(m.Sub)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.Messages = append(m.Messages, m.SenderStyle.Render("You: ")+m.Textarea.Value())
			m.Spinner.Spinner = spinner.Line
			m.Input <- m.Textarea.Value()
			m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(strings.Join(m.Messages, "\n")))
			m.Textarea.Reset()
			m.Viewport.GotoBottom()
		}

	case errMsg:
		m.Err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m ChatModel) View() string {
	return fmt.Sprintf(
		"%s\n%s%s%s%s%s%s%s%s\n%s",
		m.Viewport.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().TopLeft),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		lipgloss.NewStyle().Foreground(lipgloss.Color("247")).Render(lipgloss.ThickBorder().Top),
		m.Spinner.View(),
		m.Textarea.View(),
	)
}
