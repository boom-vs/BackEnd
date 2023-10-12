package types

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type WebServerContext struct {
	IP              string
	Port            int
	Domain          string
	CertificatePath string
	KeyPath         string
	AutoLetsEncrypt bool
	TlsConfig       *tls.Config
	PublicFolder    string
	DataBase        *gorm.DB
	Env             map[string]string
}

type WebSocketContext struct {
	Gin            *gin.Context
	Socket         *websocket.Conn
	Base           *gorm.DB
	Env            map[string]string
	ReceivedJSON   interface{}
	EmployeeToken  string
	LastController string

	//RawContext         *types.WebSocketContext
}

type GormContext struct {
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	AutoMigrate bool
	SSLMode     bool
}

type RequestContext struct {
	Socket        *websocket.Conn
	Base          *gorm.DB
	ReceivedData  interface{}
	Response      WebSocketPackage
	EmployeeToken string
	EmployeeId    uint
	Updater       func(Contoller string, Data []interface{})
	Env           map[string]string
}
