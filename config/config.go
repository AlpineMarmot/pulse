package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"pulse/util"
	"time"
)

type Config struct {
	Mongo struct {
		Address  string
		Database string
	}
	Crawler struct {
		AllowedDomains    []string `yaml:"allowed-domains"`
		AllowURLRevisit   bool     `yaml:"allow-url-revisit"`
		Async             bool     `yaml:"async"`
		CheckHead         bool     `yaml:"check-head"`
		DetectCharset     bool     `yaml:"detect-charset"`
		DisallowedDomains []string `yaml:"disallowed-domains"`
		IgnoreRobotsTxt   bool     `yaml:"ignore-robots-txt"`
		Limit             struct {
			Parallelism interface{} `yaml:"parallelism"`
		}
		MaxDepth            int           `yaml:"max-depth"`
		MaxUrlVisited       int64         `yaml:"max-url-visited"`
		SleepBetweenRequest time.Duration `yaml:"sleep-between-request"`
		Proxy               []string      `yaml:"proxy"`
		UserAgents          []string      `yaml:"user-agents"`
	}
}

// Create a Config from a yaml string
func NewConfigFromString(data string) Config {
	newConf := Config{}

	err := yaml.Unmarshal([]byte(data), &newConf)
	util.CheckError(err, "Loading config from a string")
	return newConf
}

// Create a Config from yaml file
func NewConfigFromFile(configFilePath string) Config {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		util.CheckError(err, "Checking existence of a config file")
	}

	configFileContent, err := ioutil.ReadFile(configFilePath)
	util.CheckError(err, "Loading config file content")

	return NewConfigFromString(string(configFileContent))
}

// check config file existence
func ConfigFileExists(configFilePath string) bool {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDefaultConfig() string {
	return `
mongo:
  address: "mongodb://localhost:27017"
  database: "scrapping"
crawler:
  allow-url-revisit: false
  #allowed-domains:
  #  - www.google.com
  async: true
  detect-charset: false
  #disallowed-domains:
  #  - www.google.com
  ignore-robots-txt: false
  limit:
    parallelism: "1"
  max-url-visited: 3
  random-users-agents: false
  sleep-between-request: 0
  user-agents:
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.0129.115 Safari/537.36"

`
}
