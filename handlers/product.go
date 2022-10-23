package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	productdto "waysfood/dto/product"
	dto "waysfood/dto/result"
	"waysfood/models"
	"waysfood/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var path_file = "http://localhost:8080/uploads/"

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) ShowProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.ProductRepository.ShowProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	for i, p := range products {
		products[i].Image = path_file + p.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: products}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var product models.Product
	product, err := h.ProductRepository.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product.Image = path_file + product.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	dataUpload := r.Context().Value("dataFile")
	filename := dataUpload.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	qty, _ := strconv.Atoi(r.FormValue("qty"))

	request := productdto.CreateProductRequest{
		Title: r.FormValue("title"),
		Price: price,
		Qty:   qty,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product := models.Product{
		Title:  request.Title,
		Price:  request.Price,
		Image:  filename,
		Qty:    request.Qty,
		UserID: userId,
	}

	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product, _ = h.ProductRepository.GetProductByID(product.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataContext := r.Context().Value("dataFile")
	filename := dataContext.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	qty, _ := strconv.Atoi(r.FormValue("qty"))
	request := productdto.UpdateProductRequest{
		Title: r.FormValue("title"),
		Price: price,
		Qty:   qty,
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	product := models.Product{}

	product.ID = id

	if request.Title != "" {
		product.Title = request.Title
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if request.Image != "" {
		product.Image = filename
	}

	if request.Qty != 0 {
		product.Qty = request.Qty
	}

	data, err := h.ProductRepository.UpdateProduct(product, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	product, err := h.ProductRepository.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ProductRepository.DeleteProduct(product, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}
