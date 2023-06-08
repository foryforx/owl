package conf

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbPool *sqlx.DB
var onceDB sync.Once

// GetConnPool instantiates pgx.ConnPool
func GetConnPool(config *Configuration) (*sqlx.DB, error) {
	var err error

	onceDB.Do(func() {
		dbPool, err = sqlx.Open("postgres",
			fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPwd, conf.DBDatabase, conf.DBSSLMode))
		if err == nil {
			dbPool.SetMaxOpenConns(conf.DBMaxConnection)
		}
	})
	return dbPool, err
}
