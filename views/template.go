package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func Parse(filePath string) (*Template, error) {
	tpl, err := template.ParseFiles(filepath.Join("templates/", filePath+".gohtml"))
	if err != nil {
		return &Template{}, fmt.Errorf("Error parsing %w", err)
	}

	return &Template{htmlTemplate: tpl}, nil
}

type Template struct {
	htmlTemplate *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTemplate.Execute(w, data)
	if err != nil {
		log.Printf("Error excuting %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
