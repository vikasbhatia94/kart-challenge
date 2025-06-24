# 🍔 Food Ordering API - Advanced Challenge (Go)

This is a complete, production-style Go implementation of the [Kart Advanced Backend Challenge](https://github.com/oolio-group/kart-challenge/blob/advanced-challenge/backend-challenge/README.md), based on OpenAPI 3.1.

---

## Features

- Full OpenAPI 3.1 server generated via `oapi-codegen`.
- REST endpoints for:
  - Listing all products.
  - Viewing an individual product by its ID.
  - Placing an order with an optional promo code.
- A robust promo code validation engine that:
  - Scans and processes raw `.gz` files for 8–10 character alphanumeric codes.
  - Validates a code only if it appears in at least two of the three source files.
  - Caches the final validated promo code list in a `.gob` file for extremely fast subsequent startups.
- A `--refresh-promos` command-line flag to force a re-index of the promo code files.
- A complete unit test suite for the server implementation.

---

## API Endpoints

The server exposes the following endpoints, prefixed with `/api`. For example: `http://localhost:8080/api/product`.

| Method | Path           | Description               |
|--------|----------------|---------------------------|
| GET    | `/api/product`     | List all products         |
| GET    | `/api/product/{id}`| Get product by ID         |
| POST   | `/api/order`       | Place an order with promo |

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Docker (optional, for containerized deployment)

### 1. Clone the repository
```bash
git clone https://github.com/YOUR_USERNAME/kart-challenge
cd kart-challenge
```

### 2. Set up Promo Code Files
Create a directory named `promos` in the project root and place the three `.gz` files inside it. The application will process these on first run.
```
promos/
├── couponbase1.gz
├── couponbase2.gz
└── couponbase3.gz
```

### 3. Build and Run
The application will automatically download Go modules and build the server.

```bash
# To run with cached promo codes (fastest startup):
go run ./main.go

# To force a re-scan of the promo files:
go run ./main.go --refresh-promos
```
The server will be running on `http://localhost:8080`.

---

## 🧪 Running Tests
To run the full suite of unit tests, execute the following command from the project root:
```bash
go test -v ./...
```

---

## Design Notes & Assumptions

This project is built to be a robust and well-structured service, while making a few simplifying assumptions appropriate for the scope of the challenge.

- **Promo Code Logic**: The core validation requirement (a code appearing in ≥2 of 3 files) is implemented.
- **Discount Calculation**: A flat 10% discount is applied for any valid promo code. This is defined as a constant in `impl/server.go`. A production system would likely involve a more dynamic rules engine (e.g., specific discounts per code, usage limits, expiration dates) managed in a database.
- **Persistence**: The product catalog is stored in-memory. The promo code list is persisted to a binary `.gob` file for performance. A production service would use a database for both.
- **Authentication**: The generated code includes a basic authentication middleware hook. For this challenge, it performs a simple check for a `Bearer` token but does not validate it against a real user database.

---

## 🐳 Docker Support
The project includes a `Dockerfile` for easy containerization.

**Build the Docker image:**
```bash
docker build -t kart .
```

**Run the Docker container:**
This command mounts the local `promos` directory into the container.
```bash
docker run -p 8080:8080 -v $PWD/promos:/app/promos kart
```

---

## 📎 API Reference
The API conforms to the OpenAPI specification available here: [https://orderfoodonline.deno.dev/public/openapi.html](https://orderfoodonline.deno.dev/public/openapi.html)
