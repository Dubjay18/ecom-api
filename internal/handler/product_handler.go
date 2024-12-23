package handler

import (
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	r *gin.RouterGroup
	s *service.ProductService
}

func NewProductHandler(r *gin.RouterGroup, s *service.ProductService) *ProductHandler {
	return &ProductHandler{r: r, s: s}
}
