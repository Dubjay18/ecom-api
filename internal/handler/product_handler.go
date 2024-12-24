package handler

import (
	"net/http"
	"strconv"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	r *gin.RouterGroup
	s *service.ProductService
}

func NewProductHandler(r *gin.RouterGroup, s *service.ProductService) {
	handler := &ProductHandler{
		r: r,
		s: s,
	}
	ar := r.Use(middleware.AdminMiddleware())
	ar.GET("/products", handler.ListProducts)
	ar.POST("/products", handler.CreateProduct)
	ar.GET("/products/:id", handler.GetProduct)
	ar.PUT("/products/:id", handler.UpdateProduct)
	ar.DELETE("/products/:id", handler.DeleteProduct)
}

// Create Product godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product details"
// @Success 201 {object} domain.Product
// @Failure 400 {object} response.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	product := &domain.Product{
		Name:       form.Value["name"][0],
		Price:      parseFloat64(form.Value["price"][0]),
		Stock:      parseInt(form.Value["stock"][0]),
		SKU:        form.Value["sku"][0],
		CategoryID: uint(parseInt(form.Value["category_id"][0])),
	}

	// Handle image file
	file, err := c.FormFile("image")
	if err == nil {
		// If image is provided, process it
		// You might want to implement proper file storage logic here
		product.ImageURL = file.Filename
	}

	perr := h.s.Create(c.Request.Context(), product)
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func parseFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

// Get Product godoc
// @Summary Get a product by ID
// @Description Get a product by ID
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
// @Description Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body domain.Product true "Product details"
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

	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}
	product, perr := h.s.GetByID(c.Request.Context(), uint(id))
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	if form.Value["name"] != nil {
		product.Name = form.Value["name"][0]
	}
	if form.Value["price"] != nil {
		product.Price = parseFloat64(form.Value["price"][0])
	}
	if form.Value["stock"] != nil {
		product.Stock = parseInt(form.Value["stock"][0])
	}
	if form.Value["sku"] != nil {
		product.SKU = form.Value["sku"][0]
	}
	if form.Value["category_id"] != nil {
		product.CategoryID = uint(parseInt(form.Value["category_id"][0]))
	}
	if form.Value["description"] != nil {
		product.Description = form.Value["description"][0]
	}
	if form.Value["image_url"] != nil {
		product.ImageURL = form.Value["image_url"][0]
	}

	// // Handle image file
	// file, err := c.FormFile("image")
	// if err == nil {
	// 	// If image is provided, process it
	// 	// You might want to implement proper file storage logic here
	// 	product.ImageURL = file.Filename
	// }

	perr = h.s.Update(c.Request.Context(), product)
	if perr != nil {
		response.Error(c, perr.Code, perr.Message, perr.Error())
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

// Delete Product godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
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

	response.Success(c, http.StatusNoContent, "Product deleted successfully", nil)
}

// List Products godoc
// @Summary List products
// @Description List products
// @Tags products
// @Accept json
// @Produce json
// @Param name query string false "Product name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {object} []domain.Product
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
