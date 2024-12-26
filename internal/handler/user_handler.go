package handler

import (
	"net/http"

	"github.com/Dubjay18/ecom-api/internal/domain"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/Dubjay18/ecom-api/pkg/common"
	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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
		auth.POST("/register-admin", handler.RegisterAdmin)
	}

	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtSecret))
	{
		users.GET("/me", handler.GetProfile)
		// users.PUT("/me", handler.UpdateProfile)
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.RegisterRequest true "User registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
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
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
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
// @Description Get authenticated user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} response.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.s.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err.Code, err.Message, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// RegisterAdmin godoc
// @Summary Register new admin
// @Description Register a new admin with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param admin body domain.RegisterRequest true "Admin registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register-admin [post]
func (h *UserHandler) RegisterAdmin(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			response.RenderBindingErrors(c, err.(validator.ValidationErrors))
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	user, err := h.s.RegisterAdmin(c.Request.Context(), req)
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
