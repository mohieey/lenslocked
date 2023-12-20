package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohieey/lenslocked/appctx"
	"github.com/mohieey/lenslocked/models"
)

type Users struct {
	Templates struct {
		SignUp Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
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

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(int(user.ID))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, newCookie(CookieSession, session.Token))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(int(user.ID))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, newCookie(CookieSession, session.Token))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := appctx.User(r.Context())
	if user == nil {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie(CookieSession)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	err = u.SessionService.Delete(tokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	deleCookie(w, CookieSession)

	w.Write([]byte("sucessfully signed out"))
}
