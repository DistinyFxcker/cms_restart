package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/alecthomas/kingpin.v2"
	"textprojcet/config"
)

var (
	port = kingpin.Flag("port","监听端口").Short('P').Int()
	adminPort = kingpin.Flag("aport","管理后台监听端口").Short('A').Int()
	dbHost = kingpin.Flag("dbhost","数据库地址").Short('i').Default("").String()
	dbPort = kingpin.Flag("dbport","数据库端口").Short('p').Default("0").Int()
	dbUser = kingpin.Flag("dbuser","数据库用户名").Short('u').Default("").String()
	dbPassword = kingpin.Flag("dbpassword","数据库密码").Short('w').Default("").String()
	dataDir = kingpin.Flag("datadir","数据库文件夹").Short('d').Default(config.DefaultDataDir).String()
	version = kingpin.Flag("version","版本信息").Short('v').Default("false").Bool()
)

func main() {
	serve := gin.Default()
	serve.POST("")
	kingpin.HelpFlag.Hidden()
	kingpin.Parse()
	fmt.Printf("服务器版本:%s\n",config.Version)

	cfg := getConfig()
}

func getConfig() *config.Config{
	cfg := config.Global(*dataDir)
}