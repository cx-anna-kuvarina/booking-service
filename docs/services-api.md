# Services API Documentation

## Overview

The Services API provides comprehensive management of business services, including creation, reading, updating, and deletion of services offered by business accounts.

## Base URL

```
/api/services
```

## Authentication

All endpoints require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### 1. Create Service

**POST** `/api/services/`

Creates a new service for a business account.

**Request Body:**
```json
{
  "business_account_id": "uuid-string",
  "name": "Service Name",
  "description": "Optional service description",
  "duration_minutes": 60,
  "price": 50.00,
  "currency": "USD",
  "category": "Optional category"
}
```

**Required Fields:**
- `business_account_id`: UUID of the business account
- `name`: Service name (string)
- `duration_minutes`: Service duration in minutes (positive integer)
- `price`: Service price (non-negative number)

**Optional Fields:**
- `description`: Service description
- `currency`: Currency code (3 characters, defaults to "USD")
- `category`: Service category

**Response:**
```json
{
  "id": "uuid-string",
  "business_account_id": "uuid-string",
  "name": "Service Name",
  "description": "Optional service description",
  "duration_minutes": 60,
  "price": 50.00,
  "currency": "USD",
  "category": "Optional category",
  "is_active": true,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**Status Codes:**
- `201 Created`: Service created successfully
- `400 Bad Request`: Validation error
- `404 Not Found`: Business account not found
- `500 Internal Server Error`: Server error

### 2. Get Service

**GET** `/api/services/{id}`

Retrieves a specific service by ID.

**Path Parameters:**
- `id`: Service UUID

**Response:**
```json
{
  "id": "uuid-string",
  "business_account_id": "uuid-string",
  "name": "Service Name",
  "description": "Service description",
  "duration_minutes": 60,
  "price": 50.00,
  "currency": "USD",
  "category": "Category",
  "is_active": true,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**Status Codes:**
- `200 OK`: Service retrieved successfully
- `400 Bad Request`: Missing service ID
- `404 Not Found`: Service not found
- `500 Internal Server Error`: Server error

### 3. Update Service

**PUT** `/api/services/{id}`

Updates an existing service. Only provided fields will be updated.

**Path Parameters:**
- `id`: Service UUID

**Request Body:**
```json
{
  "name": "Updated Service Name",
  "price": 75.00,
  "is_active": false
}
```

**Updatable Fields:**
- `name`: Service name
- `description`: Service description
- `duration_minutes`: Service duration in minutes
- `price`: Service price
- `currency`: Currency code
- `category`: Service category
- `is_active`: Service availability status

**Response:**
```json
{
  "id": "uuid-string",
  "business_account_id": "uuid-string",
  "name": "Updated Service Name",
  "description": "Service description",
  "duration_minutes": 60,
  "price": 75.00,
  "currency": "USD",
  "category": "Category",
  "is_active": false,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T11:00:00Z"
}
```

**Status Codes:**
- `200 OK`: Service updated successfully
- `400 Bad Request`: Validation error
- `404 Not Found`: Service not found
- `500 Internal Server Error`: Server error

### 4. Delete Service

**DELETE** `/api/services/{id}`

Deletes a service permanently.

**Path Parameters:**
- `id`: Service UUID

**Response:**
- No content

**Status Codes:**
- `204 No Content`: Service deleted successfully
- `400 Bad Request`: Missing service ID
- `404 Not Found`: Service not found
- `500 Internal Server Error`: Server error

### 5. List Services

**GET** `/api/services/`

Retrieves a paginated list of services with optional filtering.

**Query Parameters:**
- `limit`: Maximum number of services to return (default: 20, max: 100)
- `offset`: Number of services to skip (default: 0)
- `business_account_id`: Filter by business account ID
- `category`: Filter by service category
- `is_active`: Filter by active status (true/false)

**Example Request:**
```
GET /api/services/?limit=10&offset=0&business_account_id=123&is_active=true
```

**Response:**
```json
{
  "services": [
    {
      "id": "uuid-string",
      "business_account_id": "uuid-string",
      "name": "Service Name",
      "description": "Service description",
      "duration_minutes": 60,
      "price": 50.00,
      "currency": "USD",
      "category": "Category",
      "is_active": true,
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "total": 1
}
```

**Status Codes:**
- `200 OK`: Services retrieved successfully
- `400 Bad Request`: Validation error
- `500 Internal Server Error`: Server error

### 6. Get Services by Business Account

**GET** `/api/services/business-account/{business_account_id}`

Retrieves all active services for a specific business account.

**Path Parameters:**
- `business_account_id`: Business account UUID

**Response:**
```json
[
  {
    "id": "uuid-string",
    "business_account_id": "uuid-string",
    "name": "Service Name",
    "description": "Service description",
    "duration_minutes": 60,
    "price": 50.00,
    "currency": "USD",
    "category": "Category",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

**Status Codes:**
- `200 OK`: Services retrieved successfully
- `400 Bad Request`: Missing business account ID
- `500 Internal Server Error`: Server error

## Data Models

### Service

```go
type Service struct {
    ID                 string     `json:"id"`
    BusinessAccountID  string     `json:"business_account_id"`
    Name               string     `json:"name"`
    Description        *string    `json:"description,omitempty"`
    DurationMinutes    int        `json:"duration_minutes"`
    Price              float64    `json:"price"`
    Currency           string     `json:"currency"`
    Category           *string    `json:"category,omitempty"`
    IsActive           bool       `json:"is_active"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
}
```

### Create Service Request

```go
type CreateServiceRequest struct {
    BusinessAccountID string  `json:"business_account_id"`
    Name              string  `json:"name"`
    Description       *string `json:"description,omitempty"`
    DurationMinutes   int     `json:"duration_minutes"`
    Price             float64 `json:"price"`
    Currency          string  `json:"currency"`
    Category          *string `json:"category,omitempty"`
}
```

### Update Service Request

```go
type UpdateServiceRequest struct {
    Name             *string  `json:"name,omitempty"`
    Description      *string  `json:"description,omitempty"`
    DurationMinutes  *int     `json:"duration_minutes,omitempty"`
    Price            *float64 `json:"price,omitempty"`
    Currency         *string  `json:"currency,omitempty"`
    Category         *string  `json:"category,omitempty"`
    IsActive         *bool    `json:"is_active,omitempty"`
}
```

## Validation Rules

### Create Service
- `business_account_id`: Required, must be a valid UUID
- `name`: Required, non-empty string
- `duration_minutes`: Required, must be greater than 0
- `price`: Required, must be non-negative
- `currency`: Optional, must be exactly 3 characters if provided

### Update Service
- All fields are optional
- If provided, fields must meet the same validation rules as create
- `name`: Cannot be empty if provided
- `duration_minutes`: Must be greater than 0 if provided
- `price`: Must be non-negative if provided
- `currency`: Must be exactly 3 characters if provided

### List Services
- `limit`: Must be between 1 and 100
- `offset`: Must be non-negative

## Error Responses

All error responses follow this format:

```json
{
  "message": "Error description",
  "type": "ERROR",
  "code": 1
}
```

### Common Error Codes
- `InvalidRequest`: Invalid request format
- `ValidationError`: Validation failed
- `NotFound`: Resource not found
- `InternalError`: Server error

## Examples

### Creating a Haircut Service

```bash
curl -X POST /api/services/ \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "business_account_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Men\'s Haircut",
    "description": "Professional men\'s haircut and styling",
    "duration_minutes": 30,
    "price": 25.00,
    "currency": "USD",
    "category": "Hair Services"
  }'
```

### Updating Service Price

```bash
curl -X PUT /api/services/550e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 30.00
  }'
```

### Listing Services for a Business

```bash
curl -X GET "/api/services/?business_account_id=550e8400-e29b-41d4-a716-446655440000&is_active=true" \
  -H "Authorization: Bearer <jwt-token>"
```

## Notes

- Services are automatically set to active when created
- Deleting a service is permanent and cannot be undone
- The API automatically sets timestamps for creation and updates
- Currency defaults to "USD" if not specified
- All timestamps are in ISO 8601 format with timezone information

