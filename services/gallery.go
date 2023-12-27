package services

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mohieey/lenslocked/models"
)

const DefaultImagesDir = "images"

type GalleryService struct {
	DB                  *sql.DB
	ImagesDir           string
	SupportedExtensions []string
}

func (gs *GalleryService) Create(title string, userID int) (*models.Gallery, error) {
	gallery := models.Gallery{
		Title:  title,
		UserId: userID,
	}

	row := gs.DB.QueryRow(
		`
		INSERT INTO galleries(title, user_id)
		VALUES($1, $2) RETURNING id;
		`,
		gallery.Title, gallery.UserId,
	)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating gallery: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) GetById(id int) (*models.Gallery, error) {
	gallery := models.Gallery{
		ID: id,
	}

	row := gs.DB.QueryRow(
		`
		SELECT title, user_id 
		FROM galleries
		WHERE id = $1;
		`,
		gallery.ID,
	)
	err := row.Scan(&gallery.Title, &gallery.UserId)
	if err != nil {
		return nil, fmt.Errorf("error getting gallery: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) GetByUserId(userID int) ([]models.Gallery, error) {

	rows, err := gs.DB.Query(
		`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;
		`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting galleries by user id: %w", err)
	}

	var galleries []models.Gallery
	for rows.Next() {
		gallery := models.Gallery{
			UserId: userID,
		}

		err = rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("error getting galleries by user id: %w", err)
		}

		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error getting galleries by user id: %w", err)
	}

	return galleries, nil
}

func (gs *GalleryService) Update(gallery *models.Gallery) error {
	_, err := gs.DB.Exec(
		`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;
		`,
		gallery.ID, gallery.Title,
	)
	if err != nil {
		return fmt.Errorf("error updating gallery: %w", err)
	}

	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(
		`
		DELETE FROM galleries
		WHERE id = $1;
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting gallery: %w", err)
	}

	return nil
}

func (gs *GalleryService) galleryDir(id int) string {
	if gs.ImagesDir == "" {
		gs.ImagesDir = DefaultImagesDir
	}
	return filepath.Join(gs.ImagesDir, fmt.Sprintf("gallery-%d", id))
}

func (gs *GalleryService) hasExtension(imagePath string, extensions []string) bool {
	imagePath = strings.ToLower(imagePath)
	for _, extension := range extensions {
		extension = strings.ToLower(extension)
		if filepath.Ext(imagePath) == extension {
			return true
		}
	}

	return false
}

func (gs *GalleryService) Images(galleryID int) ([]models.Image, error) {
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("error getting images: %w", err)
	}

	images := []models.Image{}

	for _, file := range allFiles {
		if gs.hasExtension(file, gs.SupportedExtensions) {
			images = append(images, models.Image{
				GalleryID: galleryID,
				Path:      file,
				FileName:  filepath.Base(file),
			})
		}
	}

	return images, nil
}
