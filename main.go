package main

import (
	"flag"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"pulse/database"
	"pulse/middleware"
	"pulse/util"
	"time"
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

	db = database.NewMongoDb(pulse.config.Mongo.Address, pulse.config.Mongo.Database)
	err := db.Connect()
	util.CheckError(err, "Connecting to mongo database")

	currentSessionId = CreateSessionId()
	fmt.Println(currentSessionId)

	pulse.OnRequest(middleware.StoreRequest(db, currentSessionId))
	pulse.OnResponse(middleware.StoreResponse(db, currentSessionId))
	pulse.OnHTML(middleware.GrabImageUrlSelector(), middleware.GrabImageUrl(db, currentSessionId))

	defer pulse.CloseStorage()
	pulse.Start()
}

func CreateSessionId() interface{} {
	coll := db.Collection("sessions")
	res, _ := coll.InsertOne(db.GetQueryContext(), bson.M{
		"dt_created": time.Now(),
	})
	return res.InsertedID
}
