package middleware

import (
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"pulse/database"
	"time"
)

func StoreRequest(db database.MongoDb, parentId interface{}) func(r *colly.Request) {
	return func(r *colly.Request) {
		coll := db.Collection("request")
		_, _ = coll.InsertOne(db.GetQueryContext(), bson.M{
			"sessionId":  parentId,
			"dt_created": time.Now(),
			"method":     r.Method,
			"url":        r.Ctx.Get("url"),
			"request":    r.Headers,
		})

	}
}
