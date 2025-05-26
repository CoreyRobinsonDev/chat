package main

import (
	"encoding/json"
	"os"

	u "github.com/coreyrobinsondev/utils"
)

type Config struct {
	Model string `json:"model"`
}

func (self *Config) Create() {
	self.Model = "gemini-2.5-flash-preview-05-20" 
	bytes := u.Unwrap(json.Marshal(self))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"

	err := os.Mkdir(configPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Fatal(err)
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

func (self *Config) SetModel(model string) {
	self.Model = model
	bytes := u.Unwrap(json.Marshal(self))
	homeDir := u.Unwrap(os.UserHomeDir())
	configPath := homeDir + "/.config/search"
	configFile := u.Unwrap(os.Create(configPath + "/search.conf"))
	defer configFile.Close()
	u.Unwrap(configFile.Write(bytes))
}
