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
