package settings

import (
	"encoding/json"
	"os"

	u "github.com/coreyrobinsondev/utils"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
)



var ConfigFile Config

type Config struct {
	Model string `json:"model"`
	SystemInstruction string `json:"systemInstruction"`
	GeminiApiKey string `json:"geminiApiKey"`
	GeminiModels []string `json:"geminiModels"`
	GeminiChatHistory []*genai.Content `json:"geminiChatHistory"`
}

func (self *Config) Create() {
	godotenv.Load()
	self.GeminiApiKey = os.Getenv("GEMINI_API_KEY")
	self.GeminiModels = []string{
		"gemini-2.5-flash-preview-05-20",
		"gemini-2.0-flash",
		"gemini-2.0-flash-lite",
	}
	self.Model = self.GeminiModels[0] 
	self.SystemInstruction = "You're a senior software engineer giving short and concise answers. Include code examples"

	bytes := u.Unwrap(json.MarshalIndent(self, "", "\t"))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/chat"

	err := os.Mkdir(configPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		Logger.Fatal(err)
	}
	configFile := u.Unwrap(os.Create(configPath + "/settings.json"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}

func (self Config) IsExist() (bool, error) {
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/chat"

	_, err := os.Stat(configPath + "/settings.json")
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func (self *Config) Init() {
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/chat"

	if u.Unwrap(self.IsExist()) {
		content := u.Unwrap(os.ReadFile(configPath + "/settings.json"))
		u.Expect(json.Unmarshal(content, &self))
	} else {
		self.Create()
	}
}


func (self Config) Write() {
	bytes := u.Unwrap(json.MarshalIndent(self, "", "\t"))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/chat"
	configFile := u.Unwrap(os.Create(configPath + "/settings.json"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}

