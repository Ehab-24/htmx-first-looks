package scripts

import (
	"log"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CacheArticles() map[string]string {

	root, err := os.Getwd()
	check(err)
	dir := path.Join(root, "/static/articles")

	files, err := os.ReadDir(dir)
	check(err)

	articles := make(map[string]string)
	for _, file := range files {
		filepath := path.Join(dir, file.Name())
		b, err := os.ReadFile(filepath)
		check(err)

		content := markdown.ToHTML(b, nil, nil)
		articles[file.Name()] = string(content)
	}

	return articles
}
