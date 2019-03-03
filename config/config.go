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
			DomainRegexp string        `yaml:"domain-regexp"`
			DomainGlob   string        `yaml:"domain-glob"`
			Delay        time.Duration `yaml:"delay"`
			RandomDelay  time.Duration `yaml:"random-delay"`
			Parallelism  int           `yaml:"parallelism"`
		}
		MaxBodySize         int           `yaml:"max-body-size"`
		MaxDepth            int           `yaml:"max-depth"`
		MaxUrlVisited       int64         `yaml:"max-url-visited"`
		SleepBetweenRequest time.Duration `yaml:"sleep-between-request"`
		Storage             struct {
			ClearOnStart bool `yaml:"clear-on-start"`
			Redis        struct {
				Address  string
				Password string
				Db       int
				Prefix   string
			}
			MongoDb struct {
				Database string
				Uri      string
			}
		}
		Proxy      []string `yaml:"proxy"`
		UserAgents []string `yaml:"user-agents"`
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
