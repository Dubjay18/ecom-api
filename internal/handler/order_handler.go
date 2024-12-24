package handler

import (
	"net/http"
	"strconv"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	r *gin.RouterGroup
	s *service.OrderService
}

func NewOrderHandler(r *gin.RouterGroup, s *service.OrderService) *OrderHandler {

	handler := &OrderHandler{
		r: r,
		s: s,
	}

	handler.RegisterRoutes()

	return handler
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body CreateOrderRequest true "Order creation details"
// @Success 201 {object} Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	order, err := h.s.Create(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.Error(c, err.Code, "Failed to create order", err.Message)
		return
	}

	response.Success(c, http.StatusCreated, "Order Created Successfully", order)
}

// GetUserOrders godoc
// @Summary Get user orders
// @Description Get all orders for authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Order
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/orders [get]

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orders, err := h.s.List(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, err.Code, "Failed to fetch orders", err.Message)
		return
	}

	response.Success(c, http.StatusOK, "User Orders", orders)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}
	order, oerr := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, oerr.Code, "Failed to fetch order", oerr.Message)
		return
	}

	response.Success(c, http.StatusOK, "Order", order)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	oerr := h.s.Cancel(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, oerr.Code, "Failed to cancel order", oerr.Message)
		return
	}

	response.Success(c, http.StatusOK, "Order cancelled successfully", nil)
}

// RegisterRoutes registers routes
func (h *OrderHandler) RegisterRoutes() {
	h.r.POST("/orders", h.CreateOrder)
	h.r.GET("/orders", h.GetUserOrders)
	h.r.GET("/orders/:id", h.GetOrder)
	h.r.DELETE("/orders/:id", h.CancelOrder)
}

// SuccessResponse godoc
type SuccessResponse struct {
	Message string `json:"message"`
}
