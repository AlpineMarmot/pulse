package config

import (
	"github.com/AlpineMarmot/pulse/middleware"
	"github.com/AlpineMarmot/pulse/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

// pulse configuration structure
type Config struct {
	Pulse   configPulse   `yaml:"pulse"`
	Crawler configCrawler `yaml:"crawler"`
}

type configPulse struct {
	Mongo       configMongoDb                        `yaml:"mongo"`
	Middlewares []string                             `yaml:"middlewares"`
	Html        []middleware.HtmlAttributeDefinition `yaml:"html"`
}

type configCrawler struct {
	AllowedDomains      []string             `yaml:"allowed-domains"`
	AllowURLRevisit     bool                 `yaml:"allow-url-revisit"`
	Async               bool                 `yaml:"async"`
	CheckHead           bool                 `yaml:"check-head"`
	DetectCharset       bool                 `yaml:"detect-charset"`
	DisallowedDomains   []string             `yaml:"disallowed-domains"`
	IgnoreRobotsTxt     bool                 `yaml:"ignore-robots-txt"`
	Limit               configCrawlerLimit   `yaml:"limit"`
	MaxBodySize         int                  `yaml:"max-body-size"`
	MaxDepth            int                  `yaml:"max-depth"`
	MaxUrlVisited       int64                `yaml:"max-url-visited"`
	RandomUserAgents    bool                 `yaml:"random-user-agents"`
	SleepBetweenRequest time.Duration        `yaml:"sleep-between-request"`
	Storage             configCrawlerStorage `yaml:"storage"`
	Proxy               []string             `yaml:"proxy"`
	UserAgents          []string             `yaml:"user-agents"`
}

type configCrawlerLimit struct {
	DomainRegexp string        `yaml:"domain-regexp"`
	DomainGlob   string        `yaml:"domain-glob"`
	Delay        time.Duration `yaml:"delay"`
	RandomDelay  time.Duration `yaml:"random-delay"`
	Parallelism  int           `yaml:"parallelism"`
}

type configCrawlerStorage struct {
	ClearOnStart bool          `yaml:"clear-on-start"`
	Redis        configRedis   `yaml:"redis"`
	Mongo        configMongoDb `yaml:"mongo"`
}

type configMongoDb struct {
	Database string `yaml:"database"`
	Uri      string `yaml:"uri"`
}

type configRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

// Create a Config from a yaml string
func NewConfigFromString(data string) Config {
	newConf := Config{}
	err := yaml.Unmarshal([]byte(data), &newConf)
	//util.CheckError(err, "Loading config from a string")
	util.FatalError("There is a problem with your configuration", err)
	return newConf
}

// Create a Config from yaml file
func NewConfigFromFile(configFilePath string) Config {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		//util.CheckError(err, "Checking existence of a config file")
		util.FatalError("Cannot find the configuration file", err)
	}

	configFileContent, err := ioutil.ReadFile(configFilePath)
	//util.CheckError(err, "Loading config file content")
	util.FatalError("A problem occurred while loading the configuration file", err)

	return NewConfigFromString(string(configFileContent))
}
