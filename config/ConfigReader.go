package config

import (
	"crm-backend/types"
	"log"
	"net"
	"os"
	"strconv"
)

type ConfigReader struct {
}

func (cr *ConfigReader) log(errorMessage string) {
	log.Println(errorMessage)
}

func (cr *ConfigReader) getServerIP() string {
	serverIPStr := os.Getenv("SERVER_IP")
	if serverIPStr == "" {
		return types.ConstServerIP
	}
	serverIP := net.ParseIP(serverIPStr)
	if serverIP == nil {
		return types.ConstServerIP
	}
	return serverIP.String()

}

func (cr *ConfigReader) getServerPort() int {

	serverPortStr := os.Getenv("SERVER_PORT")
	if serverPortStr == "" {
		return types.ConstServerPort
	}
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		cr.log("getServerPort: " + err.Error())
		return types.ConstServerPort
	}
	return serverPort
}

func (cr *ConfigReader) getServerDomain() string {
	serverDomain := os.Getenv("SERVER_DOMAIN")
	if serverDomain == "" {
		return types.ConstServerDefaultDomain
	}
	return serverDomain
}

func (cr *ConfigReader) getServerCertificatePath() string {
	serverCertificatePath := os.Getenv("SERVER_CRT_PATH")
	if serverCertificatePath == "" {
		return ""
	}

	if _, err := os.Stat(serverCertificatePath); err != nil {
		cr.log("getServerCertificatePath: " + err.Error())
		return ""
	}

	return serverCertificatePath
}

func (cr *ConfigReader) getServerKeyPath() string {
	serverKeyPath := os.Getenv("SERVER_KEY_PATH")
	if serverKeyPath == "" {
		return ""
	}

	if _, err := os.Stat(serverKeyPath); err != nil {
		cr.log("getServerKeyPath: " + err.Error())
		return ""
	}

	return serverKeyPath
}

func (cr *ConfigReader) getServerLetsEncrypt() bool {
	serverLetsEncrypt := os.Getenv("SERVER_LETS_ENCRYPT")
	if serverLetsEncrypt == "" {
		return false
	}
	return true
}

func (cr *ConfigReader) getServerPublicFolder() string {
	serverPublicFolder := os.Getenv("SERVER_PUBLIC")
	if serverPublicFolder == "" {
		return types.ConstServerPublicFolder
	}

	if _, err := os.Stat(serverPublicFolder); err != nil {
		cr.log("getServerPublicFolder: " + err.Error())
		return types.ConstServerPublicFolder
	}

	return serverPublicFolder
}

func (cr *ConfigReader) GetWebServerContext() *types.WebServerContext {
	return &types.WebServerContext{
		IP:              cr.getServerIP(),
		Port:            cr.getServerPort(),
		Domain:          cr.getServerDomain(),
		CertificatePath: cr.getServerCertificatePath(),
		KeyPath:         cr.getServerKeyPath(),
		AutoLetsEncrypt: cr.getServerLetsEncrypt(),
		PublicFolder:    cr.getServerPublicFolder(),
	}
}

func (cr *ConfigReader) getDBHost() string {
	dbHostStr := os.Getenv("DB_HOST")
	if dbHostStr == "" {
		return types.ConstDBHost
	}

	return dbHostStr
}

func (cr *ConfigReader) getDBPort() int {
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		return types.ConstDBPort
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		cr.log("getDbPort: " + err.Error())
		return types.ConstDBPort
	}
	return dbPort
}

func (cr *ConfigReader) getDBUser() string {
	dbUserStr := os.Getenv("DB_USER")
	if dbUserStr == "" {
		return types.ConstDBUser
	}

	return dbUserStr
}

func (cr *ConfigReader) getDBPassword() string {
	dbPasswordStr := os.Getenv("DB_PASSWORD")
	return dbPasswordStr
}

func (cr *ConfigReader) getDBBase() string {
	dbBaseStr := os.Getenv("DB_BASE")
	return dbBaseStr
}

func (cr *ConfigReader) getDBAutoMigrate() bool {
	dbAutoMigrateStr := os.Getenv("DB_AUTO_MIGRATE")

	if dbAutoMigrateStr == "" {
		return types.ConstDBAutoMigrate
	}

	dbAutoMigrate, err := strconv.Atoi(dbAutoMigrateStr)
	if err != nil {
		cr.log("getDbAutoMigrate: " + err.Error())
		return types.ConstDBAutoMigrate
	}

	if dbAutoMigrate == 0 {
		return false
	}

	return true
}

func (cr *ConfigReader) getDBSSLMode() bool {
	dbSSLModeStr := os.Getenv("DB_SSL_MODE")

	if dbSSLModeStr == "" {
		return false
	}

	dbSSLMode, err := strconv.Atoi(dbSSLModeStr)
	if err != nil {
		return false
	}

	if dbSSLMode == 0 {
		return false
	}

	return true
}

func (cr *ConfigReader) GetGormContext() *types.GormContext {
	return &types.GormContext{
		Host:        cr.getDBHost(),
		Port:        cr.getDBPort(),
		User:        cr.getDBUser(),
		Password:    cr.getDBPassword(),
		Database:    cr.getDBBase(),
		AutoMigrate: cr.getDBAutoMigrate(),
		SSLMode:     cr.getDBSSLMode(),
	}
}

func (cr ConfigReader) GetEnv() map[string]string {
	env := make(map[string]string)

	env["TMP"] = "/tmp"
	return env
}
