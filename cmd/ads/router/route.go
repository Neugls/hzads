package router

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	realip "github.com/thanhhh/gin-gonic-realip"
	"hz.code/neugls/ads/cmd/ads/controllers"
	"hz.code/neugls/ads/emd"
	"hz.code/neugls/ads/internal/config"
)

func Build() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery(), realip.RealIP())

	//static
	fe, _ := fs.Sub(emd.ResWeb, "assets/web/static")
	r.StaticFS("/static", http.FS(fe))

	//uploads
	r.Static("/uploads", path.Join(config.V.DataDir, "uploads"))

	//html template
	ft, fterr := fs.Sub(emd.ResWeb, "assets/web/templates")
	fmt.Printf("ft %v, fterr %v\n", ft, fterr)
	if fterr != nil {
		panic("template dir not found," + fterr.Error())
	}
	htmlrender := newHTMLRender()
	htmlrender.Init(ft, template.FuncMap{})
	r.HTMLRender = htmlrender

	r.GET("/", controllers.Index)
	r.GET("/ads", controllers.GetAds)

	// session
	store := memstore.NewStore([]byte(config.V.Secret))
	store.Options(sessions.Options{
		MaxAge:   1800,
		HttpOnly: true,
		Path:     "/",
	})
	r.Use(sessions.Sessions("hzads", store))

	r.GET("/login", controllers.GetLogin)
	r.POST("/login", controllers.PostLogin)

	r.GET("/init", controllers.Init)
	r.POST("/init", controllers.PostInit)

	admin := r.Group("/admin", controllers.UserAuthMiddleware())
	admin.GET("/", controllers.Admin)
	admin.GET("/ads", controllers.Ads)
	admin.POST("/ads", controllers.AddAd)
	admin.POST("/ad/:id", controllers.UpdateAd)
	admin.GET("/logout", controllers.GetLogout)
	return r
}
