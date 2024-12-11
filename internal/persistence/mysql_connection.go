package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func CreateConnection() *sql.DB {
	var db *sql.DB

	cfg := mysql.Config{
		User:   viper.GetString("mysql.user"),
		Passwd: viper.GetString("mysql.password"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%v:%v", viper.Get("mysql.host"), viper.Get("mysql.port")),
		DBName: viper.GetString("mysql.database"),
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}
