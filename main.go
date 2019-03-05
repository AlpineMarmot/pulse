package main

import (
	"flag"
	"fmt"
	"pulse/database"
	"pulse/middleware"
	"pulse/util"
)

var db database.MongoDb
var currentSessionId interface{}

func main() {
	configFile := flag.String("c", "", "Configuration file")
	flag.Parse()
	url := flag.Arg(0)

	pulse := NewPulse()
	pulse.SetEntryPoint(url)
	pulse.LoadConfigFile(*configFile)

	// connect to database
	db = database.NewMongoDb(pulse.config.Pulse.Mongo.Uri, pulse.config.Pulse.Mongo.Database)
	err := db.Connect()
	util.CheckError(err, "Connecting to mongo database")

	// create a session id
	currentSessionId = pulse.CreateSessionID(db)
	fmt.Println(currentSessionId)

	// register middlewares
	pulse.OnRequest(middleware.StoreRequest(db, currentSessionId))
	pulse.OnResponse(middleware.StoreResponse(db, currentSessionId))

	for _, html := range pulse.config.Pulse.Html {
		htmlAttr := middleware.HtmlAttributeDefinition{
			Collection:  html.Collection,
			Selector:    html.Selector,
			Tag:         html.Tag,
			Attr:        html.Attr,
			ContextAttr: html.ContextAttr,
		}
		pulse.OnHTML(middleware.GetHtmlAttributeString(htmlAttr), middleware.HtmlAttribute(db, htmlAttr))
	}

	// defer colly internal storage if applicable
	defer pulse.CloseStorage()

	// crawl!
	pulse.Start()
}
