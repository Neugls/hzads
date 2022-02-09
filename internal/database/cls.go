package database

import (
	"github.com/jmoiron/sqlx"
	database "hz.code/hz/golib/dbo"
	"hz.code/neugls/ads/internal/config"
)

type idb uint

func (i *idb) GetDB() *sqlx.DB {
	return db
}
func (i *idb) Prefix(str string) string {
	return Prefix(str)
}

func (i *idb) UnPrefix(str string) string {
	return UnPrefix(str)
}

func (i *idb) GetPrefix() string {
	return config.V.TablePrefix
}

//GetDBO GetDBO
func GetDBO() *database.DBO {
	_idb := idb(1)
	return database.Get(&_idb)
}
