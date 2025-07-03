package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/helpers"
	"github.com/savak1990/transactions-service/app/models"
)

// Category handlers
func (h *HandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var createCategoryDto models.CreateCategoryDto
	if err := json.NewDecoder(r.Body).Decode(&createCategoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model using the new converter
	category, err := models.FromAPICreateCategory(createCategoryDto)
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
		UserID:  r.URL.Query().Get("userId"),
		GroupBy: r.URL.Query().Get("groupBy"),
		SortBy:  r.URL.Query().Get("sortBy"),
		Order:   r.URL.Query().Get("order"),
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

	// Check if groupBy parameter is set to "categoryGroup"
	if filter.GroupBy == "categoryGroup" {
		// Group categories by category group
		groupedCategories := h.groupCategoriesByGroup(results)
		WriteJSONListResponse(w, groupedCategories, "")
	} else {
		// Convert to DTOs for regular response
		categoryDtos := make([]models.CategoryDto, len(results))
		for i, category := range results {
			categoryDtos[i] = models.ToAPICategory(&category)
		}
		// Use paginated response structure (without actual pagination for now)
		WriteJSONListResponse(w, categoryDtos, "")
	}
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
	var updateCategoryDto models.UpdateCategoryDto
	if err := json.NewDecoder(r.Body).Decode(&updateCategoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate category ID format
	_, err := uuid.Parse(categoryID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category ID format")
		return
	}

	// Convert DTO to DAO model using the new converter
	category, err := models.FromAPIUpdateCategory(updateCategoryDto, categoryID)
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

// groupCategoriesByGroup groups categories by their category group
func (h *HandlerImpl) groupCategoriesByGroup(categories []models.Category) []models.CategoryGroupWithCategoriesDto {
	// Map to group categories by category group ID
	groupMap := make(map[string]*models.CategoryGroupWithCategoriesDto)

	for _, category := range categories {
		categoryDto := models.ToAPICategory(&category)

		// Get category group ID, handle cases where it might be empty
		groupID := categoryDto.CategoryGroupID
		if groupID == "" {
			groupID = "ungrouped" // Handle categories without a group
		}

		// Check if we already have this group
		if group, exists := groupMap[groupID]; exists {
			// Add category to existing group
			group.Categories = append(group.Categories, categoryDto)
		} else {
			// Create new group
			newGroup := &models.CategoryGroupWithCategoriesDto{
				CategoryGroupID:       categoryDto.CategoryGroupID,
				CategoryGroupName:     categoryDto.CategoryGroupName,
				CategoryGroupImageUrl: categoryDto.CategoryGroupImageUrl,
				CategoryGroupRank:     categoryDto.CategoryGroupRank,
				Categories:            []models.CategoryDto{categoryDto},
			}

			// Handle ungrouped categories
			if groupID == "ungrouped" {
				newGroup.CategoryGroupID = ""
				newGroup.CategoryGroupName = "Ungrouped"
				newGroup.CategoryGroupImageUrl = nil
				newGroup.CategoryGroupRank = nil
			}

			groupMap[groupID] = newGroup
		}
	}
	// Convert map to slice
	result := make([]models.CategoryGroupWithCategoriesDto, 0, len(groupMap))
	for _, group := range groupMap {
		result = append(result, *group)
	}

	// Sort groups by rank (highest first, then by name)
	sort.Slice(result, func(i, j int) bool {
		// Handle nil ranks
		rankI := result[i].CategoryGroupRank
		rankJ := result[j].CategoryGroupRank

		// If both have ranks, sort by rank (descending)
		if rankI != nil && rankJ != nil {
			if *rankI != *rankJ {
				return *rankI > *rankJ // Higher rank first
			}
		} else if rankI != nil {
			return true // Items with ranks come before items without ranks
		} else if rankJ != nil {
			return false // Items with ranks come before items without ranks
		}

		// If ranks are equal or both nil, sort by name (ascending)
		return result[i].CategoryGroupName < result[j].CategoryGroupName
	})

	return result
}
