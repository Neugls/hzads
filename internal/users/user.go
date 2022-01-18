package users

import (
	"database/sql"
	"fmt"

	"hz.code/neugls/ads/internal/database"
)

//User 用户授权信息
type User struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Deleted  int    `json:"deleted"`
	State    int    `json:"state"`
	Admin    bool   `json:"admin"` //是否为管理员

	Role        uint         `json:"role"` //所属权限角色
	CreatedAt   sql.NullTime `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	LastLogin   sql.NullTime `json:"lastlogin"`
	LastLoginIP string       `json:"lastloginip"`
}

//Load load User by column
func (u *User) Load(column string, value interface{}) error {
	if column == "id" {
		return database.DBGetByIDWithMemCache("userinfo", u, "select * from #__user where id=?", value.(uint))
	}

	return database.Get(u, fmt.Sprintf("select * from #__user where %s=?", column), value)
}
