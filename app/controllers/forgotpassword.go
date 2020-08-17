package controllers

import (
	"chat-backend/app/lib"
	"chat-backend/app/models"
	"chat-backend/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// ForgotPassword user
func ForgotPassword(response http.ResponseWriter, request *http.Request) {
	user := &models.User{}
	verification := &models.Verification{}

	// Get all request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the request
	err = json.Unmarshal(body, &user)

	// Validate the user input
	err = user.Validate("forgot-password")
	if err != nil {
		lib.Error(response, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Get the user data by E-Mail
	userData, err := user.GetUserByEmail(config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	if userData != nil {
		// Check the verification data first
		userIDString := fmt.Sprint(userData.ID)
		verificationData, _ := verification.GetVerificationByID(userIDString, "Forgot Password", config.DB)
		if verificationData != nil {
			// Delete the existing verification data
			verificationIDString := fmt.Sprint(verificationData.ID)
			_, err := verification.DeleteVerification(verificationIDString, config.DB)
			if err != nil {
				lib.Error(response, http.StatusBadRequest, err.Error())
				return
			}
		}

		// Generate the new verification data
		randomString := lib.RandStringRunes(30)
		verification.Name = "Forgot Password"
		verification.Token = randomString
		verification.UserID = userData.ID
		config.DB.Save(&verification)

		lib.Success(response, http.StatusCreated, "Verification token has been sent", verification)
		return
	}

	lib.Error(response, http.StatusNotFound, "User not found")
	return
}

// ChangePassword user
func ChangePassword(response http.ResponseWriter, request *http.Request) {
	verification := &models.Verification{}
	user := &models.User{}
	userJSON := &models.UserJSON{}
	verificationToken := mux.Vars(request)["token"]

	// Check verification token
	if verificationToken == "" {
		lib.Error(response, http.StatusNotFound, "Verification token not found")
		return
	}

	// Get verification data by token
	verificationData, _ := verification.GetVerificationByToken(verificationToken, config.DB)
	if verificationData == nil {
		lib.Error(response, http.StatusNotFound, "Verification data not found")
		return
	}

	// Get all request body
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the request
	err = json.Unmarshal(body, &user)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the user input
	err = user.Validate("change-password")
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Get the user data
	userIDString := fmt.Sprint(verificationData.UserID)
	userData, _ := userJSON.GetUser(userIDString, config.DB)
	if userData == nil {
		lib.Error(response, http.StatusNotFound, "User not found")
		return
	}

	// Update the user password
	_, err = user.ChangeUserPassword(userIDString, config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Delete verification data
	verificationIDString := fmt.Sprint(verificationData.ID)
	_, err = verification.DeleteVerification(verificationIDString, config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	lib.Success(response, http.StatusOK, "Password successfully changed", userData)
	return
}
