package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
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
			Parallelism string `yaml:"parallelism"`
		}
		MaxUrlVisited       int64         `yaml:"max-url-visited"`
		SleepBetweenRequest time.Duration `yaml:"sleep-between-request"`
		UserAgents          []string      `yaml:"user-agents"`
	}
}

// Create a Config from a yaml string
func NewConfigFromString(data string) Config {
	newConf := Config{}

	err := yaml.Unmarshal([]byte(data), &newConf)
	checkError(err, "Loading config from a string")
	return newConf
}

// Create a Config from yaml file
func NewConfigFromFile(configFilePath string) Config {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		checkError(err, "Checking existence of a config file")
	}

	configFileContent, err := ioutil.ReadFile(configFilePath)
	checkError(err, "Loading config file content")

	return NewConfigFromString(string(configFileContent))
}

// check config file existence
func ConfigFileExists(configFilePath string) bool {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Apply a Config to Colly
func ApplyConfigToColly(config Config, c *colly.Collector) {
	c.AllowedDomains = config.Crawler.AllowedDomains
	c.AllowURLRevisit = config.Crawler.AllowURLRevisit
	c.DetectCharset = config.Crawler.DetectCharset
	c.DisallowedDomains = config.Crawler.DisallowedDomains
	c.CheckHead = config.Crawler.CheckHead
	c.Async = config.Crawler.Async
	c.IgnoreRobotsTxt = config.Crawler.IgnoreRobotsTxt

	realLimitParallelism := 1
	limitParallelism := string(config.Crawler.Limit.Parallelism)
	if limitParallelism == "auto" {
		realLimitParallelism = runtime.NumCPU()
	} else {
		x, _ := strconv.Atoi(string(config.Crawler.Limit.Parallelism))
		realLimitParallelism = x
	}

	fmt.Println(realLimitParallelism)
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: realLimitParallelism,
		//RandomDelay: 15 * time.Second,
	})
	checkError(err, "Applying Config to a colly collector")
}
