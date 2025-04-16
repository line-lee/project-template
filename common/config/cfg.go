package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

const (
	// 服务名称：规范命名，可以清晰知道服务之间的层级调用

	WebGateway = "web"
	//......其他网关可以继续命名 ***Gateway，例如:OpenGateway

	VarietyService = "variety"
	TripService    = "trip"

	// ........ 最底层服务，假如之后会将 core/service/chimera/bottom 这部分拆成服务，可以命名 ****Mico,例如：FusionMico
	// FusionMico = "fusion"
)

type Config struct {
	Core        Core        `yaml:"core" json:"core"`
	MysqlConfig MysqlConfig `yaml:"mysql_config" json:"mysql_config"`
	RedisConfig RedisConfig `yaml:"redis_config" json:"redis_config"`
	KafkaConfig KafkaConfig `yaml:"kafka_config" json:"kafka_config"`
	QiNiuConfig QiNiuConfig `yaml:"qiniu_config" json:"qi_niu_config"`

	MysqlClient   *sql.DB              `yaml:"-" json:"-"`
	RedisClient   *redis.Client        `yaml:"-" json:"-"`
	KafkaProducer *kafka.Producer      `yaml:"-" json:"-"` // 生产者
	KafkaConsumer *kafka.Consumer      `yaml:"-" json:"-"` // 消费者
	ApiLimit      map[string]*ApiLimit `yaml:"-" json:"-"`
	BaiduLocation string               `yaml:"-" json:"-"` // 百度城市编码json信息
}

type Core struct {
	Gateways map[string]Register `json:"gateways" json:"gateways"`
	Services map[string]Register `yaml:"services" json:"services"`
	Micos    map[string]Register `yaml:"micos" json:"micos"`
}
type Register struct {
	Url  string `yaml:"-" json:"url"`
	Http int    `yaml:"http" json:"http"`
	Grpc int    `yaml:"grpc" json:"grpc"`
}

type MysqlConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type KafkaConfig struct {
	Servers string `yaml:"servers" json:"servers"` // 接入点的IP地址以及端口
}

type QiNiuConfig struct {
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Bucket    string `yaml:"bucket" json:"bucket"`
}

type ApiLimit struct {
	Path     string `json:"path"`
	Method   string `json:"method"`
	Describe string `json:"describe"`
	MenuId   int64  `json:"menu_id,omitempty"`
	PageId   int64  `json:"page_id,omitempty"`
	ButtonId int64  `json:"button_id,omitempty"`
}

var cfg *Config

func Info() *Config { return cfg }

func Init() *Config {
	var configFileName = "config_dev.yaml"
	var apiLimitFileName = "api_limit.json"
	var baiduLocationFileName = "baidu_code.json"
	if os.Getenv("TRIPPORTAL_ENVIRONMENT") == "PREPARE" {
		configFileName = "./config_prev.yaml"
		apiLimitFileName = "./api_limit.json"
		baiduLocationFileName = "./baidu_code.json"
	}
	if os.Getenv("TRIPPORTAL_ENVIRONMENT") == "PRODUCTION" {
		configFileName = "./config_prod.yaml"
		apiLimitFileName = "./api_limit.json"
		baiduLocationFileName = "./baidu_code.json"
	}
	// 配置文件
	configFile, err := os.OpenFile(configFileName, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("读取配置文件 os OpenFile err", err)
		os.Exit(-1)
	}
	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			fmt.Println("读取配置文件 file Close err", err)
			os.Exit(-1)
		}
	}(configFile)
	cfb, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Println("读取配置文件 io ReadAll err", err)
		os.Exit(-1)
	}
	cfg = new(Config)
	err = yaml.Unmarshal(cfb, &cfg)
	if err != nil {
		fmt.Println("读取配置文件 yaml Unmarshal err", err)
		os.Exit(-1)
	}
	for k, register := range cfg.Core.Gateways {
		register.Url = fmt.Sprintf("%s-gateway", k)
		cfg.Core.Gateways[k] = register
	}
	for k, register := range cfg.Core.Services {
		register.Url = fmt.Sprintf("%s-service", k)
		cfg.Core.Services[k] = register
	}
	for k, register := range cfg.Core.Micos {
		register.Url = fmt.Sprintf("%s-mico", k)
		cfg.Core.Micos[k] = register
	}
	// api限制文件
	apiLimitFile, err := os.OpenFile(apiLimitFileName, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("读取api限制文件 os OpenFile err", err)
		os.Exit(-1)
	}
	defer func(apiLimitFile *os.File) {
		err = apiLimitFile.Close()
		if err != nil {
			fmt.Println("读取api限制文件 file Close err", err)
			os.Exit(-1)
		}
	}(apiLimitFile)
	alb, err := io.ReadAll(apiLimitFile)
	if err != nil {
		fmt.Println("读取api限制文件 io ReadAll err", err)
		os.Exit(-1)
	}
	limits := make([]*ApiLimit, 0)
	err = json.Unmarshal(alb, &limits)
	if err != nil {
		fmt.Println("读取api限制文件 json Unmarshal err", err)
		os.Exit(-1)
	}
	alm := make(map[string]*ApiLimit)
	for _, limit := range limits {
		key := fmt.Sprintf("%s$####$%s", limit.Method, limit.Path)
		alm[key] = limit
	}
	cfg.ApiLimit = alm
	// 百度省市县三级数据
	baiduLocationFile, err := os.OpenFile(baiduLocationFileName, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("读取 baidu location 文件 os OpenFile err", err)
		os.Exit(-1)
	}
	defer func(baiduLocationFile *os.File) {
		err = baiduLocationFile.Close()
		if err != nil {
			fmt.Println("读取 baidu location 限制文件 file Close err", err)
			os.Exit(-1)
		}
	}(baiduLocationFile)
	blb, err := io.ReadAll(baiduLocationFile)
	if err != nil {
		fmt.Println("读取 baidu location 限制文件 io ReadAll err", err)
		os.Exit(-1)
	}
	cfg.BaiduLocation = string(blb)
	// 打印配置文件
	cfg.print()
	return cfg
}

func (c *Config) Open(fs ...func(c *Config)) {
	for _, f := range fs {
		f(c)
	}
}

func (c *Config) print() {
	config := Config{
		Core:        c.Core,
		MysqlConfig: c.MysqlConfig,
		RedisConfig: c.RedisConfig,
		KafkaConfig: c.KafkaConfig,
		QiNiuConfig: c.QiNiuConfig,
	}
	configBytes, _ := json.MarshalIndent(config, "", "    ")
	fmt.Println("=================================START 配置文件 START======================================")
	fmt.Println(string(configBytes))
	fmt.Println("=================================END 配置文件 END======================================")
}
