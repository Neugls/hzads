package router

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"hz.code/neugls/ads/cmd/ads/controllers"
	"hz.code/neugls/ads/emd"
	"hz.code/neugls/ads/internal/config"
)

func Build() *gin.Engine {
	r := gin.New()
	//static
	fe, _ := fs.Sub(emd.ResWeb, "static")
	r.StaticFS("/static", http.FS(fe))

	//html template
	htmlrender := newHTMLRender()
	ft, _ := fs.Sub(emd.ResWeb, "templates")
	htmlrender.Init(ft, template.FuncMap{})
	r.HTMLRender = htmlrender

	r.GET("/", controllers.Index)

	// session
	store := memstore.NewStore([]byte(config.V.Secret))
	store.Options(sessions.Options{
		MaxAge: 1800,
	})
	r.Use(sessions.Sessions("hzmingtang", store))

	r.GET("/", controllers.Main)

	admin := r.Group("/u", controllers.UserAuthMiddleware())
	admin.GET("/", controllers.Admin)
	return r
}
