package configuration

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	DB *sql.DB
}

func NewMySQL(url string, connMaxLifetime time.Duration, maxOpenConns, maxIdleConns int) MySQL {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	return MySQL{DB: db}
}

func (impl *MySQL) BuildUpdateData(data map[string]interface{}) ([]string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	for field, value := range data {
		if v, ok := value.(string); ok && v == "" {
			continue
		}

		fields = append(fields, field+" = ?")
		values = append(values, value)
	}

	return fields, values
}
