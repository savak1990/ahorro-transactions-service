package handler

import (
	"github.com/savak1990/transactions-service/app/service"
)

// HandlerImpl implements Handler and wires to TransactionsService
// Add a field for the service dependency
type HandlerImpl struct {
	Service service.Service
}

func NewHandlerImpl(svc service.Service) *HandlerImpl {
	return &HandlerImpl{Service: svc}
}

var _ Handler = (*HandlerImpl)(nil)
