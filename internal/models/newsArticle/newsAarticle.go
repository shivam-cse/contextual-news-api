package newsArticle

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsArticleDBResponse struct {
	ID              string    `bson:"_id" json:"article_id"`
    Title           string    `bson:"title" json:"title"`
    Description     string    `bson:"description" json:"description"`
    URL             string    `bson:"url" json:"url"`
    PublicationDate interface{} `bson:"publication_date" json:"publication_date"`
    SourceName      string    `bson:"source_name" json:"source_name"`
    RelevanceScore  float64   `bson:"relevance_score" json:"relevance_score"`
    Latitude        float64   `bson:"latitude" json:"latitude"`
    Longitude       float64   `bson:"longitude" json:"longitude"`
    Category        []string  `bson:"category" json:"category"`
	LLMSummary      string    `bson:"llm_summary" json:"llm_summary"`
}

type UserEvent struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ArticleID string             `json:"article_id" bson:"article_id"`
	EventType string             `json:"event_type" bson:"event_type"` // view, click, share
	Latitude  float64            `json:"latitude" bson:"latitude"`
	Longitude float64            `json:"longitude" bson:"longitude"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
}