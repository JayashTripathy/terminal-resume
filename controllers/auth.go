package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"terminal-resume.jayash.space/models"
	"terminal-resume.jayash.space/utils"
)

var jwtKey = []byte("my_secret")

var Users = map[string]string{}


func Signup(w http.ResponseWriter, r *http.Request) error {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}

	  // Manually validate the fields
	  if user.Name == "" || len(user.Name) < 3 || len(user.Name) > 32 {
        return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Name must be between 3 and 32 characters"})
    }

    if user.Email == "" || !utils.IsValidEmail(user.Email) {
        return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Invalid email address"})
    }

    if user.Password == "" || len(user.Password) < 8 {
        return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Password must be at least 8 characters long"})
    }


	var existingUser  models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "User already exists"})
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, utils.ApiError{Error: errHash.Error()})
	}

	models.DB.Create(&user)
	log.Printf("User created: %v", user)
	return utils.WriteJSON(w, http.StatusOK, utils.ApiResponse{Message: "User created"})
	
}

func Login(w http.ResponseWriter, r *http.Request) error {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "User not found"})
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Invalid password"})
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString , err := token.SignedString(jwtKey)

	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, utils.ApiError{Error: err.Error()})
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
        Value:    tokenString,
        Expires:  expirationTime,
        HttpOnly: true,
	})

	return utils.WriteJSON(w, http.StatusOK, utils.ApiResponse{Message: "Logged in"})
}


