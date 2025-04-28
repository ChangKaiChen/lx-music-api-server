package config

import (
	"context"
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	conf     *Config
	once     sync.Once
	instance *ExpirationListener
)

type Config struct {
	Log     Log      `yaml:"log"`
	Quality []string `yaml:"quality"`
	Users   []User   `yaml:"users"`
}
type User struct {
	Token        string       `yaml:"token"`
	RefreshLogin RefreshLogin `yaml:"refreshLogin"`
}
type RefreshLogin struct {
	Enable   bool          `yaml:"enable"`
	Interval time.Duration `yaml:"interval"`
}
type Log struct {
	Level    string `yaml:"level"`
	Filepath string `yaml:"filepath"`
}
type ExpirationListener struct {
	rdb      *redis.Client
	ctx      context.Context
	isPubSub bool
}

func Init() {
	go WatchConfig()
}
func GetConf() *Config {
	once.Do(initConf)
	return conf
}
func GetExpirationListener() *ExpirationListener {
	once.Do(func() {
		ctx := context.Background()
		rdb := redis.NewClient(&redis.Options{
			Addr:     global.GetConf().Redis.Addr,
			Password: global.GetConf().Redis.Password,
			DB:       5,
		})
		instance = &ExpirationListener{rdb: rdb, ctx: ctx}
	})
	return instance
}
func initConf() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get current file path")
	}
	sourceDir := filepath.Dir(filename)
	configPath := filepath.Join(sourceDir, "config.yaml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		panic(err)
	}
}
func WatchConfig() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get current file path")
	}
	sourceDir := filepath.Dir(filename)
	configPath := filepath.Join(sourceDir, "config.yaml")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err = watcher.Close()
		if err != nil {
		}
	}(watcher)
	log := logger.GetLogger()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Info("", "Config file changed: "+event.Name)
					initConf()
				}
			case err = <-watcher.Errors:
				if err != nil {
					log.Errorf("", "WatchConfig Error: %v", err)
				}
			}
		}
	}()
	if err = watcher.Add(configPath); err != nil {
		panic(err)
	}
	// 阻止 goroutine 退出
	select {}
}
