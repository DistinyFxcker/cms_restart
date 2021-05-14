package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"textprojcet/config"
	"textprojcet/runtime"
	"time"
	"upper.io/db.v3/lib/sqlbuilder"
	//"upper.io/db.v3/mysql"
	//"upper.io/db.v3/sqlite"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm/dialects/sqlite"
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

	cfg := getConfig() //配置初始化
	_ = os.Mkdir(cfg.DataDir,os.ModeDir|os.ModePerm)

	database := connDatabase(cfg)	//链接数据库,根据配置内的数据库类型链接数据库
	defer disconnectionDatabase(database) //加入defer在服务结束的时候关闭数据库链接
	//设置数据库类型
	cfg.DBUsing = cfg.DBType

	//打开文件储存库
	opts := badger.DefaultOptions("data/files")
	opts.SyncWrites = true
	fileDb , err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	defer fileDb.Close()

	//初始化运行 注册各种服务
	runtime := runtime.New(database,cfg.DBType,fileDb,cfg.Port,cfg.AdminPort,cfg.Port)
	if err := runtime.Init(); err != nil{
		panic(err)
	}

	if err := runtime.Start();err != nil {
		panic(err)
	}
	defer runtime.Stop()
	//
	////开启API服务 开启一个echo服务
	//apiServer := webapi.NewServer(runtime ,cfg.AdminPort , cfg.Port)
	//go func (){
	//	if err :=apiServer.Start() ; err != nil {
	//		panic(err)
	//	}
	//}()
	////defer关闭服务
	//defer apiServer.Shutdown()


	//开启web服务 	开启一个echo服务


}
//获取相关配置
func getConfig() *config.Config{
	//设置默认配置文件
	cfg := config.Global(*dataDir)
	// 从配置文件中获取配置;
	if err :=cfg.Load();err != nil {
		cfg.SetDefault()//如果已有的数据出现问题,直接调用默认配置
	}


	defer cfg.Save()
	//验证数据
	if err := cfg.Validate();err != nil{
		//由于未设置返回错误数据，此处使用默认panic
		panic(err.Error())
	}

	//设置端口
	if *port !=0 {
		cfg.Port = int(*port)
	}

	if *adminPort != 0 {
		cfg.AdminPort = *adminPort
	}

	if *dbHost != "" {
		cfg.DBHost = *dbHost
	}

	if *dbPort != 0 {
		cfg.DBPort = *dbPort
	}

	if *dbUser != "" {
		cfg.DBUser = *dbUser
	}

	if *dbPassword != "" {
		cfg.DBPassword = *dbPassword
	}

	return cfg
}

//链接数据库(现在支持两种数据库链接sqlite和mysql)
func connDatabase(cfg *config.Config) sqlbuilder.Database{
	switch cfg.DBType {
	case "" , "default" , gorm.DefaultTableNameHandler() :
		dburl ,err := sqlite.ParseURL(fmt.Sprintf("file://%s/database.db?cache=private&_journal=WAL",cfg.DataDir))
		if err != nil {
			panic(err)
		}

		_database , err := sqlite.Open(dburl)
		if err != nil{
			panic(err)
		}
		return _database
	case mysql.Adapter:
		dbUrl := mysql.ConnectionURL{
			User :		cfg.DBUser,
			Password:	cfg.DBPassword,
			Database:	cfg.DBName,
			Host:		fmt.Sprintf("%s:%d",cfg.DBHost,cfg.DBPort),
		}
		_database , err := mysql.Open(dbUrl)
		if err != nil {
			panic(err)
		}
		go func(){
			for{
				time.Sleep(1 * time.Hour)
				fmt.Printf("database ping()")
				_ = _database.Ping()
			}
		}()

		return _database
	}
	return nil
}

//在关闭服务的时候关闭数据库连接
func disconnectionDatabase(db sqlbuilder.Database){
	_ = db.Close()
}

