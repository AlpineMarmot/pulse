package middleware

import (
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"pulse/database"
	"time"
)

func StoreResponse(db database.MongoDb, parentId interface{}) func(r *colly.Response) {
	return func(r *colly.Response) {
		coll := db.Collection("responses")
		_, _ = coll.InsertOne(db.GetQueryContext(), bson.M{
			"sessionId":  parentId,
			"dt_created": time.Now(),
			"url":        r.Ctx.Get("url"),
			"headers":    r.Headers,
			"body":       r.Body,
		})
	}
}
