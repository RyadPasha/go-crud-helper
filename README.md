
# Go HTTP CRUD Helper

## Overview

This project provides a **generic HTTP CRUD helper** for Go that allows you to create, read, update, and delete (CRUD) resources with minimal effort. The helper is designed to be flexible and reusable for any data model, making it suitable for building RESTful APIs in Go.

### Features:
- **Generic CRUD operations** using Go reflection
- **In-memory data store** to manage resources
- **Thread-safe** operations using `sync.Mutex`
- **Easy to integrate** into any project
- Supports **any data model** via reflection, so you only need to define your data models once and can reuse the helper functions for different models.

## Setup

### Requirements:
- Go 1.18+ (due to Go reflection)
- Basic understanding of Go's `net/http` package and reflection

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ryadpasha/go-http-crud-helper.git
   cd go-http-crud-helper
   ```

2. Build and run the server:

   ```bash
   go run main.go
   ```

   The server will start and listen on port 8080 by default.

## Usage

### Endpoints:
- **POST /item**: Create a new `Item`
- **GET /item?id=<id>**: Get an `Item` by ID
- **GET /item**: Get all `Items`
- **PUT /item?id=<id>**: Update an `Item` by ID
- **DELETE /item?id=<id>**: Delete an `Item` by ID

### Example:

**POST /item**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"title": "Learn Go", "done": false}' http://localhost:8080/item
```

**GET /item?id=1**

```bash
curl http://localhost:8080/item?id=1
```

## Extending

You can extend the functionality of this helper by:
- Adding validation for data input
- Integrating with a database instead of using an in-memory store
- Adding custom logging or authentication features

Feel free to contribute or create new branches with additional features!

## License

This project is licensed under the MIT License.
