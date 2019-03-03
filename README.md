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
 
see [default.yml](https://github.com/AlpineMarmot/pulse/blob/master/default.yml)