package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"showta.cc/app/system/log"
	"showta.cc/app/system/router/api"
	"showta.cc/app/system/router/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(log.GinLogger(), log.GinRecovery(true))
	r.Use(Cors())

	r.GET("/dist/favicon.ico", api.Favicon)
	r.GET("/preference", api.GetPreference)
	r.GET("/fd/*path", api.ProxyFile)

	pa := r.Group("", middleware.PermissiveAuth)
	api.AddRouterFile(pa)
	api.AddRouterWebdav(r)

	admin := r.Group("/admin")
	admin.POST("/login", api.UserLogin)

	ea := admin.Group("", middleware.EnforcingAuth)
	ea.GET("/menu", api.GetMenu)
	ea.GET("/user/about", api.AboutUser)
	ea.POST("/user/reset_pwd", api.ResetPwd)

	sa := admin.Group("", middleware.StrictAuth)
	api.AddRouterUser(sa)
	api.AddRouterStorage(sa)
	api.AddRouterFolder(sa)
	api.AddRouterPreference(sa)

	api.EmbedWeb(r, func(handlers ...gin.HandlerFunc) {
		r.NoRoute(handlers...)
	})

	return r
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	})
}
