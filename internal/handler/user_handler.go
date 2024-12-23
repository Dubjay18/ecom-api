package handler

import (
	"github.com/Dubjay18/ecom-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	r *gin.RouterGroup
	s *service.UserService
}

func NewUserHandler(r *gin.RouterGroup, s *service.UserService) *UserHandler {
	return &UserHandler{r: r, s: s}
}
