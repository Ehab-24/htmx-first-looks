package main

import (
	"html/template"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type HTMLScript string

const (
	TailwindCSS HTMLScript = "tailwindcss"
	HTMX        HTMLScript = "htmx"
	SimpleMDE   HTMLScript = "simplemde"
)

var HTMLScripts = map[HTMLScript]template.HTML{
	TailwindCSS: `<script src="https://cdn.tailwindcss.com"></script>`,
	HTMX: `
 <script src="https://unpkg.com/htmx.org@1.9.6"
   integrity="sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni"
   crossorigin="anonymous"></script>`,
	SimpleMDE: `<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css" />
	`,
}
