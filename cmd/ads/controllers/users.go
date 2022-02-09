package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hz.code/hz/golib/language"
	"hz.code/hz/golib/utils"
	"hz.code/neugls/ads/internal/config"
	"hz.code/neugls/ads/internal/users"
)

//CRSFToken CRSFToken
const CRSFToken = "39kasdfn"

//FlashMsg FlashMsg
const FlashMsg = "flashmessage"

//LastLoginUsername LastLoginUsername
const LastLoginUsername = "lastloginusername"

//LOGINED LOGINED
const LOGINED = "loginuser.id"

//GetLogin 登录
func GetLogin(c *gin.Context) {
	token := utils.GetRandomString(12)
	session := sessions.Default(c)
	session.Set(CRSFToken, token)
	session.Save()

	flashMsg, _ := session.Get(FlashMsg).(string)
	lastUsername, _ := session.Get(LastLoginUsername).(string)
	ret := c.Query("return")
	if ret == "" {
		ret = "/"
	}
	if logined, _ := IsLogined(session); logined {
		log.Println("already logined, return to ", ret)
		c.Redirect(http.StatusSeeOther, ret)
		return
	}

	log.Printf("token: %s, lastUsername: %s, flashMsg: %s\n", token, lastUsername, flashMsg)

	c.HTML(http.StatusOK, "backend.login.html", gin.H{
		"title":        language.I18nDef("login to admin", "登录管理后台"),
		"token":        token,
		"flashMsg":     flashMsg,
		"lastUsername": lastUsername,
		"appname":      config.V.AppName,
		"return":       ret,
	})
	session.Set(FlashMsg, "")
	session.Save()
}

//PostLogin 执行登录
func PostLogin(c *gin.Context) {
	session := sessions.Default(c)
	crsfToken, _ := session.Get(CRSFToken).(string)

	username := c.PostForm("username")
	password := c.PostForm("password")
	redirect := c.PostForm("return")

	log.Printf("username: %s, password: %s, redirect: %s\n", username, password, redirect)

	if c.PostForm("token") != crsfToken {
		ShowLogin(c, redirect, language.I18nDef("invalidFormToken", "非法的提交，请返回重试"))
		return
	}

	session.Set(LastLoginUsername, username)
	session.Save()

	user := users.User{}
	if e := user.Load("username", username); e != nil {
		log.Printf("load user by username:%s fail:%s", username, e)
	}

	if user.ID <= 0 {
		ShowLogin(c, redirect, language.I18nDef("loginError", "登录失败, 请检查您的用户名和密码 "))
		return
	}

	if !user.CheckPassword(password) {
		ShowLogin(c, redirect, language.I18nDef("loginError", "登录失败, 请检查您的用户名和密码 "))
		return
	}

	Login(session, &user)

	if redirect == "" {
		redirect = "/"
	}
	c.Redirect(http.StatusSeeOther, redirect)
}

//ShowLogin 显示登录界面
func ShowLogin(c *gin.Context, loginReturn string, msg string) {
	session := sessions.Default(c)

	if msg != "" {
		session.Set(FlashMsg, msg)
		session.Save()
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/login%s", utils.IIF(loginReturn == "", "", "?return="+url.QueryEscape(loginReturn)).(string)))
}

//GetLogout 退出登录
func GetLogout(c *gin.Context) {
	session := sessions.Default(c)

	islogined, _ := IsLogined(session)
	if islogined {
		Logout(session)
		session.Clear()
	}

	ShowLogin(c, "", "")
}

//IsLogined 判断当前是否已经登录了，并返回所对应的uid
func IsLogined(s sessions.Session) (logined bool, user *users.User) {
	uidinf := s.Get(LOGINED)
	if uidinf == nil {
		log.Println("get logined tag fail from session")
		return false, nil
	}
	user = &users.User{}
	id := uidinf.(uint)
	if e := user.Load("id", id); e != nil {
		log.Printf("load user by id: %d fail: %s\n", id, e.Error())
		return false, nil
	}

	return true, user
}

//LoginedUser return the logined user
func LoginedUser(c *gin.Context) (user *users.User) {
	s := sessions.Default(c)
	logined, u := IsLogined(s)
	if logined {
		return u
	}
	return nil
}

//Login 设置登录状态
func Login(s sessions.Session, user *users.User) {
	s.Set(LOGINED, user.ID)
	s.Set("loginuser.nickname", user.Nickname)
	s.Save()

	user.LastLogin = time.Now().Unix()
	if err := user.SaveToDatabase(); err != nil {
		log.Printf("save user: %d last login time fail: %s\n", user.ID, err.Error())
	}
}

//Logout 设置不登录
func Logout(s sessions.Session) {
	s.Set(LOGINED, nil)
	s.Set("nickname", nil)
	s.Save()
}

//UserAuthMiddleware UserAuthMiddleware
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if logined, user := IsLogined(session); logined {
			c.Set("user", user)
			c.Next()
		} else {
			log.Println("user not logined yet at admin middleware")
			NeedLoginFirst(c)
			c.Abort()
		}
	}
}

func APIAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		session := sessions.Default(c)
		if logined, user := IsLogined(session); logined {
			c.Set("user", user)
			c.Next()
		} else {
			c.JSON(401, map[string]interface{}{
				"error": "401 Unauthozied",
			})
			c.Abort()
		}
	}
}

// //AdminAuthMiddleware AdminAuthMiddleware
// func AdminAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session := sessions.Default(c)
// 		if logined, user := IsLogined(session); logined {
// 			c.Set("user", user)

// 			if acl.HaveAccess(user, acl.ResAdminDashboard, acl.CRUDRead) {
// 				c.Next()
// 			} else {
// 				c.String(http.StatusOK, "you have not rights to view this page")
// 			}

// 		} else {
// 			NeedLoginFirst(c)
// 			c.Abort()
// 		}
// 	}
// }

//NeedLoginFirst need login first
func NeedLoginFirst(c *gin.Context) {
	ShowLogin(c, c.Request.RequestURI, language.I18nDef("please login first", "请先登录后操作"))

}
