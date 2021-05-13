package config

import (
	"encoding/json"
	"os"
	"path"
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
	Port int `json:"port"`
	AdminPort int `json:"adminport"`
	DBUsing int `json:"-"`
	DBType int `json:"dbType"`

	DBHost int `json:"dbHost"`
	DBPort int `json:"dbport"`
	DBUser int `json:"dbPassword"`
	DataDir string `json:"-"`
}

func Global(dataDir string) *Config{
	cfg.DataDir = dataDir
	return cfg
}
//
func (c *Config) Load() error{
	cfgPath := path.Join(cfg.DataDir, "cms.conf")
	return cfg.saveToFaile(cfgPath)
}

func (c *Config) saveToFaile(cfgPath string) error{
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