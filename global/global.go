package global

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	IsLocal bool  `yaml:"isLocal"`
	Etcd    Etcd  `yaml:"etcd"`
	Redis   Redis `yaml:"redis"`
	OTel    OTel  `yaml:"otel"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}
type Etcd struct {
	Addr              string        `yaml:"addr"`
	AuthEnable        bool          `yaml:"authEnable"`
	Username          string        `yaml:"username"`
	Password          string        `yaml:"password"`
	MaxIdleTimeout    time.Duration `yaml:"maxIdleTimeout"`
	MinIdlePerAddress int           `yaml:"minIdlePerAddress"`
}
type OTel struct {
	Enable   bool   `yaml:"enable"`
	Endpoint string `yaml:"endpoint"`
	Headers  string `yaml:"headers"`
}

func Init() {
	once.Do(initConf)
}
func initConf() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(workDir, "global/global.yaml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		panic(err)
	}
	if !conf.IsLocal {
		return
	}
	conf.Etcd.Addr = strings.Replace(conf.Etcd.Addr, "etcd", "localhost", 1)
	conf.Redis.Addr = strings.Replace(conf.Redis.Addr, "redis", "localhost", 1)
}
func GetConf() *Config {
	once.Do(initConf)
	return conf
}
