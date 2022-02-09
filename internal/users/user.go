package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"hz.code/neugls/ads/internal/database"
)

//User 用户授权信息
type User struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastLogin int64  `json:"last_login"`
}

//Load load User by column
func (u *User) Load(column string, value interface{}) error {
	if column == "id" {
		return database.DBGetByIDWithMemCache("userinfo", u, "select * from #__users where id=?", value.(uint))
	}

	return database.Get(u, fmt.Sprintf("select * from #__users where %s=?", column), value)
}

func (u *User) Add() error {
	_, err := database.Insert("insert into #__users(username, password, name, last_login) values(?, ?, ?, ?)", u.Username, u.Password, u.Nickname, u.LastLogin)
	return err
}

func (u *User) SaveToDatabase() error {
	err := database.Update("update #__users set username=?, password=?, name=?, last_login=? where id=?", u.Username, u.Password, u.Nickname, u.LastLogin, u.ID)
	return err
}

func (u *User) SetPassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pwd)
	return nil
}

//CheckPassword 检查密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
