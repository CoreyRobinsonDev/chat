package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	u "github.com/coreyrobinsondev/utils"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportCaller: true,
	ReportTimestamp: true,
	TimeFormat: time.TimeOnly,
})

type Model string
const (
	GEMINI25FLASHPREVIEW0520 Model = "gemini-2.5-flash-preview-05-20"
	GEMINI20FLASH Model = "gemini-2.0-flash"
	GEMINI20FLASHLITE Model = "gemini-2.0-flash-lite"
)


func main() {
	u.SetErrorHandler(func (err error) {
		logger.Fatal(err)
	})	
	u.Expect(godotenv.Load())

	ctx := context.Background()
	client := u.Unwrap(genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	}))

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			"You're a senior software engineer giving short and concise answers. Include code examples", 
			genai.RoleUser,
		),
	}

	history := []*genai.Content{}

	chat := u.Unwrap(client.Chats.Create(ctx, string(GEMINI25FLASHPREVIEW0520), config, history))
	stream := chat.SendMessageStream(ctx, genai.Part{Text: "How to make a while loop in Go"})

	for chunk := range stream {
		part := chunk.Candidates[0].Content.Parts[0]
		fmt.Print(part.Text)
	}
}





