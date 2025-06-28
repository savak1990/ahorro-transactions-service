package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/models"
	repo "github.com/savak1990/transactions-service/app/repo"
)

type ServiceImpl struct {
	repo repo.Repository
}

func NewServiceImpl(repo repo.Repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.ID = uuid.New()
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *ServiceImpl) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, transactionID)
}

func (s *ServiceImpl) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	return s.repo.ListTransactions(ctx, filter)
}

func (s *ServiceImpl) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateTransaction(ctx, tx)
}

func (s *ServiceImpl) DeleteTransaction(ctx context.Context, transactionID string) error {
	return s.repo.DeleteTransaction(ctx, transactionID)
}

func (s *ServiceImpl) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	balance.ID = uuid.New()
	return s.repo.CreateBalance(ctx, balance)
}

func (s *ServiceImpl) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	return s.repo.GetBalance(ctx, balanceID)
}

func (s *ServiceImpl) ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error) {
	return s.repo.ListBalances(ctx, filter)
}

func (s *ServiceImpl) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	balance.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateBalance(ctx, balance)
}

func (s *ServiceImpl) DeleteBalance(ctx context.Context, balanceID string) error {
	return s.repo.DeleteBalance(ctx, balanceID)
}

func (s *ServiceImpl) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	category.ID = uuid.New()
	return s.repo.CreateCategory(ctx, category)
}

func (s *ServiceImpl) ListCategories(ctx context.Context, filter models.ListCategoriesInput) ([]models.Category, error) {
	categories, _, err := s.repo.ListCategories(ctx, filter)
	return categories, err
}

func (s *ServiceImpl) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	return s.repo.GetCategory(ctx, categoryID)
}

func (s *ServiceImpl) UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	category.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateCategory(ctx, category)
}

func (s *ServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	return s.repo.DeleteCategory(ctx, categoryID)
}

func (s *ServiceImpl) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	merchant.ID = uuid.New()
	return s.repo.CreateMerchant(ctx, merchant)
}

func (s *ServiceImpl) GetMerchant(ctx context.Context, merchantID string) (*models.Merchant, error) {
	return s.repo.GetMerchant(ctx, merchantID)
}

func (s *ServiceImpl) ListMerchants(ctx context.Context, filter models.ListMerchantsFilter) ([]models.Merchant, string, error) {
	return s.repo.ListMerchants(ctx, filter)
}

func (s *ServiceImpl) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	merchant.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateMerchant(ctx, merchant)
}

func (s *ServiceImpl) DeleteMerchant(ctx context.Context, merchantID string) error {
	return s.repo.DeleteMerchant(ctx, merchantID)
}

func (s *ServiceImpl) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsFilter) ([]models.TransactionEntry, string, error) {
	return s.repo.ListTransactionEntries(ctx, filter)
}

// Ensure ServiceImpl implements Service
var _ Service = (*ServiceImpl)(nil)
