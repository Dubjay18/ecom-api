package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Dubjay18/ecom-api/internal/config"
	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/Dubjay18/ecom-api/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	r      *gin.RouterGroup
	s      *service.ProductService
	logger *logrus.Logger
	cf     config.APIKeysConfig
}

func NewProductHandler(r *gin.RouterGroup, s *service.ProductService, logger *logrus.Logger, secretKey string, cfg config.APIKeysConfig) {
	handler := &ProductHandler{
		r:      r,
		s:      s,
		logger: logger,
		cf:     cfg,
	}
	r.Use(middleware.AuthMiddleware(secretKey))
	ar := r.Use(middleware.AdminMiddleware())

	ar.GET("/products", handler.ListProducts)
	ar.POST("/products", handler.CreateProduct)
	ar.GET("/products/:id", handler.GetProduct)
	ar.PUT("/products/:id", handler.UpdateProduct)
	ar.DELETE("/products/:id", handler.DeleteProduct)
}

// Create Product godoc
// @Summary Create a new product
// @Description Create a new product using form data
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product name"
// @Param price formData number true "Product price"
// @Param stock formData int true "Product stock"
// @Param sku formData string true "Product SKU"
// @Param category formData string true "Product category"
// @Param image formData file true "Product image"
// @Success 201 {object} domain.Product
// @Failure 400 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest

	if err := c.ShouldBind(&req); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			response.RenderBindingErrors(c, err.(validator.ValidationErrors))
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	log.Println("Received form data:", req)

	file, err := c.FormFile("image")
	if err != nil {
		h.logger.Error(err.Error())
		response.Error(c, http.StatusBadRequest, "Image file is required", "image file is required")
		return
	}

	imagePath, err := upload.UploadImage(c, file, h.cf.CloudinaryCloudName, h.cf.CloudinaryKey, h.cf.CloudinarySecret)
	if err != nil {
		h.logger.Error(err.Error())
		response.Error(c, http.StatusBadRequest, "Failed to upload image", err.Error())
		return
	}

	product := &domain.Product{
		Name:     req.Name,
		Price:    req.Price,
		Stock:    req.Stock,
		SKU:      req.SKU,
		Category: req.Category,
		ImageURL: imagePath,
	}

	perr := h.s.Create(c.Request.Context(), product)
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

// Get Product godoc
// @Summary Get a product by ID
// @Description Returns the product that matches the given ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	product, perr := h.s.GetByID(c.Request.Context(), uint(id))
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

// Update Product godoc
// @Summary Update a product
// @Description Updates a product by ID using form data
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param name formData string false "Product name"
// @Param price formData number false "Product price"
// @Param stock formData int false "Product stock"
// @Param sku formData string false "Product SKU"
// @Param category formData string false "Product category"
// @Param image formData file false "Product image"
// @Success 200 {object} domain.Product
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var req domain.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			response.RenderBindingErrors(c, err.(validator.ValidationErrors))
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	existingProduct, perr := h.s.GetByID(c.Request.Context(), uint(id))
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Price != 0 {
		existingProduct.Price = req.Price
	}
	if req.Stock != 0 {
		existingProduct.Stock = req.Stock
	}
	if req.SKU != "" {
		existingProduct.SKU = req.SKU
	}
	if req.Category != "" {
		existingProduct.Category = req.Category
	}

	file, err := c.FormFile("image")
	if err == nil {
		imagePath, err := upload.UploadImage(c, file, h.cf.CloudinaryCloudName, h.cf.CloudinaryKey, h.cf.CloudinarySecret)
		if err != nil {
			h.logger.Error(err.Error())
			response.Error(c, http.StatusBadRequest, "Failed to upload image", err.Error())
			return
		}
		existingProduct.ImageURL = imagePath
	}

	perr = h.s.Update(c.Request.Context(), existingProduct)
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", existingProduct)
}

// Delete Product godoc
// @Summary Delete a product
// @Description Deletes a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	perr := h.s.Delete(c.Request.Context(), uint(id))
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

// List Products godoc
// @Summary List products
// @Description Lists products with optional filtering
// @Tags products
// @Accept json
// @Produce json
// @Param name query string false "Product name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {array} domain.Product
// @Failure 400 {object} response.Response
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var filter domain.ProductFilter
	filter.Name = c.Query("name")
	filter.MinPrice = parseFloat64(c.Query("min_price"))
	filter.MaxPrice = parseFloat64(c.Query("max_price"))

	products, perr := h.s.List(c.Request.Context(), filter)
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}

func parseFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
