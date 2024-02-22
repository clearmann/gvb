package g

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Mode          string // debug | release
		Port          string
		DbType        string // mysql | sqlite
		DbAutoMigrate bool   // 是否自动迁移数据库表结构
		DbLogMode     string // silent | error | warn | info
		StartTime     string
		MachineID     int64
	}
	Mysql struct {
		Host     string // 服务器地址
		Port     string // 端口
		DSN      string // 高级配置
		Dbname   string // 数据库名
		Username string // 数据库用户名
		Password string // 数据库密码
	}
	Redis struct {
		DB       int    // 指定 Redis 数据库
		Addr     string // 服务器地址:端口
		Password string // 密码
	}
	Log struct {
		DebugFileName string `json:"debugFileName" yaml:"debug_filename"`
		InfoFileName  string `json:"infoFileName" yaml:"info_filename"`
		WarnFileName  string `json:"warnFileName" yaml:"warn_filename"`
		MaxSize       int    `json:"maxsize" yaml:"maxsize"`
		MaxAge        int    `json:"max_age" yaml:"max_age"`
		MaxBackups    int    `json:"max_backups" yaml:"max_backups"`
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
	JWT struct {
		Secret string
		Expire int
		Issuer string
	}
}

var Conf *Config

func ReadConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}
	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}
	log.Println("配置文件内容加载成功.....")
	return Conf
}
