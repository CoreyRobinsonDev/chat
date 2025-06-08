# chat
Terminal based AI chat client running on Google Gemini models
<br>

![demo](https://vhs.charm.sh/vhs-4nLDenhbKScUPvy58sEDWt.gif)

## Installation
1. Build from source
    ```bash
    git clone https://github.com/CoreyRobinsonDev/chat
    cd chat
    go build
    ```
1. Add Gemini API key. This can be done in 2 ways:
    - Within the project directory create a `.env` file, add `GEMINI_API_KEY=your_api_key`, then run `go build`. This will be added to your settings file at `~/.config/chat/settings.json`. `GEMINI_API_KEY` can also be exported as system environment variable
    - Alternitively, you can update the `geminiApiKey` field in your `~/.config/chat/settings.json` file directly

## Commands
- `config`: Select the AI model to use
    - Other [Gemini models](https://ai.google.dev/gemini-api/docs/models) can be added to `geminiModels` in `~/.config/chat/settings.json`

## License
[The Unlicense](./LICENSE)

---
Powered with [Bubbletea](https://github.com/charmbracelet/bubbletea) • [Glamour](https://github.com/charmbracelet/glamour)
