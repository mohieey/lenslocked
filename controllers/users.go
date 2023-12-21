package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mohieey/lenslocked/appctx"
	"github.com/mohieey/lenslocked/services"
)

type Users struct {
	Templates struct {
		SignUp Template
	}
	UserService          *services.UserService
	SessionService       *services.SessionService
	PasswordResetService *services.PasswordResetService
	EmailService         *services.EmailService
}

func (u *Users) SignUp(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) SignIn(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := appctx.User(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *Users) SignOut(w http.ResponseWriter, r *http.Request) {
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

func (u *Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	pwReset, err := u.PasswordResetService.Create(email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": []string{pwReset.Token},
	}

	resetUrl := fmt.Sprintf("www.localhost.com:3000/reset_password?%v", vals.Encode())

	err = u.EmailService.ForgotPassword(email, resetUrl)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("check your email"))
}

func (u *Users) ResetPassord(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	// password := r.FormValue("password")
	// confirmPassword := r.FormValue("confirmPassword")

	user, err := u.PasswordResetService.Consume(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
