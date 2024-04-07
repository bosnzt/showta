package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/log"
	"showta.cc/app/web"
	"strings"
)

var (
	static    fs.FS = web.Static
	rawHtml   string
	indexHtml string
)

func Favicon(c *gin.Context) {
	if conf.SiteFavicon == "" {
		ico, err := ioutil.ReadFile("app/web/dist/favicon.ico")
		if err != nil {
			return
		}

		c.Header("Cache-Control", "public, max-age=7776000")
		c.Writer.WriteString(string(ico))
	} else {
		c.Redirect(302, conf.SiteFavicon)
	}
}

func EmbedWeb(r *gin.Engine, noRoute func(handlers ...gin.HandlerFunc)) {
	loadStatic()
	loadHtml()

	folders := []string{"css", "js", "img"}
	r.Use(func(c *gin.Context) {
		for _, folder := range folders {
			if strings.HasPrefix(c.Request.RequestURI, fmt.Sprintf("/dist/%s/", folder)) {
				c.Header("Cache-Control", "public, max-age=7776000")
			}
		}
	})

	for i, folder := range folders {
		sub, err := fs.Sub(static, folder)
		if err != nil {
			log.Errorf("can't find folder: %s", folder)
		} else {
			r.StaticFS(fmt.Sprintf("/dist/%s/", folders[i]), http.FS(sub))
		}
	}

	noRoute(func(c *gin.Context) {
		log.Debug("noRoute: ", c.Request.URL.Path)
		c.Data(http.StatusOK, "text/html;charset=utf-8", []byte(indexHtml))
	})
}

func loadHtml() {
	indexFile, err := static.Open("index.html")
	if err != nil {
		log.Errorf("failed to read index.html: %+v", err)
		return
	}

	defer func() {
		indexFile.Close()
	}()

	index, err := io.ReadAll(indexFile)
	if err != nil {
		log.Errorf("failed to read dist/index.html: %+v", err)
		return
	}

	rawHtml = string(index)
	UpdateHtml()
}

func UpdateHtml() {
	indexHtml = rawHtml
	replacements := map[string]string{
		"<title>ShowTa</title>": "<title>" + conf.SiteTitle + "</title>",
	}
	for oldStr, newStr := range replacements {
		indexHtml = strings.Replace(indexHtml, oldStr, newStr, 1)
	}
}

func loadStatic() {
	dist, err := fs.Sub(static, "dist")
	if err != nil {
		log.Errorf("failed to read dist dir: %+v", err)
		return
	}

	static = dist
}
