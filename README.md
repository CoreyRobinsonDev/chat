# search
Terminal based AI chat client running on Google Gemini models

![demo](https://vhs.charm.sh/vhs-4nLDenhbKScUPvy58sEDWt.gif)

## Installation
1. Build from source
    ```bash
    git clone https://github.com/CoreyRobinsonDev/search
    cd search
    go build
    ```
1. Add Gemini API key. This can be done in 2 ways:
    - Within the project directory create a `.env` file, add `GEMINI_API_KEY=your_api_key`, then run `go build`. This will be added to your settings file at `~/.config/search/settings.json`. This file will be created initial run
    - Alternitively, you can update the `geminiApiKey` field in your `~/.config/search/settings.json` file directly

## Commands
- `config`: Select the AI model to use
    - Other [Gemini models](https://ai.google.dev/gemini-api/docs/models) can be added to `ui/list`
    ```go
	items := []list.Item{
		item(string("gemini-2.5-flash-preview-05-20")),
		item(string("gemini-2.0-flash")),
		item(string("gemini-2.0-flash-lite")),
	}
    ```
## License
[The Unlicense](./LICENSE)

---
Powered with [Bubbletea](https://github.com/charmbracelet/bubbletea) • [Glamour](https://github.com/charmbracelet/glamour)
