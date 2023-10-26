package photoscontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dennis1645/go-api/helper"
	"github.com/dennis1645/go-api/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Create(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"error": "ID tidak valid"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Mengambil inputan json
	var userInput models.Photo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"pesan": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	userInput.UserID = uint(id)
	// Insert data ke database
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"pesan": "Gagal menyimpan data"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response := map[string]string{"pesan": "Data berhasil disimpan!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	// ambil id user sama id foto
	vars := mux.Vars(r)
	idParam := vars["id"]
	userParam := vars["userId"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return
	}

	userId, err := strconv.Atoi(userParam)
	if err != nil {
		return
	}

	// Ambil data foto berdasarkan ID
	var photos models.Photo
	if err := models.DB.Where("id = ?", id).Select("id, title, photo_url, caption, user_id").First(&photos).Error; err != nil {
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

	// Cek apakah user_id di database cocok dengan userId dari parameter
	if photos.UserID == uint(userId) {
		// Mengambil inputan json
		var userInput models.Photo
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userInput); err != nil {
			response := map[string]string{"pesan": err.Error()}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
		defer r.Body.Close()

		if models.DB.Model(&userInput).Where("id = ?", id).Updates(&userInput).RowsAffected == 0 {
			response := map[string]string{"pesan": "Gagal mengupdate data"}
			helper.ResponseJSON(w, http.StatusForbidden, response)
			return
		}
		response := map[string]string{"pesan": "Data berhasil di update!"}
		helper.ResponseJSON(w, http.StatusOK, response)
	} else {
		response := map[string]string{"pesan": "Anda tidak memiliki izin untuk mengupdate foto ini"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
	}
}
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	userParam := vars["userId"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := map[string]string{"error": "ID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	userId, err := strconv.Atoi(userParam)
	if err != nil {
		response := map[string]string{"error": "UserID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	var photos models.Photo

	if err := models.DB.Where("id = ?", id).Select("id, title, photo_url, caption, user_id").First(&photos).Error; err != nil {
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

	// Cek apakah user_id di database cocok dengan userId dari parameter
	if photos.UserID == uint(userId) {
		if models.DB.Delete(&photos, id).RowsAffected == 0 {
			response := map[string]string{"pesan": "Gagal menghapus data"}
			helper.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"pesan": "Data berhasil dihapus!"}
		helper.ResponseJSON(w, http.StatusOK, response)
	} else {
		response := map[string]string{"pesan": "Anda tidak memiliki izin untuk menghapus foto ini"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	// Ambil data foto berdasarkan ID
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
		"id":          photos.ID,
		"user id":     photos.UserID,
		"photo title": photos.Title,
		"photo url":   photos.PhotoURL,
	}

	helper.ResponseJSON(w, http.StatusOK, data)
}
