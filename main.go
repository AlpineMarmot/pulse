package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/mongodb/mongo-go-driver/bson"
	"time"
)

func main() {
	configFile := flag.String("c", DefaultConfigFile, "Configuration file")
	flag.Parse()
	url := flag.Arg(0)

	pulse := NewPulse()
	pulse.SetEntryPoint(url)
	pulse.LoadConfigFile(*configFile)

	extensions.RandomUserAgent(pulse.colly)

	db := NewMongoDb(pulse.config.Mongo.Address, pulse.config.Mongo.Database)
	err := db.connect()
	checkError(err, "Connecting to mongo database")

	responseCollection := db.Collection("sessions")
	res, _ := responseCollection.InsertOne(db.GetQueryContext(), bson.M{
		"dt_created": time.Now(),
	})
	sessionId := res.InsertedID
	fmt.Println(sessionId)

	pulse.OnHTML("img[data-src]", func(e *colly.HTMLElement) {
		src := e.Attr("data-src")
		fmt.Println(src)
	})

	pulse.Start()
}
