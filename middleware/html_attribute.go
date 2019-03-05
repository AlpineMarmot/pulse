package middleware

import (
	"github.com/AlpineMarmot/pulse/database"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type HtmlAttributeDefinition struct {
	Collection  string `yaml:"collection"`
	Selector    string `yaml:"selector"`
	Tag         string `yaml:"tag"`
	Attr        string `yaml:"attr"`
	ContextAttr string `yaml:"context-attr"`
}

func GetHtmlAttributeString(htmlAttribute HtmlAttributeDefinition) string {
	if len(htmlAttribute.Selector) > 0 {
		return htmlAttribute.Selector
	}
	return htmlAttribute.Tag + "[" + htmlAttribute.Attr + "]"
}

func HtmlAttribute(db database.MongoDb, htmlAttribute HtmlAttributeDefinition) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		coll := db.Collection(htmlAttribute.Collection)
		value := e.Attr(htmlAttribute.Attr)
		context := ""
		if len(htmlAttribute.ContextAttr) > 0 {
			context = e.Attr(htmlAttribute.ContextAttr)
		}
		_, _ = coll.InsertOne(db.GetQueryContext(), bson.M{
			"dt_created": time.Now(),
			"selector":   GetHtmlAttributeString(htmlAttribute),
			"context":    context,
			"value":      value,
		})
	}
}
