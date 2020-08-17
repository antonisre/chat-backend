package controllers

import (
	"chat-backend/config"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"chat-backend/app/lib"
	"chat-backend/app/models"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllRoles getting all users
func GetAllRoles(response http.ResponseWriter, request *http.Request) {
	role := &models.Role{}

	// Count total of roles
	total, err := role.CountRoles(config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Paginate the roles
	queryParams := request.URL.Query()
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	nameParam := queryParams.Get("name")
	if limit < 1 {
		limit = 10
	}
	page, begin := lib.Pagination(request, limit)
	// @info total variable's from counting the roles in the model
	pages := total / limit
	if (total % limit) != 0 {
		pages++
	}

	// Return the paginate
	roles, err := role.GetRoles(begin, limit, nameParam, config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	mapRoles := lib.PaginationResponse(request, page, pages, limit, total, roles)

	lib.Success(response, http.StatusOK, "Roles list", mapRoles)
	return
}

// CreateRole create a new role
func CreateRole(response http.ResponseWriter, request *http.Request) {
	role := &models.Role{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		lib.Error(response, http.StatusUnprocessableEntity, err.Error())
		return
	}

	newRole, err := role.Create(config.DB)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	lib.Success(response, http.StatusCreated, "Role successfully created", newRole)
}

// GetRole get role by ID
func GetRole(response http.ResponseWriter, request *http.Request) {
	role := &models.Role{}
	id := mux.Vars(request)["id"]

	roleData, _ := role.GetRoleByID(id, config.DB)

	if roleData == nil {
		lib.Error(response, http.StatusNotFound, "Role not found")
		return
	}

	lib.Success(response, http.StatusOK, "Role Detail", roleData)
}

// UpdateRole selected role
func UpdateRole(response http.ResponseWriter, request *http.Request) {
	role := &models.Role{}
	id := mux.Vars(request)["id"]

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		lib.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		lib.Error(response, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Check role data and update the data
	roleData, _ := role.GetRoleByID(id, config.DB)
	if roleData != nil {
		if err := config.DB.Debug().Table("roles").First(&roleData).Update("name", role.Name).Error; err != nil {
			lib.Error(response, http.StatusBadRequest, err.Error())
			return
		}
		config.DB.Save(&roleData)

		lib.Success(response, http.StatusOK, "Role successfully updated", roleData)
		return
	}

	lib.Error(response, http.StatusNotFound, "Role not found")
	return
}

// DeleteRole delete selected role
func DeleteRole(response http.ResponseWriter, request *http.Request) {
	role := &models.Role{}
	id := mux.Vars(request)["id"]

	roleData, _ := role.GetRoleByID(id, config.DB)
	if roleData != nil {
		_, err := role.Delete(roleData.ID, config.DB)
		if err != nil {
			lib.Error(response, http.StatusBadRequest, err.Error())
			return
		}
		lib.Success(response, http.StatusOK, "Role successfully deleted", roleData)
		return
	}

	lib.Error(response, http.StatusNotFound, "Role not found")
	return
}
