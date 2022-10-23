package repositories

import (
	"waysfood/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	ShowProducts() ([]models.Product, error)
	GetProductByID(ID int) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(product models.Product, ID int) (models.Product, error)
	DeleteProduct(product models.Product, ID int) (models.Product, error)
}

func RepositoryProduct(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ShowProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("User").Find(&products).Error

	return products, err
}

func (r *repository) GetProductByID(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("User").First(&product, ID).Error

	return product, err
}

func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
	err := r.db.Create(&product).Error

	return product, err
}

func (r *repository) UpdateProduct(product models.Product, ID int) (models.Product, error) {
	err := r.db.Model(&product).Where("id=?", ID).Updates(&product).Error

	return product, err
}

func (r *repository) DeleteProduct(Product models.Product, ID int) (models.Product, error) {
	err := r.db.Delete(&Product).Error

	return Product, err
}
