package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	dto "waysfood/dto/result"
	usersdto "waysfood/dto/users"
	"waysfood/models"
	"waysfood/pkg/bcrypt"
	"waysfood/repositories"

	"github.com/gorilla/mux"
)

type handler struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handler {
	return &handler{UserRepository}
}

func (h *handler) ShowUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepository.ShowUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.UserRepository.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Status: http.StatusNotFound, Message: "ID: " + strconv.Itoa(id) + " not found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	UserResponse := usersdto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Location: user.Location,
		Image:    user.Image,
		Role:     user.Role,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: UserResponse}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(usersdto.UpdateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user := models.User{}

	user.ID = id

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Password != "" {
		password, err := bcrypt.HashingPassword(request.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
		}
		user.Password = password
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}

	if request.Location != "" {
		user.Location = request.Location
	}

	if request.Image != "" {
		user.Image = request.Image
	}

	if request.Gender != "" {
		user.Gender = request.Gender
	}

	data, err := h.UserRepository.UpdateUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := h.UserRepository.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.UserRepository.DeleteUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: convertResponse(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponse(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		FullName: u.FullName,
		Gender:   u.Gender,
		Phone:    u.Phone,
		Location: u.Location,
		Role:     u.Role,
		Image:    u.Image,
	}
}
