package controllers

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"hz.code/hz/golib/language"
	"hz.code/hz/golib/utils"
	"hz.code/neugls/ads/internal/ads"
	"hz.code/neugls/ads/internal/config"
)

func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "backend.index.html", gin.H{
		"title":   language.I18nDef("Ads dashboard", "好智广告管理系统后台"),
		"appname": config.V.AppName,
	})
}

func Ads(c *gin.Context) {
	adsLists, count, err := ads.List(0, 1000)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  count,
		"data":   adsLists,
	})
}

func AddAd(c *gin.Context) {
	name := c.PostForm("name")
	adType := c.PostForm("type")
	content := ""

	if adType != "webpage" {
		contentFile, err := c.FormFile("content")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  "fetch upload file error: " + err.Error(),
			})
			return
		}
		content = "uploads/" + utils.GetRandomString(20) + filepath.Ext(contentFile.Filename)
		abs := path.Join(config.V.DataDir, content)
		if c, e := filepath.Abs(abs); e == nil {
			abs = c
		}
		if err := os.MkdirAll(path.Dir(abs), 0777); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  "mkdir error: " + err.Error(),
			})
			return
		}

		if err := c.SaveUploadedFile(contentFile, abs); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  "save upload file error: " + err.Error(),
			})
			return
		}

	} else {
		content = c.PostForm("content")
	}

	ad := &ads.Ads{
		Name:     name,
		Type:     adType,
		Content:  content,
		Position: "main",
		Status:   ads.AdsStatusValid,
		Sort:     0,
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}

	if err := ads.Insert(ad); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   ad,
	})
}

func UpdateAd(c *gin.Context) {
	id, err := utils.AsInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"error":  "invalid id " + err.Error(),
		})
		return
	}
	action := c.PostForm("action")

	if action == "invalid" {
		if err := ads.Invalid(uint(id)); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}

	if action == "del" {
		if err := ads.Delete(uint(id)); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}

	if action == "valid" {
		if err := ads.Valid(uint(id)); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}

}
