// File: main.go
// Author: Mohamed Riyad
// Email: mohamed.riyad@example.com
// Date: November 2024
// License: MIT
// Description: This file demonstrates how to use the generic HTTP CRUD helper defined in crud_helper.go.
// It registers a basic data model ("Item") and exposes a set of RESTful CRUD operations for it via HTTP.
// It also shows how to start a simple HTTP server using the standard Go net/http package.

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Item represents a generic data model for demonstration purposes.
type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	// Create a new instance of the generic Store
	store := NewStore()

	// Register CRUD operations for the "Item" data model
	http.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(store, reflect.TypeOf(Item{}), w, r)
	})

	// Start the HTTP server on port 8080
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
