package test

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq" // initialize posgresql for test
	"github.com/qiangxue/go-restful-api/internal/config"
	"github.com/qiangxue/go-restful-api/pkg/log"
	"path"
	"runtime"
	"testing"
)

var db *dbx.DB

// DB returns the database connection for testing purpose.
func DB(t *testing.T) *dbx.DB {
	if db != nil {
		return db
	}
	logger, _ := log.NewForTest()
	dir := getSourcePath()
	cfg, err := config.Load(dir+"/../../config/local.yml", logger)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db, err = dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	db.LogFunc = logger.Infof
	return db
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *dbx.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.TruncateTable(table).Execute()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

// getSourcePath returns the directory containing the source code that is calling this function.
func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
