package middleware

import (
	"github.com/AlpineMarmot/pulse/database"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func GrabImageUrlSelector() string {
	return "img[src]"
}

func GrabImageUrl(db database.MongoDb, parentId interface{}) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		coll := db.Collection("images")

		src := e.Attr("src")
		alt := e.Attr("alt")

		if len(e.Attr("data-src")) > 0 {
			src = e.Attr("data-src")
		}

		_, _ = coll.InsertOne(db.GetQueryContext(), bson.M{
			"sessionId":  parentId,
			"dt_created": time.Now(),
			"image":      src,
			"title":      alt,
		})

	}
}
