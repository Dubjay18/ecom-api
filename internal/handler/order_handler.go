package handler

import (
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	r *gin.RouterGroup
	s *service.OrderService
}

func NewOrderHandler(r *gin.RouterGroup, s *service.OrderService) *OrderHandler {
	return &OrderHandler{r: r, s: s}
}
