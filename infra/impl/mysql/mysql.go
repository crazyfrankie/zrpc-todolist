package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func New() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("mysql open, dsn: %s, err: %w", dsn, err)
	}

	return db, nil
}
