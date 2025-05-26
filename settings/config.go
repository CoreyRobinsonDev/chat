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
	GeminiApiKey string `json:"geminiApiKey"`
	GeminiChatHistory []*genai.Content `json:"geminiChatHistory"`
}

func (self *Config) Create() {
	u.Expect(godotenv.Load())
	self.Model = "gemini-2.5-flash-preview-05-20" 
	self.GeminiApiKey = os.Getenv("GEMINI_API_KEY")

	bytes := u.Unwrap(json.Marshal(self))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"

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
	configPath := homeDir + "/.config/search"

	_, err := os.Stat(configPath + "/settings.json")
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func (self *Config) Init() {
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"

	if u.Unwrap(self.IsExist()) {
		content := u.Unwrap(os.ReadFile(configPath + "/settings.json"))
		u.Expect(json.Unmarshal(content, &self))
	} else {
		self.Create()
	}
}


func (self Config) Write() {
	bytes := u.Unwrap(json.Marshal(self))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"
	configFile := u.Unwrap(os.Create(configPath + "/settings.json"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}

