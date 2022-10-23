package routes

import (
	"waysfood/handlers"
	"waysfood/pkg/mysql"
	"waysfood/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	r.HandleFunc("/users", h.ShowUsers).Methods("GET")
	r.HandleFunc("/user/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/user/{id}", h.UpdateUser).Methods("PATCH")
	r.HandleFunc("/user/{id}", h.DeleteUser).Methods("DELETE")
}
