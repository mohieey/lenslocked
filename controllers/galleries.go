package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mohieey/lenslocked/appctx"
	"github.com/mohieey/lenslocked/models"
	"github.com/mohieey/lenslocked/services"
)

type Galleries struct {
	GalleryService *services.GalleryService
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	user := appctx.User(r.Context())
	title := r.FormValue("title")

	gallery, err := g.GalleryService.Create(title, int(user.ID))
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gallery)
}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return
	}

	title := r.FormValue("title")

	if !isGalleryOwner(w, r, gallery) {
		return
	}

	gallery.Title = title

	err = g.GalleryService.Update(gallery)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gallery)
}

func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := appctx.User(r.Context())

	galleries, err := g.GalleryService.GetByUserId(int(user.ID))
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(galleries)
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gallery)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return
	}

	if !isGalleryOwner(w, r, gallery) {
		return
	}

	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("gallery deleted successfully"))
}

func (g *Galleries) getGalleryById(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.GalleryService.GetById(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return nil, err
	}

	return gallery, nil
}

func isGalleryOwner(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) bool {
	user := appctx.User(r.Context())

	if gallery.UserId != int(user.ID) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	return true
}
