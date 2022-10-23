package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	authdto "waysfood/dto/auth"
	dto "waysfood/dto/result"
	usersdto "waysfood/dto/users"
	"waysfood/models"
	"waysfood/pkg/bcrypt"
	jwtToken "waysfood/pkg/jwt"
	"waysfood/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(usersdto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	user := models.User{
		Email:    request.Email,
		Password: password,
		FullName: request.FullName,
		Gender:   request.Gender,
		Phone:    request.Phone,
		Role:     request.Role,
	}

	data, _ := h.AuthRepository.Register(user)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.AuthRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := h.AuthRepository.Login(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Email not registered!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Wrong password!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	gnrtToken := jwt.MapClaims{}
	gnrtToken["id"] = user.ID
	gnrtToken["exp"] = time.Now().Add(time.Hour * 3).Unix()

	token, err := jwtToken.GenerateToken(&gnrtToken)
	if err != nil {
		log.Println(err)
		fmt.Println("Unauthorize")
		return
	}

	AuthResponse := authdto.AuthResponse{
		FullName: user.FullName,
		Email:    user.Email,
		Password: user.Password,
		Token:    token,
	}

	w.Header().Set("Content-Type", "application/json")
	response := dto.SuccessResult{Status: http.StatusOK, Data: AuthResponse}
	json.NewEncoder(w).Encode(response)

}
