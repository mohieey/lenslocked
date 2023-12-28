package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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

	gallery.Images, err = g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("gallery deleted successfully"))
}

func (g *Galleries) ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, image.Path)
}

func (g *Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	gallery, err := g.getGalleryById(w, r)
	if err != nil {
		return
	}

	if !isGalleryOwner(w, r, gallery) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("image deleted successfully"))
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
