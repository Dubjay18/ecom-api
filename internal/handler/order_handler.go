package handler

import (
	"net/http"
	"strconv"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
// @Summary Place an order for one or more products
// @Description Create a new order for authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body domain.CreateOrderRequest true "Order details"
// @Success 201 {object} domain.Order
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RenderBindingErrors(c, err.(validator.ValidationErrors))
		return
	}

	userID, _ := c.Get("userID")
	order, err := h.s.PlaceOrder(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.Error(c, err.Code, "Failed to create order", err.Message)
		return
	}

	response.Success(c, http.StatusCreated, "Order created successfully", order)
}

// GetUserOrders godoc
// @Summary List all orders for a user
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Order
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orders, err := h.s.ListUserOrders(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, err.Code, "Failed to fetch orders", err.Message)
		return
	}

	response.Success(c, http.StatusOK, "Orders retrieved successfully", orders)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order if it is still pending
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/orders/:id [delete]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	oerr := h.s.CancelOrder(c.Request.Context(), uint(id))
	if oerr != nil {
		response.Error(c, oerr.Code, "Failed to cancel order", oerr.Message)
		return
	}

	response.Success(c, http.StatusOK, "Order cancelled successfully", nil)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param status body domain.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} domain.Order
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/orders/:id/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req domain.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid order ID", err.Error())
		return
	}

	oerr := h.s.UpdateOrderStatus(c.Request.Context(), uint(id), req.Status)
	if oerr != nil {
		response.Error(c, oerr.Code, "Failed to update order status", oerr.Message)
		return
	}

	response.Success(c, http.StatusOK, "Order status updated successfully", nil)
}

// RegisterRoutes registers order-related routes
func (h *OrderHandler) RegisterRoutes() {
	h.r.POST("/orders", h.CreateOrder)
	h.r.GET("/orders", h.GetUserOrders)
	h.r.DELETE("/orders/:id", h.CancelOrder)

	adminRoutes := h.r.Group("/")
	adminRoutes.Use(middleware.AdminMiddleware())
	adminRoutes.PUT("/orders/:id/status", h.UpdateOrderStatus)
}

// SuccessResponse godoc
type SuccessResponse struct {
	Message string `json:"message"`
}
