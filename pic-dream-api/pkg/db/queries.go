package db

import (
	"context"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
)

type Database interface{}

func CreateAccount(account types.Account) error {
	_, err := Bun.NewInsert().Model(&account).Exec(context.Background())
	return err
}
