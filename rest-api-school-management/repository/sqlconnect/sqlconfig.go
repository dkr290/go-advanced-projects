// Package sqlconnect - for connecting and testing connection to the mariadb database server
package sqlconnect

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(conf *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBDatabase,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	var count int
	for {

		err := db.Ping()
		if err != nil {
			fmt.Println("Error connection to mariadb", err)
			if count > 10 {
				return nil, fmt.Errorf("error connection for more then 10 times %v", err)
			}
			count++
			time.Sleep(2 * time.Second)
			continue
		} else {
			fmt.Println("Connected to mariadb")
			break
		}

	}
	return db, nil
}
