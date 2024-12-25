package handler

import (
	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	r      *gin.RouterGroup
	s      service.UserService
	logger *logrus.Logger
}

func NewUserHandler(r *gin.RouterGroup, s service.UserService, logger *logrus.Logger, jwtSecret string) {
	handler := &UserHandler{
		r:      r,
		s:      s,
		logger: logger,
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtSecret))
	{
		users.GET("/me", handler.GetProfile)
		// users.PUT("/me", handler.UpdateProfile)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.RegisterRequest true "User registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			response.RenderBindingErrors(c, err.(validator.ValidationErrors))
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	user, err := h.s.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Error(err)
		switch err {
		case &common.ErrEmailExists:
			response.Error(c, http.StatusConflict, "Registration failed", err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, "Failed to register user", err.Error())
		}
		return

	}
	response.Success(c, http.StatusCreated, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	resp, err := h.s.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err.Code, "Invalid credentials", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.s.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err.Code, err.Message, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
