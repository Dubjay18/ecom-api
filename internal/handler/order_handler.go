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

func NewOrderHandler(r *gin.RouterGroup, s *service.OrderService, secretKey string) *OrderHandler {
	handler := &OrderHandler{
		r: r,
		s: s,
	}
	r.Use(middleware.AuthMiddleware(secretKey))
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
// @Example JSON Response - Success
//
//	{
//	  "status": 201,
//	  "message": "Order created successfully",
//	  "data": {
//	    "id": 1,
//	    "user_id": 42,
//	    "status": "pending",
//	    "items": [
//	      {
//	        "product": {..},
//	        "quantity": 2,
//	        "price": 29.99
//	      }
//	    ]
//	  }
//	}
//
// @Failure 400 {object} response.ErrorResponse
// @Example JSON Response - Error
//
//	{
//	  "status": 400,
//	  "message": "Invalid input",
//	  "error": [
//	    {"items": "required"}
//	  ]
//	}
//
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
// @Example JSON Response - Success
//
//	{
//	  "status": 200,
//	  "message": "Orders retrieved successfully",
//	  "data": [
//	    {
//	      "id": 1,
//	      "status": "delivered",
//	      "items": []
//	    }
//	  ]
//	}
//
// @Failure 401 {object} response.ErrorResponse
// @Example JSON Response - Error
//
//	{
//	  "status": 401,
//	  "message": "Unauthorized",
//	  "error": "token is invalid"
//	}
//
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	// ...
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order if it is still pending
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} response.Response
// @Example JSON Response - Success
//
//	{
//	  "status": 200,
//	  "message": "Order cancelled successfully"
//	}
//
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Example JSON Response - Error
//
//	{
//	  "status": 404,
//	  "message": "Failed to cancel order",
//	  "error": "Order not found"
//	}
//
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
// @Security Bearer
// @Security JWT
// @Param id path int true "Order ID"
// @Param status body domain.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} domain.Order
// @Example JSON Response - Success
//
//	{
//	  "status": 200,
//	  "message": "Order status updated successfully",
//	  "data": {
//	    "id": 1,
//	    "status": "shipped"
//	  }
//	}
//
// @Failure 400 {object} response.ErrorResponse
// @Example JSON Response - Error
//
//	{
//	  "status": 400,
//	  "message": "Invalid order ID",
//	  "error": " invalid syntax"
//	}
//
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
