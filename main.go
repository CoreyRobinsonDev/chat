package main

import (
	"os"

	"github.com/coreyrobinsondev/chat/settings"
	"github.com/coreyrobinsondev/chat/ui"
	u "github.com/coreyrobinsondev/utils"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func main() {
	u.SetErrorHandler(func(err error) {
		settings.Logger.Fatal(err)
	})
	settings.ConfigFile.Init()
	settings.ConfigFile.GeminiChatHistory = []*genai.Content{}
	settings.ConfigFile.Write()
	if len(settings.ConfigFile.GeminiApiKey) == 0 {
		godotenv.Load()
		settings.ConfigFile.GeminiApiKey = os.Getenv("GEMINI_API_KEY")
		settings.ConfigFile.Write()
		if len(settings.ConfigFile.GeminiApiKey) == 0 {
			settings.Logger.Fatal("Please provide your API key to 'geminiApiKey' in ~/.config/chat/settings.json")
		}
	}

	args := os.Args[1:]

	if len(args) != 0 {
		switch args[0] {
		case "config":
			ui.RunList()
			os.Exit(0)
		}
	}
	ui.RunChat()
}
