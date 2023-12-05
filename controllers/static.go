package controllers

import (
	"net/http"

	"github.com/mohieey/lenslocked/views"
)

func StaticHandler(tmpl *views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func FAQ(tmpl *views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Is it free?",
			Answer:   "Yes",
		},
		{
			Question: "How may images can I upload?",
			Answer:   "100 images per month",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, questions)
	}
}
