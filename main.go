package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/log"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/model"
	"showta.cc/app/system/router"
)

func main() {
	conf.InitConf()
	log.InitCore(conf.AppConf.Log)

	model.InitDb(conf.AppConf.Database)
	logic.Init()
	gin.SetMode(gin.ReleaseMode)

	routerInit := router.InitRouter()
	endPoint := fmt.Sprintf("%s:%d", conf.AppConf.Server.Host, conf.AppConf.Server.Port)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routerInit,
	}

	var err error
	if conf.AppConf.Server.Https {
		log.StdInfof("start https server listening %s", endPoint)
		sslCertFile := conf.AbsPath(conf.AppConf.Server.SSLCertPem)
		sslKeyFile := conf.AbsPath(conf.AppConf.Server.SSLKeyPem)
		err = server.ListenAndServeTLS(sslCertFile, sslKeyFile)
	} else {
		log.StdInfof("start http server listening %s", endPoint)
		err = server.ListenAndServe()
	}

	if err != nil {
		log.StdError(err)
		os.Exit(0)
	}
}
