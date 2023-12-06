package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		SignUp Template
	}
}

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) {
	//******************************** STUDY NOTES ****************************************************************
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// fmt.Fprint(w, r.PostForm)
	// fmt.Fprint(w, r.PostForm.Get("email"))
	//******************************** STUDY NOTES ****************************************************************

	fmt.Fprintln(w, r.FormValue("email"))
	fmt.Fprintln(w, r.FormValue("password"))
}
