package templates

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

const InternalPublic = template.HTML("internal error<br>try again later")

// // // //
const relative = "../../web/templates"

// // // //
var (
	Message = load("/msg.html")
	Main    = load("/index.html")
)

// // // //

type Template struct {
	tmpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data string) error {
	return t.tmpl.Execute(w, template.HTML(data))
}

func (t Template) Internal(w http.ResponseWriter) error {
	return t.tmpl.Execute(w, InternalPublic)
}

func load(files ...string) *Template {
	for i := range files {
		files[i] = relative + files[i]
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("template: parse files '%s': %s", strings.Join(files, ","), err)
	}

	return &Template{tmpl: t}
}
