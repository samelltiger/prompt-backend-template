// config/config.go
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	NewAPI struct {
		Domain   string `yaml:"domain"`
		AdminKey string `yaml:"admin_key"`
	} `yaml:"new_api"`

	Database struct {
		GatewayDSN string `yaml:"gateway_dsn"`
		NewAPIDSN  string `yaml:"new_api_dsn"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	RateLimit struct {
		BillingQueryLimit int `yaml:"billing_query_limit"` // 每分钟查询次数
		LogQueryLimit     int `yaml:"log_query_limit"`     // 每分钟日志查询次数
	} `yaml:"rate_limit"`

	Log struct {
		RetentionDays int `yaml:"retention_days"` // 日志保留天数
	} `yaml:"log"`

	ModelMapping map[string][]string `yaml:"model_mapping"` // 模型映射关系
	Logger       Logger              `yaml:"logger"`
	OSS          OSS                 `yaml:"oss"`
}

var AppConfig Config

type Logger struct {
	Level      string `yaml:"level"` // debug, info, warn, error
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxsize"`    // MB
	MaxBackups int    `yaml:"maxbackups"` // 文件个数
	MaxAge     int    `yaml:"maxage"`     // 天数
	Compress   bool   `yaml:"compress"`   // 是否压缩
}

type OSS struct {
	Aliyun       AliyunOSS        `yaml:"aliyun"`
	OSSProxySrvs []OSSProxyServer `yaml:"oss_proxy_srv"`
	OssCacheDir  string           `yaml:"oss_cache_dir"`
}

type OSSProxyServer struct {
	IP     string `yaml:"ip"`
	Weight int    `yaml:"weight"`
}

type AliyunOSS struct {
	AccessID  string `yaml:"SZ_ALIYUN_ACCESS_ID"`
	AccessKey string `yaml:"SZ_ALIYUN_ACCESS_KEY"`
	Bucket    string `yaml:"SZ_ALIYUN_BUCKET"`
	Endpoint  string `yaml:"SZ_ALIYUN_ENDPOINT"`
	SSL       bool   `yaml:"SZ_ALIYUN_SSL"`
	IsCNAME   bool   `yaml:"SZ_ALIYUN_IS_CNAME"`
	Debug     bool   `yaml:"SZ_ALIYUN_DEBUG"`
	ImgPath   string `yaml:"IMG_PATH"`
	BookPath  string `yaml:"BOOK_PATH"`
}

func LoadConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	return err
}
