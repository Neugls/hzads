package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"hz.code/hz/golib/language"
	"hz.code/neugls/ads/internal/ads"
	"hz.code/neugls/ads/internal/config"
	"hz.code/neugls/ads/internal/database"
	"hz.code/neugls/ads/internal/users"
)

//Index 用于加载基本的html
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "front.index.html", gin.H{
		"title": language.I18nDef("welcome to use", "欢迎使用"),
	})
}

//GetAds 获取广告列表
func GetAds(c *gin.Context) {
	adss, count, err := ads.ListValid()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"error":   err.Error(),
			"message": language.I18nDef("get ads list failed", "获取广告列表失败"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  count,
		"data":   adss,
	})
}

func Init(c *gin.Context) {
	if true || database.IsFirstUse() {
		c.HTML(http.StatusOK, "backend.first.html", gin.H{
			"title":   language.I18nDef("First to use", "首次使用"),
			"appname": config.V.AppName,
			"error":   "",
		})
		return
	}
	c.AbortWithError(http.StatusNotFound, errors.New(language.I18nDef("not found", "未找到")))
}

func PostInit(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	password1 := c.PostForm("password1")

	if password != password1 {
		c.HTML(http.StatusOK, "backend.first.html", gin.H{
			"title":   language.I18nDef("First to use", "首次使用"),
			"appname": config.V.AppName,
			"error":   language.I18nDef("password not same", "密码不一致"),
		})
		return
	}

	user := users.User{}
	user.Username = username
	user.Nickname = language.I18nDef("admin", "管理员")
	user.SetPassword(password)
	if err := user.Add(); err != nil {
		c.HTML(http.StatusOK, "backend.first.html", gin.H{
			"title":   language.I18nDef("First to use", "首次使用"),
			"appname": config.V.AppName,
			"error":   language.I18nDef("Add user fail", "添加用户失败") + err.Error(),
		})
		return
	}

	ShowLogin(c, "/admin", language.I18nDef("Initilization successed, please login", "初始化成功，请登录"))
}
