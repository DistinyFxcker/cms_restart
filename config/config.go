package config

import (
	"encoding/json"
	"os"
	"path"
	//"upper.io/db.v3/gosqlite"
	"upper.io/db.v3/sqlite"
)

const (
	Version = "V2.0.18"

	UsernameLengthMin = 4

	UserPasswordLengthMin = 6

	JwtHMACTokenKey = "8yFGJs4qFbWmuf3E1tlIfbEHZKo78QMm"

	JwtECCTokenKeySeed = "shijianfeichuan"

	DefaultAdminServePort = 25824

	// // DefaultP2PPort APK分发P2P分享UDP端口;
	// DefaultP2PPort = 46565
	//ServePort
	DefaultServePort = 5824

	// DefaultDataDir 默认配置文件路径;
	DefaultDataDir = "./data"

	// ForumBaseUrl = "http://192.168.2.181:37685"
	ForumBaseUrl = "http://bbs.jilangtv.com:37685"

	// ForumTopicsUrl 论坛url基础路径;
	ForumTopicsUrl   = ForumBaseUrl + "/api/v1/topics/details2"
	NewVersionUrlFmt = "http://jilangtv.com:10000/api/version/%s" //ForumBaseUrl + "/api/v1/versioninfos/code/%s"
)
var cfg = &Config{}

type Config struct {
	Port 		int `json:"port"`
	AdminPort 	int `json:"adminport"`

	DBUsing 	string `json:"-"`

	DBType 		string `json:"dbType"`

	DBHost 		string `json:"dbHost"`
	DBPort 		int `json:"dbport"`
	DBUser 		string `json:"dbUser"`
	DBPassword 	string `json:"dbPassword"`
	DBName 		string `json:"dbName"`

	DataDir 	string `json:"-"`
}

//设置数据库和数据文件存储位置
func Global(dataDir string) *Config{
	cfg.DataDir = dataDir
	return cfg
}

//加载基础配置  默认配置文件为cms.conf
func (c *Config) Load() error{
	cfgPath := path.Join(cfg.DataDir, "cms.conf")
	return cfg.loadFromFile(cfgPath)
}
//加载数据
func (c *Config) loadFromFile(cfgPath string)error{
	cfgFile ,err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer cfgFile.Close()

	if err := json.NewDecoder(cfgFile).Decode(c); err != nil{
		return err
	}

	return nil
}

//保存最新配置信息
func (c *Config) Save() error{
	cfgPath :=path.Join(cfg.DataDir,"cms.conf")
	return cfg.saveToFile(cfgPath)
}

//保存文件配置信息
func (c *Config) saveToFile(cfgPath string) error{
	cfgFile , err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	//json.NewEncoder序列化
	encoder := json.NewEncoder(cfgFile)
	encoder.SetIndent("","\t")

	if err := encoder.Encode(c);err != nil{
		return err
	}

	return nil
}
//验证参数
func (c *Config) Validate() error {
	c.Normalize()

	return nil
}
//验证未传入的数据使用默认数据
func (c *Config) Normalize(){
	if c.Port == 0{
		c.Port = DefaultServePort	//设置APP默认端口
	}

	if c.AdminPort == 0 {
		c.AdminPort = DefaultAdminServePort //设置管理CMS默认端口
	}

	if c.DBType == ""{
		c.DBType = sqlite.Adapter //默认数据库为sqlite
	}

	if c.DBHost ==  ""{
		c.DBHost = "127.0.0.1" //默认数据库地址为本地
	}

	if c.DBPort == 0{
		c.DBPort = 3306 //默认数据库接口为3306
	}

	if c.DBName == ""{
		c.DBName = "cms"
	}

	if c.DataDir == ""{
		c.DataDir = DefaultDataDir
		_ = os.MkdirAll(c.DataDir,os.ModeDir|os.ModePerm) //打开全部相关目录 ModeDir 文件夹访问模式 ModePerm创建权限为最高0777
	}

}
//设置默认配置
func (c *Config) SetDefault(){
	c.Port = DefaultServePort
	c.AdminPort = DefaultAdminServePort
	c.DBType = sqlite.Adapter
	c.DBHost = "127.0.0.1"
	c.DBPort = 3306
	c.DBUser = "root"
	c.DBPassword = ""
	c.DBName = "cms"
	c.DataDir = DefaultDataDir
}