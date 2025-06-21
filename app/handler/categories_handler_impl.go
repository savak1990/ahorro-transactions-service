package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/savak1990/transactions-service/app/service"
)

// CategoriesHandlerImpl provides implementation for CategoriesHandler interface
type CategoriesHandlerImpl struct {
	service service.CategoriesService
}

func NewCategoriesHandlerImpl(svc service.CategoriesService) *CategoriesHandlerImpl {
	return &CategoriesHandlerImpl{service: svc}
}

func (h *CategoriesHandlerImpl) ListCategoriesForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		WriteJSONError(w, http.StatusBadRequest, "BadRequest", "Missing user_id in path")
		return
	}

	limit := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil {
			limit = n
		}
	}
	startKey := r.URL.Query().Get("start_key")

	input := models.ListCategoriesInput{
		UserID:   userID,
		Limit:    limit,
		StartKey: startKey,
	}

	categories, nextKey, err := h.service.ListCategoriesForUser(r.Context(), input)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "InternalServerError", err.Error())
		return
	}
	WriteJSONListResponse(w, categories, nextKey)
}

var _ CategoriesHandler = (*CategoriesHandlerImpl)(nil)
