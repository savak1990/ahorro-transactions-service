package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/helpers"
	"github.com/savak1990/transactions-service/app/models"
)

// CategoryGroup handlers
func (h *HandlerImpl) CreateCategoryGroup(w http.ResponseWriter, r *http.Request) {
	var categoryGroupDto models.CategoryGroupDto
	if err := json.NewDecoder(r.Body).Decode(&categoryGroupDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model
	categoryGroup, err := models.FromAPICategoryGroup(categoryGroupDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category group data: "+err.Error())
		return
	}

	created, err := h.Service.CreateCategoryGroup(r.Context(), *categoryGroup)
	if err != nil {
		h.handleServiceError(w, err, "CreateCategoryGroup")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategoryGroup(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) ListCategoryGroups(w http.ResponseWriter, r *http.Request) {
	filter := models.ListCategoryGroupsInput{
		SortBy: r.URL.Query().Get("sortBy"),
		Order:  r.URL.Query().Get("order"),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if n, err := helpers.ParseInt(limit); err == nil {
			filter.Limit = n
		}
	}

	results, err := h.Service.ListCategoryGroups(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListCategoryGroups")
		return
	}

	// Convert to DTOs for response
	categoryGroupDtos := make([]models.CategoryGroupDto, len(results))
	for i, categoryGroup := range results {
		categoryGroupDtos[i] = models.ToAPICategoryGroup(&categoryGroup)
	}

	// Use paginated response structure (without actual pagination for now)
	WriteJSONListResponse(w, categoryGroupDtos, "")
}

func (h *HandlerImpl) GetCategoryGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryGroupID := vars["category_group_id"]
	if categoryGroupID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_group_id")
		return
	}

	categoryGroup, err := h.Service.GetCategoryGroup(r.Context(), categoryGroupID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "category group", categoryGroupID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetCategoryGroup")
		return
	}

	if categoryGroup == nil {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, "Category group not found")
		return
	}

	// Convert to DTO for response
	responseDto := models.ToAPICategoryGroup(categoryGroup)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) UpdateCategoryGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryGroupID := vars["category_group_id"]
	var categoryGroupDto models.CategoryGroupDto
	if err := json.NewDecoder(r.Body).Decode(&categoryGroupDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse category group ID and set it in DTO
	id, err := uuid.Parse(categoryGroupID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category group ID format")
		return
	}
	categoryGroupDto.CategoryGroupId = id.String()

	// Convert DTO to DAO model
	categoryGroup, err := models.FromAPICategoryGroup(categoryGroupDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category group data: "+err.Error())
		return
	}

	updated, err := h.Service.UpdateCategoryGroup(r.Context(), *categoryGroup)
	if err != nil {
		h.handleServiceError(w, err, "UpdateCategoryGroup")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategoryGroup(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) DeleteCategoryGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryGroupID := vars["category_group_id"]
	if categoryGroupID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_group_id")
		return
	}
	err := h.Service.DeleteCategoryGroup(r.Context(), categoryGroupID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteCategoryGroup")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
