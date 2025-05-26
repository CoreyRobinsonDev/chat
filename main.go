package main

import (
	"os"

	"github.com/coreyrobinsondev/search/settings"
	"github.com/coreyrobinsondev/search/ui"
	u "github.com/coreyrobinsondev/utils"
	"google.golang.org/genai"
)


func main() {
	u.SetErrorHandler(func (err error) {
		settings.Logger.Fatal(err)
	})	
	settings.ConfigFile.Init()
	settings.ConfigFile.GeminiChatHistory = []*genai.Content{}
	settings.ConfigFile.Write()

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
