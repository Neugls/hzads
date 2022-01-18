package database

import (
	"fmt"
	"reflect"
	"time"

	"hz.code/hz/golib/service/memservice"
)

//Select select data
func Select(dest interface{}, query string, args ...interface{}) error {
	return db.Select(dest, Prefix(query), args...)
}

//Get get object from db
func Get(dest interface{}, query string, args ...interface{}) error {
	return db.Get(dest, Prefix(query), args...)
}

//DBGetByIDWithMemCache DBGetByIDWithMemCache
func DBGetByIDWithMemCache(key string, dest interface{}, query string, id uint) error {
	mk := fmt.Sprintf("%s-%d", key, id)
	if v, e := memservice.GetValue(mk); e == nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v).Elem())
		//dest = v
		return nil
	}

	err := db.Get(dest, Prefix(query), id)
	if err == nil {
		memservice.SetValue(mk, dest, time.Hour)
		return nil
	}
	return err
}
