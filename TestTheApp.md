# Testing Simpler Project API

## Overview
This document provides instructions on testing the REST API endpoints using Postman (or similar tools) both locally and on the live server at: http://165.232.66.26:8080/products


## Testing Endpoints

### 1. **Create a Product**

- **Method**: `POST`
- **URL**:
  - **Local**: `http://localhost:8080/products`
  - **Deployed**: `http://165.232.66.26:8080/products`
- **Headers**: `Content-Type: application/json`
- **Body**:
    ```json
    {
      "name": "Product 1",
      "price": 29.99,
      "quantity": 10,
      "description": "A sample product description"
    }
    ```
- **Expected Response**: `201 Created`

---

### 2. **Retrieve All Products**

- **Method**: `GET`
- **URL**:
  - **Local**: `http://localhost:8080/products`
  - **Deployed**: `http://165.232.66.26:8080/products`
- **Expected Response**: A list of all products (`200 OK`).

---

### 3. **Pagination**

#### Fetch Products with Pagination

- **Method**: `GET`
- **URL**:
  - **Local**: `http://localhost:8080/products?limit=5&page=2`
  - **Deployed**: `http://165.232.66.26:8080/products?limit=5&page=2`
- **Expected Response**:
    ```json
    {
      "data": [ ... ],
      "total": 15, // Example total number of products
      "page": 2,
      "limit": 5
    }
    ```

#### Edge Cases
- **Limit=0**: `http://localhost:8080/products?limit=0`
  - Expected default limit (`10`).
- **Non-existent page**: `http://localhost:8080/products?page=1000`
  - Should return an empty list of products.

---

### 4. **Retrieve a Product by ID**

- **Method**: `GET`
- **URL**:
  - **Local**: `http://localhost:8080/products/{id}`
  - **Deployed**: `http://165.232.66.26:8080/products/{id}`
  - Replace `{id}` with the actual product ID.
- **Expected Response**: `200 OK` if the product exists, `404 Not Found` if not.

---

### 5. **Update a Product**

- **Method**: `PUT`
- **URL**:
  - **Local**: `http://localhost:8080/products/{id}`
  - **Deployed**: `http://165.232.66.26:8080/products/{id}`
- **Body**:
    ```json
    {
      "name": "Updated Product Name",
      "price": 49.99,
      "quantity": 20,
      "description": "Updated description"
    }
    ```
- **Expected Response**: `200 OK` if the product is updated, `404 Not Found` if the product does not exist.

#### Edge Cases
- **Invalid ID**: `http://localhost:8080/products/invalid-id`
  - Should return `400 Bad Request`.

---

### 6. **Delete a Product**

- **Method**: `DELETE`
- **URL**:
  - **Local**: `http://localhost:8080/products/{id}`
  - **Deployed**: `http://165.232.66.26:8080/products/{id}`
- **Expected Response**: `200 OK` if the product is deleted, `404 Not Found` if the product does not exist.

#### Edge Cases
- **Invalid ID**: `http://localhost:8080/products/invalid-id`
  - Should return `400 Bad Request`.

---

## Error Scenarios

- **Invalid JSON Body**: Missing or invalid fields during `POST` or `PUT` requests should return `400 Bad Request`.
- **Non-existent Product ID**: Accessing, updating, or deleting a non-existent product ID should return `404 Not Found`.

