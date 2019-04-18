package HTMLHandlers

import (
	"html/template"
	"net/http"
)

type Templates map[string]*template.Template

func (t Templates) RenderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, ok := t[name]
	if !ok {
		http.Error(w, "can't find template", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
