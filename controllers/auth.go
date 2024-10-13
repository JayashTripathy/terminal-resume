package controllers

import (
	"encoding/json"
	"net/http"

	"terminal-resume.jayash.space/models"
	"terminal-resume.jayash.space/utils"
)

// var jwtKey = []byte("my_secret")

var Users = map[string]string{}

func Login(w http.ResponseWriter, r *http.Request) error {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if(err != nil) {
		utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	// models.DB.Where()
	
	return nil
}