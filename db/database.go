package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// CREATE TABLE `results` (
//  `id` bigint(20) NOT NULL AUTO_INCREMENT,
//  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
//  `instance` varchar(50) DEFAULT NULL,
//  `endpoint` varchar(100) DEFAULT NULL,
//  `description` varchar(255) DEFAULT NULL,
//  `threshold` float DEFAULT NULL,
//  `instant_rate` float DEFAULT NULL,
//  `threshold_rate` float DEFAULT NULL,
//  `wholesome_rate` float DEFAULT NULL,
//  `errors` int(11) DEFAULT NULL,
//  `timing` int(11) DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=latin1
var db *sql.DB

func connect() error {
	if db != nil {
		return errors.Err("db connection already initialized")
	}
	var err error
	password := os.Getenv("BENCHMARK_PASSWORD")
	host := os.Getenv("BENCHMARK_HOST")
	user := os.Getenv("BENCHMARK_USER")
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/benchmark", user, password, host))
	if err != nil {
		return errors.Err(err)
	}
	return nil
}

type Results struct {
	Instance      string
	Endpoint      string
	Description   string
	Tolerance     int
	InstantRate   float64
	ThresholdRate float64
	WholesomeRate float64
	Errors        int
	Timing        int64
}

func StoreResults(instance string, endpoint string, description string, results Results) error {
	if db == nil {
		err := connect()
		if err != nil {
			return err
		}
	}

	_, err := db.Query("INSERT INTO benchmark.results(`instance`,`endpoint`,`description`,`threshold`,`instant_rate`,`threshold_rate`,`wholesome_rate`,`errors`,`timing`)"+
		" values(?,?,?,?,?,?,?,?,?)", instance, endpoint, description, results.Tolerance, results.InstantRate, results.ThresholdRate, results.WholesomeRate, results.Errors, results.Timing)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}
