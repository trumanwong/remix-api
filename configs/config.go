package configs

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GlobalConfig struct {
	HttpPort      int           `yaml:"http_port"`
	Env           string        `yaml:"env"`
	MaxThrottle   float64       `yaml:"max_throttle"`
	QiNiu         QiNiu         `yaml:"qiniu"`
	MySQL         MySQL         `yaml:"mysql"`
	Redis         Redis         `yaml:"redis"`
	Cache         Cache         `yaml:"cache"`
	Log           Log           `yaml:"log"`
	Elasticsearch Elasticsearch `yaml:"elasticsearch"`
	RabbitMQ      RabbitMQ      `yaml:"rabbitmq"`
	AliOss        AliOss        `yaml:"ali_oss"`
	Command       Command       `yaml:"command"`
	Task          Task          `yaml:"task"`
	Other         Other         `yaml:"other"`
	Wechat        Wechat        `yaml:"wechat"`
}

type MySQL struct {
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Database    string `yaml:"database"`
	TablePrefix string `yaml:"table_prefix"`
}

// QiNiu 七牛配置
type QiNiu struct {
	AccessKey   string `yaml:"access_key"`
	SecretKey   string `yaml:"secret_key"`
	Bucket      string `yaml:"bucket"`
	SpaceDomain string `yaml:"space_domain"`
	PicDomain   string `yaml:"pic_domain"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

type Cache struct {
	Prefix   string `yaml:"prefix"`
	Database int    `yaml:"database"`
}

// Log 日志
type Log struct {
	Path   string `yaml:"path"`
	MaxDay int    `yaml:"max_day"`
	Driver string `yaml:"driver"`
}

type RabbitMQ struct {
	Url string `yaml:"url"`
}

type Elasticsearch struct {
	UserName  string   `yaml:"user_name"`
	Password  string   `yaml:"password"`
	Addresses []string `yaml:"addresses"`
}

// AliOss 阿里配置
type AliOss struct {
	AccessKeyId      string `yaml:"access_key_id"`
	AccessKeySecret  string `yaml:"access_key_secret"`
	EndPoint         string `yaml:"end_point"`
	EndPointShenzhen string `yaml:"end_point_shenzhen"`
}

type Command struct {
	Python Python `yaml:"python"`
	Ffmpeg string `yaml:"ffmpeg"`
}

type Python struct {
	Command         string `yaml:"command"`
	GifRevertScript string `yaml:"gif_revert_script"`
	NokiaSmsScript  string `yaml:"nokia_sms_script"`
	RemixScript     string `yaml:"remix_script"`
}

type Task struct {
	Remix Remix `yaml:"remix"`
}

type Remix struct {
	TemplatePath    string     `yaml:"template_path"`
	TemplateFolders []string   `yaml:"template_folders"`
	Sentences       [][]string `yaml:"sentences"`
}

type Wechat struct {
	MiniProgram MiniProgram `yaml:"mini_program"`
}

type MiniProgram struct {
	AppId     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type Other struct {
	WeChatRobotUrl  string `yaml:"wechat_robot_url"`
	AllowOrigins    string `yaml:"allow_origins"`
	JwtKey          string `yaml:"jwt_key"`
	StorePath       string `yaml:"store_path"`
	AllowPictureExt string `yaml:"allow_picture_ext"`
}

var Config *GlobalConfig

var c = flag.String("c", "", "配置文件路径")

func init() {
	flag.Parse()
	var data []byte
	var err error
	if len(*c) == 0 {
		data, err = ioutil.ReadFile("config.yaml")
	} else {
		data, err = ioutil.ReadFile(*c)
	}
	if err != nil {
		log.Fatalf("Fail to read config file %v", err.Error())
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("yaml unmarshal fail %v", err.Error())
	}
	Config.Redis.Prefix += ":"
}
