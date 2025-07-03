package service

import (
	"context"
	"fmt"
	"time"

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
	tx.ID = models.NewTransactionID()
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *ServiceImpl) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, transactionID)
}

func (s *ServiceImpl) ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error) {
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
	balance.ID = models.NewBalanceID()
	return s.repo.CreateBalance(ctx, balance)
}

func (s *ServiceImpl) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	return s.repo.GetBalance(ctx, balanceID)
}

func (s *ServiceImpl) ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error) {
	return s.repo.ListBalances(ctx, filter)
}

func (s *ServiceImpl) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	balance.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateBalance(ctx, balance)
}

func (s *ServiceImpl) DeleteBalance(ctx context.Context, balanceID string) error {
	return s.repo.DeleteBalance(ctx, balanceID)
}

func (s *ServiceImpl) DeleteBalancesByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteBalancesByUserId(ctx, userId)
}

func (s *ServiceImpl) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	category.ID = models.NewCategoryID()
	return s.repo.CreateCategory(ctx, category)
}

func (s *ServiceImpl) ListCategories(ctx context.Context, filter models.ListCategoriesInput) ([]models.Category, error) {
	return s.repo.ListCategories(ctx, filter)
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

func (s *ServiceImpl) DeleteCategoriesByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteCategoriesByUserId(ctx, userId)
}

func (s *ServiceImpl) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	// Check if a merchant with the same name already exists for this user
	existingMerchant, err := s.repo.GetMerchantByNameAndUserId(ctx, merchant.Name, merchant.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing merchant: %w", err)
	}
	if existingMerchant != nil {
		return nil, fmt.Errorf("merchant with name '%s' already exists for this user", merchant.Name)
	}

	merchant.ID = models.NewMerchantID()
	return s.repo.CreateMerchant(ctx, merchant)
}

func (s *ServiceImpl) GetMerchant(ctx context.Context, merchantID string) (*models.Merchant, error) {
	return s.repo.GetMerchant(ctx, merchantID)
}

func (s *ServiceImpl) ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error) {
	return s.repo.ListMerchants(ctx, filter)
}

func (s *ServiceImpl) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	merchant.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateMerchant(ctx, merchant)
}

func (s *ServiceImpl) DeleteMerchant(ctx context.Context, merchantID string) error {
	return s.repo.DeleteMerchant(ctx, merchantID)
}

func (s *ServiceImpl) DeleteMerchantsByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteMerchantsByUserId(ctx, userId)
}

func (s *ServiceImpl) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error) {
	return s.repo.ListTransactionEntries(ctx, filter)
}

// CategoryGroup service methods

func (s *ServiceImpl) CreateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	categoryGroup.ID = models.NewCategoryGroupID()
	return s.repo.CreateCategoryGroup(ctx, categoryGroup)
}

func (s *ServiceImpl) ListCategoryGroups(ctx context.Context, filter models.ListCategoryGroupsInput) ([]models.CategoryGroup, error) {
	return s.repo.ListCategoryGroups(ctx, filter)
}

func (s *ServiceImpl) GetCategoryGroup(ctx context.Context, categoryGroupID string) (*models.CategoryGroup, error) {
	return s.repo.GetCategoryGroup(ctx, categoryGroupID)
}

func (s *ServiceImpl) UpdateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	categoryGroup.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateCategoryGroup(ctx, categoryGroup)
}

func (s *ServiceImpl) DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error {
	return s.repo.DeleteCategoryGroup(ctx, categoryGroupID)
}

// GetTransactionStats retrieves aggregated transaction statistics
func (s *ServiceImpl) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) (*models.TransactionStatsResponseDto, error) {
	// Get raw stats from repository
	rawStats, err := s.repo.GetTransactionStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := &models.TransactionStatsResponseDto{
		Totals: rawStats,
	}

	return response, nil
}

// Ensure ServiceImpl implements Service
var _ Service = (*ServiceImpl)(nil)
