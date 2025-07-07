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

	// Parse multiple balanceIds
	if balanceIds := ParseQueryStringArray(query, "balanceId"); len(balanceIds) > 0 {
		for _, balanceId := range balanceIds {
			if _, err := uuid.Parse(balanceId); err != nil {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid balanceId format: "+balanceId)
				return
			}
		}
		input.BalanceID = balanceIds
	}

	// Parse multiple categoryIds
	if categoryIds := ParseQueryStringArray(query, "categoryId"); len(categoryIds) > 0 {
		for _, categoryId := range categoryIds {
			if _, err := uuid.Parse(categoryId); err != nil {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid categoryId format: "+categoryId)
				return
			}
		}
		input.CategoryId = categoryIds
	}

	// Parse multiple categoryGroupIds
	if categoryGroupIds := ParseQueryStringArray(query, "categoryGroupId"); len(categoryGroupIds) > 0 {
		for _, categoryGroupId := range categoryGroupIds {
			if _, err := uuid.Parse(categoryGroupId); err != nil {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid categoryGroupId format: "+categoryGroupId)
				return
			}
		}
		input.CategoryGroupId = categoryGroupIds
	}

	// Parse multiple merchantIds
	if merchantIds := ParseQueryStringArray(query, "merchantId"); len(merchantIds) > 0 {
		for _, merchantId := range merchantIds {
			if _, err := uuid.Parse(merchantId); err != nil {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid merchantId format: "+merchantId)
				return
			}
		}
		input.MerchantId = merchantIds
	}

	// Parse multiple transactionIds
	if transactionIds := ParseQueryStringArray(query, "transactionId"); len(transactionIds) > 0 {
		for _, transactionId := range transactionIds {
			if _, err := uuid.Parse(transactionId); err != nil {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid transactionId format: "+transactionId)
				return
			}
		}
		input.TransactionId = transactionIds
	}

	// Parse single userId
	if userId := query.Get("userId"); userId != "" {
		if _, err := uuid.Parse(userId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid userId format")
			return
		}
		input.UserID = userId
	}

	// Parse single groupId
	if groupId := query.Get("groupId"); groupId != "" {
		if _, err := uuid.Parse(groupId); err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid groupId format")
			return
		}
		input.GroupID = groupId
	}

	// Parse multiple transaction types
	if transactionTypes := ParseQueryStringArray(query, "type"); len(transactionTypes) > 0 {
		for _, transactionType := range transactionTypes {
			if transactionType != "init" && transactionType != "expense" && transactionType != "income" && transactionType != "move_in" && transactionType != "move_out" {
				WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid type, must be 'init', 'expense', 'income', 'move_in', or 'move_out': "+transactionType)
				return
			}
		}
		input.Type = transactionTypes
	}

	// Parse startTime
	if startTime := query.Get("startTime"); startTime != "" {
		parsed, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid startTime format, must be RFC3339")
			return
		}
		input.StartTime = parsed
	}

	// Parse endTime
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
