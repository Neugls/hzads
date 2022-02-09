package database

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	_ "github.com/go-sql-driver/mysql" //used link with mysql drive
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	migrate "github.com/rubenv/sql-migrate"
	"hz.code/neugls/ads/emd"
	"hz.code/neugls/ads/internal/config"
	_ "modernc.org/sqlite" //sqlite 3 import
)

var db *sqlx.DB
var firstUse = false

//IsFirstUse 判断是否为第一次使用
func IsFirstUse() bool {
	if firstUse {
		return firstUse
	}

	return false
}

//Setup setup the models
func Setup() error {

	dbf := path.Join(config.V.DataDir, config.V.DatabaseName)
	if _, e := os.Stat(dbf); e != nil {
		firstUse = true
		if f, e := os.Create(dbf); e == nil {
			f.Close()
		}
	}

	db = sqlx.MustConnect("sqlite", dbf)

	if db == nil {
		panic("database can not load, please check the configuation.")
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	//do the database migration
	return migration()
}

func migration() error {
	sqlDir, err := fs.Sub(emd.ResSQLs, "assets/sqls")
	if err != nil {
		return err
	}

	migrationSource := &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(sqlDir),
	}

	migrations, err := migrate.Exec(db.DB, "sqlite3", migrationSource, migrate.Up)
	if err != nil {
		log.Printf("migration d  atabase fail: %s", err)
	} else {
		log.Printf("migration database success with %d", migrations)
	}
	return err
}

//Close Close the db
func Close() {
	if db != nil {
		db.Close()
	}

}

//DB get the db
func DB() *sqlx.DB {
	return db
}

//Prefix change the relative sql to real sql with prefix
func Prefix(str string) string {
	return strings.Replace(str, "#__", config.V.TablePrefix, -1)
}

//UnPrefix change the real sql with prefix to relative one
func UnPrefix(str string) string {
	return strings.Replace(str, config.V.TablePrefix, "#__", 1)
}
