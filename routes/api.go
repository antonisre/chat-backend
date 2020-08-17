package routes

import (
	"chat-backend/app/controllers"
	"chat-backend/app/middlewares"
	socket "chat-backend/app/socket"
	"chat-backend/config"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	Router *mux.Router
}

// ServeRoutes handle the public routes
func (api *Api) ServeRoutes() {
	api.Router = mux.NewRouter()

	// Server static file
	var imgServer = http.FileServer(http.Dir(config.UploadPath))
	api.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", imgServer))

	//init websocket
	api.Router.HandleFunc("/ws", socket.SocketEndPoint)

	// Route List
	PublicRouter := api.Router.PathPrefix("/api").Subrouter()
	ProtectedRouter := api.Router.PathPrefix("/api/v1").Subrouter()
	ProtectedRouterHighAdminRouter := api.Router.PathPrefix("/api/v1").Subrouter()

	// Middleware
	PublicRouter.Use(middlewares.SetContentTypeHeader)
	ProtectedRouterHighAdminRouter.Use(middlewares.SetContentTypeHeader)
	ProtectedRouter.Use(middlewares.SetContentTypeHeader)
	ProtectedRouter.Use(middlewares.AuthJwtVerify)
	ProtectedRouterHighAdminRouter.Use(middlewares.AuthJwtVerify)
	ProtectedRouterHighAdminRouter.Use(middlewares.OnlyHighAdmin)

	// Open Routes
	PublicRouter.HandleFunc("/register", controllers.Register).Methods("POST")
	PublicRouter.HandleFunc("/login", controllers.Login).Methods("POST")
	PublicRouter.HandleFunc("/forgot-password", controllers.ForgotPassword).Methods("POST")
	PublicRouter.HandleFunc("/change-password/{token}", controllers.ChangePassword).Methods("PATCH")

	// High Admin Routes
	ProtectedRouterHighAdminRouter.HandleFunc("/roles", controllers.GetAllRoles).Methods("GET")
	ProtectedRouterHighAdminRouter.HandleFunc("/roles", controllers.CreateRole).Methods("POST")
	ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.GetRole).Methods("GET")
	ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.UpdateRole).Methods("PATCH")
	ProtectedRouterHighAdminRouter.HandleFunc("/roles/{id}", controllers.DeleteRole).Methods("DELETE")

	// Protected Routes
	ProtectedRouter.HandleFunc("/users/me", controllers.GetAuthenticatedUser).Methods("GET")
	ProtectedRouter.HandleFunc("/users/me/upload-image", controllers.UploadUserImage).Methods("PATCH")
	ProtectedRouter.HandleFunc("/users/me/delete-image", controllers.DeleteImage).Methods("DELETE")
	ProtectedRouter.HandleFunc("/users/me/delete-image", controllers.DeleteImage).Methods("DELETE")
	ProtectedRouter.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
}
