# 📚 API Documentation

This document provides detailed information about all available API endpoints, request/response formats, and usage examples.

## 🔗 Base URL

```
http://localhost:8080/api/v1
```

## 🔐 Authentication

All authenticated endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 📋 Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "status": 200,
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "status": 400,
  "message": "Error message",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ]
}
```

## 🔐 Authentication Endpoints

### Register User

Creates a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "status": 201,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Validation Rules:**
- `email`: Required, valid email format, unique
- `password`: Required, minimum 6 characters
- `first_name`: Required
- `last_name`: Required

---

### Login User

Authenticates a user and returns a JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "status": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user"
    }
  }
}
```

---

### Register Admin

Creates a new admin user account.

**Endpoint:** `POST /auth/register-admin`

**Request Body:**
```json
{
  "email": "admin@example.com",
  "password": "securepassword",
  "first_name": "Admin",
  "last_name": "User"
}
```

**Response:** Same as regular registration but with `role: "admin"`

## 👥 User Endpoints

### Get Current User Profile

Returns the profile of the authenticated user.

**Endpoint:** `GET /users/me`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "status": 200,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

## 🛍️ Product Endpoints

All product endpoints require admin authentication.

### List Products

Retrieves a list of all products with optional filtering.

**Endpoint:** `GET /products`

**Headers:** `Authorization: Bearer <admin-token>`

**Query Parameters:**
- `name` (optional): Filter by product name
- `min_price` (optional): Minimum price filter
- `max_price` (optional): Maximum price filter

**Example:** `GET /products?name=laptop&min_price=500&max_price=2000`

**Response:**
```json
{
  "status": 200,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Gaming Laptop",
      "description": "High-performance gaming laptop",
      "price": 1299.99,
      "sku": "LAPTOP-001",
      "stock": 10,
      "category": "Electronics",
      "image_url": "https://example.com/laptop.jpg",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### Create Product

Creates a new product with optional image upload.

**Endpoint:** `POST /products`

**Headers:** 
- `Authorization: Bearer <admin-token>`
- `Content-Type: multipart/form-data`

**Form Data:**
- `name`: Product name (required)
- `description`: Product description (optional)
- `price`: Product price (required, > 0)
- `sku`: Stock Keeping Unit (required, unique)
- `stock`: Stock quantity (required, > 0)
- `category`: Product category (optional)
- `image`: Image file (optional)

**Response:**
```json
{
  "status": 201,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "price": 1299.99,
    "sku": "LAPTOP-001",
    "stock": 10,
    "category": "Electronics",
    "image_url": "https://cloudinary.com/image/upload/v1234567890/laptop.jpg",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### Get Product by ID

Retrieves a specific product by its ID.

**Endpoint:** `GET /products/{id}`

**Headers:** `Authorization: Bearer <admin-token>`

**Response:**
```json
{
  "status": 200,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "price": 1299.99,
    "sku": "LAPTOP-001",
    "stock": 10,
    "category": "Electronics",
    "image_url": "https://cloudinary.com/image/upload/v1234567890/laptop.jpg",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### Update Product

Updates an existing product.

**Endpoint:** `PUT /products/{id}`

**Headers:** 
- `Authorization: Bearer <admin-token>`
- `Content-Type: multipart/form-data`

**Form Data:** (all optional)
- `name`: Product name
- `description`: Product description
- `price`: Product price (> 0)
- `sku`: Stock Keeping Unit
- `stock`: Stock quantity (> 0)
- `category`: Product category
- `image`: New image file

**Response:** Same as create product response with updated data.

---

### Delete Product

Deletes a product.

**Endpoint:** `DELETE /products/{id}`

**Headers:** `Authorization: Bearer <admin-token>`

**Response:**
```json
{
  "status": 200,
  "message": "Product deleted successfully"
}
```

## 📦 Order Endpoints

### Place Order

Creates a new order for the authenticated user.

**Endpoint:** `POST /orders`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 2,
      "quantity": 1
    }
  ],
  "shipping_address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "country": "USA",
    "postal_code": "10001",
    "is_default": true
  },
  "payment_method": "credit_card"
}
```

**Response:**
```json
{
  "status": 201,
  "message": "Order created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "status": "pending",
    "total_amount": 2599.98,
    "payment_status": "pending",
    "shipping_address_id": 1,
    "items": [
      {
        "id": 1,
        "quantity": 2,
        "price": 1299.99,
        "product": {
          "id": 1,
          "name": "Gaming Laptop",
          "sku": "LAPTOP-001"
        }
      }
    ],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

### Get User Orders

Retrieves all orders for the authenticated user.

**Endpoint:** `GET /orders`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "status": 200,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "status": "pending",
      "total_amount": 2599.98,
      "payment_status": "pending",
      "shipping_address_id": 1,
      "items": [
        {
          "id": 1,
          "quantity": 2,
          "price": 1299.99,
          "product": {
            "id": 1,
            "name": "Gaming Laptop",
            "sku": "LAPTOP-001"
          }
        }
      ],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### Cancel Order

Cancels an order (only if status is "pending").

**Endpoint:** `DELETE /orders/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "status": 200,
  "message": "Order cancelled successfully"
}
```

## 🏥 Health Check

### Health Check

Returns the health status of the API.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h30m45s"
}
```

## 📊 HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Invalid or missing authentication |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation errors |
| 500 | Internal Server Error - Server error |

## 🔍 Error Handling

### Common Error Responses

**Validation Error (400):**
```json
{
  "status": 400,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    },
    {
      "field": "password",
      "message": "Password must be at least 6 characters"
    }
  ]
}
```

**Authentication Error (401):**
```json
{
  "status": 401,
  "message": "Unauthorized",
  "error": "Invalid or expired token"
}
```

**Permission Error (403):**
```json
{
  "status": 403,
  "message": "Forbidden",
  "error": "Admin access required"
}
```

**Not Found Error (404):**
```json
{
  "status": 404,
  "message": "Resource not found",
  "error": "Product with ID 999 not found"
}
```

## 📝 Examples

### Complete Order Flow

1. **Register a user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123",
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

2. **Login and get token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

3. **Place an order:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 1
      }
    ],
    "shipping_address": {
      "street": "456 Oak Ave",
      "city": "Los Angeles",
      "state": "CA",
      "country": "USA",
      "postal_code": "90210"
    },
    "payment_method": "credit_card"
  }'
```

### Admin Product Management

1. **Register admin:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register-admin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "adminpassword",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

2. **Create product with image:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <admin-token>" \
  -F "name=Wireless Mouse" \
  -F "description=Ergonomic wireless mouse" \
  -F "price=29.99" \
  -F "sku=MOUSE-001" \
  -F "stock=50" \
  -F "category=Electronics" \
  -F "image=@mouse.jpg"
```

---

For more information, visit the [Swagger UI](http://localhost:8080/swagger/index.html) for interactive API documentation.