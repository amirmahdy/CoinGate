package db

import (
	"database/sql"
	"os"
	"testing"
	"utils"

	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")
	if err != nil {
		panic(err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	testStore = NewSQLStore(testDB)
	os.Exit(m.Run())
}
