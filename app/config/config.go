package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(multipleConfig)

type multipleConfig struct {
	*AppConfig       `mapstructure:"app"`
	*LogConfig       `mapstructure:"log"`
	*MySQLConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*RateLimitConfig `mapstructure:"rate_limit"`
}

type AppConfig struct {
	Name             string        `mapstructure:"name"`
	Mode             string        `mapstructure:"mode"`
	Version          string        `mapstructure:"version"`
	Port             string        `mapstructure:"port"`
	StartTime        string        `mapstructure:"start_time"`
	PageNum          int64         `mapstructure:"page_num"`
	PageSize         int64         `mapstructure:"page_size"`
	MachineID        int64         `mapstructure:"machine_id"`
	JwtExpireTime    time.Duration `mapstructure:"jwt_expire_time"`
	JwtSecret        string        `mapstructure:"jwt_secret"`  // JWT密钥
	JwtIssuer        string        `mapstructure:"jwt_issuer"`  // JWT签发人
	JwtSubject       string        `mapstructure:"jwt_subject"` // JWT签发对象
	UploadModel      string        `mapstructure:"uploadModel"`
	CacheExpireTime  time.Duration `mapstructure:"cache_expire_time"` // 缓存过期时间
	SnowflakeEpoch   string        `mapstructure:"snowflake_epoch"`   // 雪花算法起始时间
	ScheduleInterval time.Duration `mapstructure:"schedule_interval"` // 定时任务间隔
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	FilePath   string `mapstructure:"filepath"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         string `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RateLimitConfig struct {
	MaxRequests    int           `mapstructure:"max_requests"`    // 最大请求次数
	BanDuration    time.Duration `mapstructure:"ban_duration"`    // 限流持续时间
	WindowDuration time.Duration `mapstructure:"window_duration"` // 滑动窗口时间
}

func InitSettings() (err error) {
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.SetConfigFile("./config/config.yaml")
	// viper.AddConfigPath("./config/")
	// viper.AddConfigPath("./")
	// viper.SetConfigFile(filePath)
	var filePath string
	flag.StringVar(&filePath, "filePath", "./config/config.yaml", "配置文件")
	//解析命令行参数
	flag.Parse()
	fmt.Println(filePath)
	fmt.Println(flag.Args())
	fmt.Println(flag.NArg())
	fmt.Println(flag.NFlag())
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 初始化配置结构体指针
	Conf.AppConfig = &AppConfig{}
	Conf.LogConfig = &LogConfig{}
	Conf.MySQLConfig = &MySQLConfig{}
	Conf.RedisConfig = &RedisConfig{}
	Conf.RateLimitConfig = &RateLimitConfig{}

	// 配置信息的反序列化
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
