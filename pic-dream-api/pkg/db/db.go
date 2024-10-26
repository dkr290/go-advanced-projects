package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dkr290/go-advanced-projects/pic-dream-api/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type PictureDatabase interface {
	CreateAccount(account *types.Account) error
	GetAccountByUserID(userID uuid.UUID) (types.Account, error)
	UpdateAccount(account *types.Account) error
}

type SupabasePostgresql struct {
	Bun *bun.DB
}

func CreateDatabase(
	dbname string,
	dbuser string,
	dbpassword string,
	dbhost string,
) (*sql.DB, error) {
	hostArr := strings.Split(dbhost, ":")
	host := hostArr[0]
	port := "5432"
	if len(hostArr) > 1 {
		port = hostArr[1]
	}
	uri := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbuser,
		dbpassword,
		dbname,
		host,
		port,
	)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Init() (*bun.DB, error) {
	var Bun *bun.DB
	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		log.Fatal("DB_HOST is mandatory")
	}
	user := os.Getenv("DB_USER")
	if len(user) == 0 {
		log.Fatal("DB_USER is mandatory")
	}
	pass := os.Getenv("DB_PASSWORD")
	if len(pass) == 0 {
		log.Fatal("DB_PASSWORD is mandatory")
	}
	dbname := os.Getenv("DB_NAME")
	if len(dbname) == 0 {
		log.Fatal("DB_NAME is mandatory")
	}
	db, err := CreateDatabase(dbname, user, pass, host)
	if err != nil {
		return &bun.DB{}, err
	}
	if err := db.Ping(); err != nil {
		return &bun.DB{}, err
	}
	Bun = bun.NewDB(db, pgdialect.New())
	if len(os.Getenv("APP_DEBUG")) > 0 {
		Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return Bun, nil
}
