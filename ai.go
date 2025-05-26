package main

import (
	"context"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/coreyrobinsondev/search/settings"
	"github.com/coreyrobinsondev/search/ui"
	u "github.com/coreyrobinsondev/utils"
	"google.golang.org/genai"
)

func RunGemini() {
	if len(settings.ConfigFile.GeminiApiKey) == 0 {
		settings.Logger.Fatal("Please provide your API key to 'geminiApiKey' in ~/.config/search/search.conf")
	}

	ctx := context.Background()
	client := u.Unwrap(genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: settings.ConfigFile.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	}))

	aiConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			"You're a senior software engineer giving short and concise answers. Include code examples", 
			genai.RoleUser,
			),
	}

	history := []*genai.Content{}

	chat := u.Unwrap(client.Chats.Create(ctx, settings.ConfigFile.Model, aiConfig, history))
	result := u.Unwrap(chat.SendMessage(ctx, genai.Part{Text: "How to make a while loop in Go"}))

	if len(result.Candidates) > 0 {
		part := result.Candidates[0].Content.Parts[0]
		ui.Chat.Messages = append(ui.Chat.Messages, ui.Chat.SenderStyle.Render("Gemini: ")+part.Text)
		ui.Chat.Viewport.SetContent(lipgloss.NewStyle().Width(ui.Chat.Viewport.Width).Render(strings.Join(ui.Chat.Messages, "\n")))
		ui.Chat.Textarea.Reset()
		ui.Chat.Viewport.GotoBottom()
	} else {
	}
}
