package ai

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/coreyrobinsondev/search/settings"
	u "github.com/coreyrobinsondev/utils"
	"google.golang.org/genai"
)

func RunGemini(input chan string, sub chan struct {}, res chan string) tea.Cmd {

	ctx := context.Background()
	client := u.Unwrap(genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: settings.ConfigFile.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	}))

	aiConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			settings.ConfigFile.SystemInstruction, 
			genai.RoleUser,
		),
	}



	return func() tea.Msg {
		for {
			in := <- input
			chat := u.Unwrap(client.Chats.Create(ctx, settings.ConfigFile.Model, aiConfig, settings.ConfigFile.GeminiChatHistory))
			result := u.Unwrap(chat.SendMessage(ctx, genai.Part{Text: in}))

			settings.ConfigFile.GeminiChatHistory = append(
				settings.ConfigFile.GeminiChatHistory,
				genai.NewContentFromText(in, genai.RoleUser),
				result.Candidates[0].Content,
			)
			settings.ConfigFile.Write()
			sub <- struct{}{}
			res <- u.Unwrap(glamour.Render(result.Candidates[0].Content.Parts[0].Text, "dark"))
		}
	}
}
