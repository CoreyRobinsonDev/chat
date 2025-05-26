package settings

import (
	"encoding/json"
	"fmt"
	"os"

	u "github.com/coreyrobinsondev/utils"
	"github.com/joho/godotenv"
)



var ConfigFile Config

type Config struct {
	Model string `json:"model"`
	GeminiApiKey string `json:"geminiApiKey"`
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
	configFile := u.Unwrap(os.Create(configPath + "/search.conf"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}

func (self Config) IsExist() (bool, error) {
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"

	_, err := os.Stat(configPath + "/search.conf")
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func (self *Config) Init() {
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"

	if u.Unwrap(self.IsExist()) {
		content := u.Unwrap(os.ReadFile(configPath + "/search.conf"))
		u.Expect(json.Unmarshal(content, &self))
	} else {
		self.Create()
	}
}


func (self *Config) Set(key string, val any) {
	switch key {
	case "model":
		self.Model = fmt.Sprintf("%v", val)
	case "geminiApiKey":
		self.GeminiApiKey = fmt.Sprintf("%v", val)
	}

	bytes := u.Unwrap(json.Marshal(self))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"
	configFile := u.Unwrap(os.Create(configPath + "/search.conf"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}

