# Simpler Project - Product Management REST API

## Overview
This is a RESTful API for managing products in a store. It supports CRUD operations and is built using Go, Gin for routing, and PostgreSQL as the database. The application allows creating, retrieving, updating, and deleting products, along with pagination and validation.

The live server is deployed on Digital Ocean and accessible at: http://165.232.66.26:8080/products


---

## Features
- Create, retrieve, update, and delete products.
- Validation of input fields (name, price, quantity).
- Pagination for fetching products.
- Error handling for non-existent resources and invalid inputs.

---

## Unit Tests

The application contains unit tests for the following scenarios:
- **Product Creation**: Validates the creation of a new product and tests error scenarios for missing or invalid fields.
- **Product Retrieval**: Tests fetching all products with and without pagination, including edge cases like non-existent pages.
- **Product Retrieval by ID**: Retrieves a product by ID and handles cases where the product does not exist.
- **Product Update**: Validates the updating of a product, including error handling for incomplete or invalid data and non-existent IDs.
- **Product Deletion**: Tests product deletion by ID and checks responses for invalid or non-existent product IDs.

---

## Setup and Deployment with Docker

### Prerequisites
Ensure you have Docker and Docker Compose installed.

### Steps to Deploy Locally
1. **Clone the Repository**
    ```bash
    git clone <repository-url>
    cd simpler-project
    ```

2. **Create Environment File**
   Create a `.env` file in the root directory with the following content:
    ```env
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=oOYyyha5lFkEyiWsy855
    POSTGRES_DB=products_db
    DB_HOST=db
    DB_PORT=5432
    GO_ENV=production
    ```

3. **Build and Run the Containers**
    ```bash
    docker-compose up --build
    ```
   This will start both the database and the application containers. The app will be available at: http://localhost:8080/products



