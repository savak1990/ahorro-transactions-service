package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/helpers"
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

	// Parse single transaction type
	if transactionType := query.Get("type"); transactionType != "" {
		if !_isSupportedType(transactionType) {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid type format")
			return
		}
		input.Type = transactionType
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

	// Get grouping query parameter
	if grouping := query.Get("grouping"); grouping != "" {
		switch grouping {
		case "categoryGroup":
			input.Grouping = models.GroupingCategoryGroup
		case "category":
			input.Grouping = models.GroupingCategory
		case "merchant":
			input.Grouping = models.GroupingMerchant
		case "balance":
			input.Grouping = models.GroupingBalance
		case "currency":
			input.Grouping = models.GroupingCurrency
		case "quarter":
			input.Grouping = models.GroupingQuarter
		case "year":
			input.Grouping = models.GroupingYear
		case "month":
			input.Grouping = models.GroupingMonth
		case "week":
			input.Grouping = models.GroupingWeek
		case "day":
			input.Grouping = models.GroupingDay
		default:
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid grouping parameter")
			return
		}
	} else {
		// Default grouping is by category
		input.Grouping = models.GroupingCategory
	}

	// Get limit with default 10
	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if n, err := helpers.ParseInt(limitStr); err == nil {
			limit = n
		} else {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid limit format: "+err.Error())
			return
		}
	}
	input.Limit = limit

	// Parse sort parameter
	input.Sort = "amount"
	if sort := query.Get("sort"); sort != "" {
		if sort == "amount" || sort == "count" || sort == "label" {
			input.Sort = sort
		} else {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid sort format, must be 'amount', 'count' or 'label'")
			return
		}
	}

	// Parse order parameter
	input.Order = "desc"
	if order := query.Get("order"); order != "" {
		if order == "asc" || order == "desc" {
			input.Order = order
		} else {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid order format, must be 'asc' or 'desc'")
			return
		}
	}

	// Parse displayCurrency, nil if not provided
	if displayCurrency := query.Get("displayCurrency"); displayCurrency != "" {
		if !helpers.IsValidCurrencyCode(displayCurrency) {
			WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid displayCurrency format")
			return
		}
		input.DisplayCurrency = strings.ToUpper(displayCurrency)
	} else {
		input.DisplayCurrency = ""
	}

	// Call service
	statsList, err := h.Service.GetTransactionStats(r.Context(), input)
	if err != nil {
		h.handleServiceError(w, err, "GetTransactionStats")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	WriteJSONListResponse(w, statsList, "")
}

func _isSupportedType(transactionType string) bool {
	switch transactionType {
	case "init", "expense", "income", "move_in", "move_out":
		return true
	default:
		return false
	}
}
