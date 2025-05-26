package main

import (
	"context"
	"fmt"

	u "github.com/coreyrobinsondev/utils"
	"google.golang.org/genai"
)

func RunGemini() {
	if len(config.GeminiApiKey) == 0 {
		logger.Fatal("Please provide your API key to 'geminiApiKey' in ~/.config/search/search.conf")
	}

	ctx := context.Background()
	client := u.Unwrap(genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	}))

	aiConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			"You're a senior software engineer giving short and concise answers. Include code examples", 
			genai.RoleUser,
			),
	}

	history := []*genai.Content{}

	chat := u.Unwrap(client.Chats.Create(ctx, config.Model, aiConfig, history))
	stream := chat.SendMessageStream(ctx, genai.Part{Text: "How to make a while loop in Go"})

	for chunk := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}
}
