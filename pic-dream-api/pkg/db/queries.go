package db

import (
	"context"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
	"github.com/google/uuid"
)

func (s *SupabasePostgresql) CreateAccount(account *types.Account) error {
	_, err := s.Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}

func (s *SupabasePostgresql) GetAccountByUserID(userID uuid.UUID) (types.Account, error) {
	var account types.Account
	err := s.Bun.NewSelect().Model(&account).Where("user_id = ?", userID).Scan(context.Background())

	return account, err
}

func (s *SupabasePostgresql) UpdateAccount(account *types.Account) error {
	_, err := s.Bun.NewUpdate().Model(account).WherePK().Exec(context.Background())
	return err
}
