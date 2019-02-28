# Pulse

Pulse is a crawler build on top of [gocolly/colly](https://github.com/gocolly/colly)

Features:
 - Yml configuration
 - MongoDB support 
    
### Installation

Go modules must be enabled

 > $ go build
 
### Usage

 > $ pulse -c conf.yml https://www.example.com
 
 ### Configuration example

```yml
mongo:
  address: "mongodb://localhost:27017"
  database: "scrapping"
crawler:
  allow-url-revisit: false
  #allowed-domains:
  # - www.google.com
  async: true
  detect-charset: false
  #disallowed-domains:
  # - www.google.com
  ignore-robots-txt: false
  limit:
    parallelism: "16"
  max-url-visited: 1000
  random-users-agents: false
  sleep-between-request: 0
  user-agents:
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.0129.115 Safari/537.36"
```

