package cache

import (
	"log"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
)

var articles map[string]string = make(map[string]string)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	root, err := os.Getwd()
	check(err)
	dir := path.Join(root, "/static/articles")

	files, err := os.ReadDir(dir)
	check(err)

	for _, file := range files {
		filepath := path.Join(dir, file.Name())
		b, err := os.ReadFile(filepath)
		check(err)

		content := markdown.ToHTML(b, nil, nil)
		articles[file.Name()] = string(content)
	}
}

func GetArticle(slug string) string {
	return articles[slug]
}
