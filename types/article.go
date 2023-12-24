package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Slug        string             `bson:"slug"`
	Description string             `bson:"description"`
	Content     string             `bson:"content"`
	Author      primitive.ObjectID `bson:"author"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type SortDirection int

type ArticleField string

type ArticleFilter struct {
	Author       string
	Title        string
	Description  string
	Content      string
	CreatedAtMin time.Time
	CreatedAtMax time.Time
	UpdatedAtMax time.Time
	UpdatedAtMin time.Time
}

type ArticleOptions struct {
	Limit  int
	Skip   int
	SortBy ArticleField
	SortDirection
}
