package ai

import (
	"context"

	"github.com/coreyrobinsondev/search/settings"
	u "github.com/coreyrobinsondev/utils"
	"google.golang.org/genai"
)

func RunGemini(input string) string {
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
	result := u.Unwrap(chat.SendMessage(ctx, genai.Part{Text: input}))

	if len(result.Candidates) > 0 {
		return result.Candidates[0].Content.Parts[0].Text
	} else { return "" } 
}
