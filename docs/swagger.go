package docs

// @title           E-commerce API
// @version         1.0
// @description     A RESTful API for an e-commerce application

// @contact.name   Jay
// @contact.email  jejeniyi7@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// swagger definitions for User endpoints
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User registration details"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
