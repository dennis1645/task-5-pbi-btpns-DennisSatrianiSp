package authcontroller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dennis1645/go-api/config"
	"github.com/dennis1645/go-api/helper"
	"github.com/dennis1645/go-api/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Mengambil inputan json
	var userInput models.Users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"pesan": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// ambil data user berdasarkan username
	var user models.Users
	if err := models.DB.Where("username = ?", userInput.Username).Select("id,username,password").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := map[string]string{"pesan": "Username dan password salah!"}
			helper.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		} else {
			response := map[string]string{"pesan": err.Error()}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		response := map[string]string{"pesan": "Username dan password salah!"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Proses pembuatan jwt
	expTime := time.Now().Add(time.Minute * 10)
	claims := &config.JWTClaim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Issuer:    "go-jwt-mux",
		},
	}

	// Mendeklarasikan algoritma yang digunakan untuk penandatanganan token
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Tandatangani token
	token, err := tokenAlgo.SignedString([]byte(config.JWT_KEY))
	if err != nil {
		response := map[string]string{"pesan": err.Error()}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}
	// Set token dalam cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
	})

	response := map[string]string{"pesan": "Login Berhasil!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Mengambil inputan json
	var userInput models.Users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"pesan": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// hash password menggunakan bcrypt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		response := map[string]string{"pesan": "Register gagal (Gagal mengenkripsi kata sandi)!"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	userInput.Password = string(hashPassword)

	// Insert data ke database
	if err := models.DB.Create(&userInput).Error; err != nil {
		response := map[string]string{"pesan": "Register gagal (Gagal menyimpan data)!"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}
	response := map[string]string{"pesan": "Register berhasil!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Hapus token dalam cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"pesan": "Logout Berhasil!"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
