package controllers

import (
	"chat-backend/app/lib"
	"chat-backend/app/models"
	"chat-backend/config"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

// GetAuthenticatedUser getting one user
func GetAuthenticatedUser(response http.ResponseWriter, request *http.Request) {
	user := &models.UserJSON{}

	userIDFromToken := fmt.Sprint(context.Get(request, "UserID"))
	userData, _ := user.GetUser(userIDFromToken, config.DB)
	if userData == nil {
		lib.Error(response, http.StatusBadRequest, "User not found")
		return
	}

	lib.Success(response, http.StatusOK, "Hi "+userData.Name, userData)
	return
}
