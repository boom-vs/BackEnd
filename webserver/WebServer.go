package webserver

import (
	"crm-backend/types"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path"
)

type WebServer struct {
	handler            *http.Server
	ginEngine          *gin.Engine
	dataBaseConnection *gorm.DB
}

func (ws *WebServer) log(errorMessage string) {
	log.Println(errorMessage)
}

func (ws *WebServer) checkCertificate(ctx *types.WebServerContext) {
	if ctx.CertificatePath == "" || ctx.KeyPath == "" {
		return
	}

	cert, err := tls.LoadX509KeyPair(ctx.CertificatePath, ctx.KeyPath)
	if err != nil {
		ws.log("checkCertificate: " + err.Error())
		return
	}

	ctx.TlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

func (ws *WebServer) getLetsEncrypt() {

}

func (ws *WebServer) ginLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			logMap := make(map[string]interface{})

			logMap["status_code"] = params.StatusCode
			logMap["path"] = params.Path
			logMap["method"] = params.Method
			logMap["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			logMap["remote_addr"] = params.ClientIP
			logMap["response_time"] = params.Latency.String()

			s, _ := json.Marshal(logMap)
			return string(s) + "\n"
		},
	)
}

func (ws *WebServer) getRoutes(ctx *types.WebServerContext) {
	//ws.ginEngine.Static("/assets", "./public/assets/")
	//ws.ginEngine.Static("/svg", "./public/svg/")

	sm := &SocketManager{}
	sm.Init()
	ws.ginEngine.GET("/ws", func(ginContext *gin.Context) {
		webSocketContext := &types.WebSocketContext{
			Gin:  ginContext,
			Base: ctx.DataBase,
			Env:  ctx.Env,
		}
		sm.webSocketUpdater(webSocketContext)
	})
	ws.ginEngine.Static("/assets/", path.Join(ctx.PublicFolder, "/assets", ""))
	ws.ginEngine.Static("/err", "./stderr")
	ws.ginEngine.NoRoute(func(ginContext *gin.Context) {
		ginContext.File(path.Join(ctx.PublicFolder, types.ConstDefaultPage))
	})
}

func (ws *WebServer) Start(ctx *types.WebServerContext) {

	ws.dataBaseConnection = ctx.DataBase

	log.Println("WebServer DB ", ws.dataBaseConnection)

	ws.ginEngine = gin.New()
	ws.ginEngine.Use(gin.Recovery())
	ws.ginEngine.Use(ws.ginLogger())
	ws.getRoutes(ctx)

	ws.checkCertificate(ctx)

	ws.handler = &http.Server{
		Addr:      fmt.Sprintf("%s:%d", ctx.IP, ctx.Port),
		Handler:   ws.ginEngine,
		TLSConfig: ctx.TlsConfig,
	}

	if ctx.TlsConfig != nil {
		//https
		err := ws.handler.ListenAndServeTLS("", "")
		if err != nil {
			ws.log("Start: " + err.Error())
		}
	} else {
		//http
		err := ws.handler.ListenAndServe()
		if err != nil {
			ws.log("Start: " + err.Error())
		}
	}
}
