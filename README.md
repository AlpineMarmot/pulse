# Pulse

Pulse is a crawler build on top of [gocolly/colly](https://github.com/gocolly/colly)

Features:
 - Expose all golly/colly options to a yml configuration
 - Create rule(s) that export crawling data to MongoDB 
    
### Installation

Go modules must be enabled

 > $ go build
 
### Usage

 > $ pulse -c conf.yml https://www.example.com
 
 ### Configuration example
 
see [default.yml](https://github.com/AlpineMarmot/pulse/blob/master/default.yml)

### Grab HTML data
This rule below will add to mongodb collection "images" the value of `src` attribute for all tag `img`. The `context-attr` is also added as images metadata.

```yaml
collection: "images"
tag: "img"
attr: "src"
context-attr: "alt"
```

You can also grab html attributes with a `selector` instead of `tag`. 

```yaml
collection: "images-test"
selector: "img[data-src]"
attr: "data-src"
context-attr: "alt"
```

More infos about selector here: [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)