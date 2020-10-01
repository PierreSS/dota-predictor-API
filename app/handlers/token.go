package handlers

import (
	"dota-predictor/app/config"
	"dota-predictor/app/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// check the validity of the token and decrement a call if bool is true
func isValidToken(w http.ResponseWriter, token string, decrementCall bool) bool {
	var user models.Users

	err := config.DB.Where("access_token = ?", token).Find(&user).Error
	if err == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "Unknown token: " + err.Error()})
		return false
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "There was a problem retrieving token from the database: " + err.Error()})
		return false
	}

	if decrementCall {
		if user.NBCallsLeft == 0 {
			w.WriteHeader(http.StatusLocked)
			json.NewEncoder(w).Encode(models.Response{Code: -1, Message: "No call left available: " + err.Error()})
			return false
		}
		user.NBCallsLeft--
		config.DB.Save(&user)
	}

	log.Println("User " + token + " verified.")
	return true
}
