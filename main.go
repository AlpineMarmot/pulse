package main

import (
	"flag"
	"github.com/AlpineMarmot/pulse/database"
	"github.com/AlpineMarmot/pulse/logger"
	"github.com/AlpineMarmot/pulse/middleware"
	"github.com/AlpineMarmot/pulse/util"
)

var db database.MongoDb
var currentSessionId interface{}

func main() {

	configFile := *flag.String("c", "", "Configuration file")
	noLog := flag.Bool("no-logging", false, "Turn off file logging")
	flag.Parse()
	url := flag.Arg(0)

	if false == *noLog {
		logger.New("pulse.log")
	} else {
		logger.New(nil)
	}

	pulse := NewPulse()
	pulse.SetEntryPoint(url)
	pulse.LoadConfigFile(configFile)

	// connect to database
	logger.Println("Connecting to mongodb ... ")
	db = database.NewMongoDb(pulse.config.Pulse.Mongo.Uri, pulse.config.Pulse.Mongo.Database)
	err := db.Connect()
	util.FatalError("Failed to connect to mongodb "+pulse.config.Pulse.Mongo.Uri, err)

	// create a session id
	currentSessionId = pulse.CreateSessionID(db)
	logger.Println("Session:", currentSessionId)

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

	// defer colly internal storage and pulse file logger if applicable
	defer pulse.CloseStorage()
	defer logger.CloseLogFile()

	// crawl!
	pulse.Start()
}
