// File: crud_helper.go
// Author: Mohamed Riyad
// Email: mohamed.riyad@example.com
// Date: November 2024
// License: MIT
// Description: This file contains a generic HTTP CRUD helper for Go that can be reused across any project.
// It provides basic Create, Read, Update, and Delete (CRUD) functionality for any type of data model
// using reflection. The helper uses an in-memory store (map) to manage items, and the functions
// are thread-safe using a mutex to ensure concurrent access.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

// Store is a generic structure to hold and manage items in memory.
type Store struct {
	data    map[int]interface{}
	nextID  int
	itemMux sync.Mutex
}

// NewStore creates a new instance of Store.
func NewStore() *Store {
	return &Store{
		data:   make(map[int]interface{}),
		nextID: 1,
	}
}

// Create adds a new item to the store and returns the item with an assigned ID.
func (s *Store) Create(item interface{}) interface{} {
	s.itemMux.Lock()
	defer s.itemMux.Unlock()

	// Assign a new ID and store the item
	id := s.nextID
	s.nextID++
	itemValue := reflect.ValueOf(item).Elem()
	itemValue.FieldByName("ID").SetInt(int64(id))

	s.data[id] = item
	return item
}

// Get retrieves an item by its ID.
func (s *Store) Get(id int, result interface{}) bool {
	s.itemMux.Lock()
	defer s.itemMux.Unlock()

	item, exists := s.data[id]
	if !exists {
		return false
	}

	// Populate result struct with the found item
	itemValue := reflect.ValueOf(item)
	reflect.ValueOf(result).Elem().Set(itemValue)
	return true
}

// GetAll retrieves all items in the store.
func (s *Store) GetAll(result interface{}) {
	s.itemMux.Lock()
	defer s.itemMux.Unlock()

	// Populate result slice with all items
	itemSlice := reflect.ValueOf(result).Elem()
	for _, item := range s.data {
		itemSlice.Set(reflect.Append(itemSlice, reflect.ValueOf(item)))
	}
}

// Update updates an existing item in the store.
func (s *Store) Update(id int, updatedItem interface{}) bool {
	s.itemMux.Lock()
	defer s.itemMux.Unlock()

	_, exists := s.data[id]
	if !exists {
		return false
	}

	// Update the item
	s.data[id] = updatedItem
	return true
}

// Delete removes an item by its ID.
func (s *Store) Delete(id int) bool {
	s.itemMux.Lock()
	defer s.itemMux.Unlock()

	_, exists := s.data[id]
	if !exists {
		return false
	}

	delete(s.data, id)
	return true
}

// handleRequest handles HTTP requests for CRUD operations on any data model.
func handleRequest(store *Store, modelType reflect.Type, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Create item
		newItem := reflect.New(modelType).Interface()
		if err := json.NewDecoder(r.Body).Decode(newItem); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		createdItem := store.Create(newItem)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdItem)

	case http.MethodGet:
		// Get all items
		if r.URL.Query().Get("id") == "" {
			result := reflect.New(reflect.SliceOf(modelType)).Interface()
			store.GetAll(result)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			return
		}

		// Get item by ID
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		result := reflect.New(modelType).Interface()
		if store.Get(id, result) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		} else {
			http.Error(w, "Item not found", http.StatusNotFound)
		}

	case http.MethodPut:
		// Update item by ID
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		updatedItem := reflect.New(modelType).Interface()
		if err := json.NewDecoder(r.Body).Decode(updatedItem); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if store.Update(id, updatedItem) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedItem)
		} else {
			http.Error(w, "Item not found", http.StatusNotFound)
		}

	case http.MethodDelete:
		// Delete item by ID
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if store.Delete(id) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(w, "Item not found", http.StatusNotFound)
		}

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}
