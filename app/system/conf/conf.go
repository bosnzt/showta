package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"showta.cc/app/lib/util"
)

type Server struct {
	Host       string `ini:"host"`
	Port       int    `ini:"port"`
	Https      bool   `ini:"https"`
	SSLCertPem string `ini:"ssl_cert_pem"`
	SSLKeyPem  string `ini:"ssl_key_pem"`
}

type Log struct {
	Enable     bool   `ini:"enable"`
	Filename   string `ini:"filename"`
	MaxSize    int    `ini:"max_size"`
	MaxBackups int    `ini:"max_backups"`
	MaxAge     int    `ini:"max_age"`
	Compress   bool   `ini:"compress"`
	Level      int    `ini:"Level"`
	Stdout     bool   `ini:"stdout"`
}

type Database struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	Dbname   string `ini:"dbname"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

type Secure struct {
	TokenExpire int    `ini:"token_expire"`
	JwtSecret   string `ini:"jwt_secret"`
	SignKey     string `ini:"sign_key"`
}

type Config struct {
	Server   `ini:"server"`
	Database `ini:"database"`
	Log      `ini:"log"`
	Secure   `ini:"secure"`
}

var (
	AppConf     = &Config{}
	AppPath     string
	IniFileName string
)

func InitConf() {
	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("get exePath error: %v", err))
	}

	AppPath = filepath.Dir(exePath)
	IniFileName = AbsPath("config.ini")
	ok, err := util.PathExist(IniFileName)
	if err != nil {
		panic(fmt.Sprintf("Read cfg file error: %v", err))
	}

	if !ok {
		genLocalConf()
		return
	}

	config, err := ini.Load(IniFileName)
	if err != nil {
		panic(fmt.Sprintf("Load cfg file fail: %v", err))
	}

	//Read and map to the structure
	err = config.MapTo(AppConf)
	if err != nil {
		panic(fmt.Sprintf("Load cfg file fail %v", err))
	}
}

func genLocalConf() {
	AppConf.Server = Server{
		Host:       "0.0.0.0",
		Port:       8888,
		SSLCertPem: "cert.pem",
		SSLKeyPem:  "key.pem",
	}

	AppConf.Log = Log{
		Enable:     true,
		Filename:   "runtime/logs/run.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
	}

	AppConf.Database = Database{
		Dbname: "runtime/data/nano.db",
	}

	AppConf.Secure = Secure{
		TokenExpire: 72,
		JwtSecret:   util.GenRandStr(16),
		SignKey:     util.GenRandStr(16),
	}

	createIniFile()
}

func createIniFile() {
	cfg := ini.Empty()
	// Format the structure into an INI file
	err := cfg.ReflectFrom(AppConf)
	if err != nil {
		fmt.Printf("Failed to reflect values: %v", err)
		return
	}

	// Write INI
	err = cfg.SaveTo(IniFileName)
	if err != nil {
		fmt.Printf("Failed to save file: %v", err)
		return
	}
}

func AbsPath(rpath string) string {
	if AppVersion == "dev" {
		return rpath
	}

	return filepath.Join(AppPath, rpath)
}
