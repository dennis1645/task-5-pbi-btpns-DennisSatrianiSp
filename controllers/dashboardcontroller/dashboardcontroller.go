package dashboardcontroller

import (
	"net/http"
	"strconv"

	"github.com/dennis1645/go-api/helper"
	"github.com/dennis1645/go-api/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	// Ambil data user berdasarkan ID
	var user models.Users
	if err := models.DB.Where("id = ?", id).Select("id, username, email").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"pesan": "User tidak ditemukan!"}
			helper.ResponseJSON(w, http.StatusNotFound, response)
			return
		} else {
			response := map[string]string{"pesan": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Ambil data user berdasarkan ID
	var photos models.Photo
	if err := models.DB.Where("user_id = ?", id).Select("id, title, photo_url,caption,user_id").First(&photos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"pesan": "Photo tidak ditemukan!"}
			helper.ResponseJSON(w, http.StatusNotFound, response)
			return
		} else {
			response := map[string]string{"pesan": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Jika data ditemukan, kirimkannya sebagai respons JSON
	data := map[string]interface{}{
		"id":          user.ID,
		"nama":        user.Username,
		"email":       user.Email,
		"photo title": photos.Title,
		"photo_url":   photos.PhotoURL,
	}

	helper.ResponseJSON(w, http.StatusOK, data)
}
