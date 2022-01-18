package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hz.code/hz/golib/language"
	"hz.code/neugls/ads/internal/database"
)

//Index 用于加载基本的html
func Index(c *gin.Context) {
	if true || database.IsFirstUse() {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    language.I18nDef("welcome to use", "欢迎使用"),
			"howtouse": language.I18nDef("lets get started", "系统初始化时用户名和密码皆为admin, 更多详细使用教程请查看官方文档"),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}
