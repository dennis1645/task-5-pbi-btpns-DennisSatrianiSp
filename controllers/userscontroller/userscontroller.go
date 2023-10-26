package userscontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dennis1645/go-api/helper"
	"github.com/dennis1645/go-api/models"
	"github.com/gorilla/mux"
)

func Update(w http.ResponseWriter, r *http.Request) {
	// ambil id user sama id foto
	vars := mux.Vars(r)
	userParam := vars["userId"]
	userId, err := strconv.Atoi(userParam)
	if err != nil {
		return
	}
	// Mengambil inputan json
	var userInput models.Users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"pesan": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	if models.DB.Model(&userInput).Where("id = ?", userId).Updates(&userInput).RowsAffected == 0 {
		response := map[string]string{"pesan": "Gagal mengupdate data"}
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}
	response := map[string]string{"pesan": "Data User berhasil di update!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userParam := vars["userId"]
	userId, err := strconv.Atoi(userParam)
	if err != nil {
		response := map[string]string{"error": "UserID tidak valid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	var users models.Users
	if models.DB.Delete(&users, userId).RowsAffected == 0 {
		response := map[string]string{"pesan": "Gagal menghapus data"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	// Hapus token dalam cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"pesan": "Hapus Akun berhasil :: Logout Berhasil!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
