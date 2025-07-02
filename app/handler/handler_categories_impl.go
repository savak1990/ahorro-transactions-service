package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/helpers"
	"github.com/savak1990/transactions-service/app/models"
)

// Category handlers
func (h *HandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryDto models.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&categoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category data: "+err.Error())
		return
	}

	created, err := h.Service.CreateCategory(r.Context(), *category)
	if err != nil {
		h.handleServiceError(w, err, "CreateCategory")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategory(created)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) ListCategories(w http.ResponseWriter, r *http.Request) {
	filter := models.ListCategoriesInput{
		UserID: r.URL.Query().Get("userId"),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if n, err := helpers.ParseInt(limit); err == nil {
			filter.Limit = n
		}
	}
	results, err := h.Service.ListCategories(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListCategories")
		return
	}

	// Convert to DTOs for response
	categoryDtos := make([]models.CategoryDto, len(results))
	for i, category := range results {
		categoryDtos[i] = models.ToAPICategory(&category)
	}

	// Use paginated response structure (without actual pagination for now)
	WriteJSONListResponse(w, categoryDtos, "")
}

func (h *HandlerImpl) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]
	if categoryID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_id")
		return
	}

	category, err := h.Service.GetCategory(r.Context(), categoryID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "category", categoryID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetCategory")
		return
	}

	if category == nil {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, "Category not found")
		return
	}

	// Convert to DTO for response
	responseDto := models.ToAPICategory(category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]
	var categoryDto models.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&categoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse category ID and set it in DTO
	id, err := uuid.Parse(categoryID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category ID format")
		return
	}
	categoryDto.CategoryID = id.String()

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category data: "+err.Error())
		return
	}

	updated, err := h.Service.UpdateCategory(r.Context(), *category)
	if err != nil {
		h.handleServiceError(w, err, "UpdateCategory")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategory(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]

	if categoryID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_id parameter")
		return
	}

	err := h.Service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteCategory")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) DeleteCategoriesByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing userId parameter")
		return
	}

	err := h.Service.DeleteCategoriesByUserId(r.Context(), userId)
	if err != nil {
		h.handleServiceError(w, err, "DeleteCategoriesByUserId")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
