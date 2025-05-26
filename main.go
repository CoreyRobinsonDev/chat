package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	u "github.com/coreyrobinsondev/utils"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	TimeFormat: time.TimeOnly,
})

type Model string
const (
	GEMINI_25_FLASH_PREVIEW_0520 Model = "gemini-2.5-flash-preview-05-20"
	GEMINI_20_FLASH Model = "gemini-2.0-flash"
	GEMINI_20_FLASH_LITE Model = "gemini-2.0-flash-lite"
)

var config Config


func main() {
	u.SetErrorHandler(func (err error) {
		logger.Fatal(err)
	})	
	config.Init()

	args := os.Args[1:]

	if len(args) != 0 {
		switch args[0] {
		case "config":
			items := []list.Item{
				item(string(GEMINI_25_FLASH_PREVIEW_0520)),
				item(string(GEMINI_20_FLASH)),
				item(string(GEMINI_20_FLASH_LITE)),
			}

			const defaultWidth = 20

			l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
			l.Title = fmt.Sprintf("Active AI model: %s\nSwitching to...", config.Model)
			l.SetShowStatusBar(false)
			l.SetFilteringEnabled(false)
			l.Styles.Title = titleStyle
			l.Styles.PaginationStyle = paginationStyle
			l.Styles.HelpStyle = helpStyle

			m := model{list: l}

			u.Unwrap(tea.NewProgram(m).Run())
			os.Exit(0)
		}
	}
	RunGemini()
}


