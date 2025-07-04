package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/models"
)

// GET /transactions/stats
func (h *HandlerImpl) GetTransactionStats(w http.ResponseWriter, r *http.Request) {
	var input models.TransactionStatsInput

	// Parse query parameters
	query := r.URL.Query()

	if balanceId := query.Get("balanceId"); balanceId != "" {
		if _, err := uuid.Parse(balanceId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid balanceId format")
			return
		}
		input.BalanceID = &balanceId
	}

	if categoryId := query.Get("categoryId"); categoryId != "" {
		if _, err := uuid.Parse(categoryId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid categoryId format")
			return
		}
		input.CategoryId = &categoryId
	}

	if categoryGroupId := query.Get("categoryGroupId"); categoryGroupId != "" {
		input.CategoryGroupId = &categoryGroupId
	}

	if merchantId := query.Get("merchantId"); merchantId != "" {
		if _, err := uuid.Parse(merchantId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid merchantId format")
			return
		}
		input.MerchantId = &merchantId
	}

	if userId := query.Get("userId"); userId != "" {
		if _, err := uuid.Parse(userId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid userId format")
			return
		}
		input.UserID = &userId
	}

	if groupId := query.Get("groupId"); groupId != "" {
		if _, err := uuid.Parse(groupId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid groupId format")
			return
		}
		input.GroupID = &groupId
	}

	if transactionType := query.Get("type"); transactionType != "" {
		if transactionType != "expense" && transactionType != "income" {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid type, must be 'expense' or 'income'")
			return
		}
		input.Type = &transactionType
	}

	if startTime := query.Get("startTime"); startTime != "" {
		parsed, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid startTime format, must be RFC3339")
			return
		}
		input.StartTime = &parsed
	}

	if endTime := query.Get("endTime"); endTime != "" {
		parsed, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid endTime format, must be RFC3339")
			return
		}
		input.EndTime = &parsed
	}

	// Call service
	stats, err := h.Service.GetTransactionStats(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err, "GetTransactionStats")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
