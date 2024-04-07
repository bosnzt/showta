package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"showta.cc/app/internal/webdav"
	"showta.cc/app/system/logic"
)

var webdavHandler *webdav.Handler

func AddRouterWebdav(r *gin.Engine) {
	webdavHandler = &webdav.Handler{
		Prefix:     "/dav",
		FileSystem: webdav.Dir("mnt"),
		LockSystem: webdav.NewMemLS(),
	}

	g := r.Group("/dav")
	g.Use(webdavAuth)
	g.Any("/*path", webdavHandle)
	g.Any("", webdavHandle)
	g.Handle("PROPFIND", "/*path", webdavHandle)
}

func webdavAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}

	if _, err := logic.Auth(c, username, password); err != nil {
		http.Error(c.Writer, "WebDAV: need authorized!", http.StatusUnauthorized)
		c.Abort()
		return
	}

	c.Next()
}

func webdavHandle(c *gin.Context) {
	webdavHandler.ServeHTTPOverride(c.Writer, c.Request)
}
