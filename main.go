package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/extensions"
	"github.com/mongodb/mongo-go-driver/bson"
	"pulse/database"
	"pulse/util"
	"time"
)

func main() {
	configFile := flag.String("c", "", "Configuration file")
	flag.Parse()
	url := flag.Arg(0)

	pulse := NewPulse()
	pulse.SetEntryPoint(url)
	pulse.LoadConfigFile(*configFile)

	extensions.RandomUserAgent(pulse.colly)

	db := database.NewMongoDb(pulse.config.Mongo.Address, pulse.config.Mongo.Database)
	err := db.Connect()
	util.CheckError(err, "Connecting to mongo database")

	responseCollection := db.Collection("sessions")
	res, _ := responseCollection.InsertOne(db.GetQueryContext(), bson.M{
		"dt_created": time.Now(),
	})
	sessionId := res.InsertedID
	fmt.Println(sessionId)

	//countRequest := 0
	//pulse.OnRequest(func(request *colly.Request) {
	//	countRequest++
	//	fmt.Println("Request #", countRequest)
	//})

	//pulse.OnHTML("img[data-src]", func(e *colly.HTMLElement) {
	//	src := e.Attr("data-src")
	//	fmt.Println(src)
	//})

	defer pulse.CloseStorage()
	pulse.Start()
}
