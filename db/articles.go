package db

import (
	"context"
	"log"

	"github.com/Ehab-24/test/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ASC  types.SortDirection = 1
	DESC types.SortDirection = -1
)

const (
	CreatedAt types.ArticleField = "created_at"
	UpdatedAt types.ArticleField = "updated_at"
)

func GetArticles(filter *types.ArticleFilter, opts *types.ArticleOptions) ([]types.Article, error) {
	_filter := bson.M{}
	if filter.Author != "" {
		_filter["author"] = filter.Author
	}
	if filter.Title != "" {
		_filter["title"] = bson.M{"$regex": filter.Title}
	}
	if filter.Description != "" {
		_filter["description"] = bson.M{"$regex": filter.Title}
	}
	// if filter.Content != "" {
	// 	_filter["content"] = bson.M{"$regex": filter.Content}
	// }
	if !filter.CreatedAtMin.IsZero() {
		_filter["created_at"] = bson.M{"$gte": filter.CreatedAtMin}
	}
	if !filter.CreatedAtMax.IsZero() {
		_filter["created_at"] = bson.M{"$lte": filter.CreatedAtMax}
	}
	if !filter.UpdatedAtMin.IsZero() {
		_filter["updated_at"] = bson.M{"$gte": filter.UpdatedAtMin}
	}
	if !filter.UpdatedAtMax.IsZero() {
		_filter["updated_at"] = bson.M{"$lte": filter.UpdatedAtMax}
	}

	log.Println(opts.Limit, opts.Skip)

	_options := options.Find()
	if opts.Limit != 0 {
		_options.SetLimit(int64(opts.Limit))
	}
	if opts.Skip != 0 {
		_options.SetSkip(int64(opts.Skip))
	}
	if opts.SortBy != "" {
		_options.SetSort(bson.M{string(opts.SortBy): opts.SortDirection})
	}

	log.Println(_options.Limit, _options.Skip)

	_articles, err := Db.Collection("posts").Find(context.Background(), _filter)
	if err != nil {
		return nil, err
	}
	var articles []types.Article
	if err := _articles.All(context.Background(), &articles); err != nil {
		return nil, err
	}

	return articles, nil
}
