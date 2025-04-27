package config

import (
	"context"
	"fmt"
	"github.com/ChangKaiChen/lx-music-api-server/app/tx/refresh"
	"github.com/ChangKaiChen/lx-music-api-server/global"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/consts"
	"github.com/ChangKaiChen/lx-music-api-server/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	conf         *Config
	onceConf     sync.Once
	onceInstance sync.Once
	instance     *ExpirationListener
)

type ExpirationListener struct {
	rdb      *redis.Client
	ctx      context.Context
	isPubSub bool
}
type Config struct {
	Log        Log           `yaml:"log"`
	Quality    []string      `yaml:"quality"`
	ExpireTime time.Duration `yaml:"expireTime"`
	Guid       string        `json:"guid"`
	Users      []User        `yaml:"users"`
}
type User struct {
	Uin          string       `yaml:"uin"`
	QQMusicKey   string       `yaml:"qqMusicKey"`
	RefreshLogin RefreshLogin `yaml:"refreshLogin"`
}
type RefreshLogin struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}
type Log struct {
	Level    string `yaml:"level"`
	Filepath string `yaml:"filepath"`
}

func Init() {
	go WatchConfig()
}
func GetConf() *Config {
	onceConf.Do(initConf)
	return conf
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
	go updateUser(conf.Users)
}
func GetExpirationListener() *ExpirationListener {
	onceInstance.Do(func() {
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
func updateUser(users []User) {
	log := logger.GetLogger()
	listener := GetExpirationListener()
	_, err := listener.rdb.ConfigSet(listener.ctx, "notify-keyspace-events", "Ex").Result()
	if err != nil {
		log.Errorf("", "Failed to configure Redis: %v", err)
		return
	}
	for _, user := range users {
		if !user.RefreshLogin.Enable || user.Uin == "" || user.QQMusicKey == "" {
			continue
		}
		key := consts.TxServiceName + "->" + user.Uin + "->" + user.QQMusicKey
		exists, e := listener.rdb.Exists(listener.ctx, key).Result()
		if e != nil {
			log.Errorf("", "Failed to check existence of %s: %v", key, e)
			return
		}
		if exists != 0 {
			continue
		}
		err = listener.rdb.Set(listener.ctx, key, "", time.Duration(user.RefreshLogin.Interval)).Err()
		if err != nil {
			log.Errorf("", "Failed to set key: %v", err)
			return
		}
	}
	if listener.isPubSub {
		return
	}
	go pubSub()
	listener.isPubSub = true
}
func pubSub() {
	log := logger.GetLogger()
	listener := GetExpirationListener()
	sub := listener.rdb.PSubscribe(listener.ctx, "__keyevent@5__:expired")
	defer func(pubSub *redis.PubSub) {
		err := pubSub.Close()
		if err != nil {
			log.Errorf("", "Failed to close redis PubSub: %v", err)
		}
	}(sub)
	ch := sub.Channel()
	go func() {
		for msg := range ch {
			expiredKey := strings.Split(msg.Payload, "->")
			if expiredKey[0] == consts.TxServiceName {
				var flag bool
				for _, user := range conf.Users {
					if user.Uin == expiredKey[1] && user.QQMusicKey == expiredKey[2] {
						flag = true
						break
					}
				}
				if !flag {
					continue
				}
				log.Info("", "try refreshing musicKey: uin: "+expiredKey[1]+" key: "+expiredKey[2])
				newKey := refresh.QQMusicKeyRefresh(expiredKey[1], expiredKey[2])
				if newKey == "" {
					continue
				}
				updateUserConfig(expiredKey[1], newKey)
			}
		}
	}()
	select {}
}
func updateUserConfig(uin, newKey string) {
	log := logger.GetLogger()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get current file path")
	}
	sourceDir := filepath.Dir(filename)
	configPath := filepath.Join(sourceDir, "config.yaml")
	for i := range conf.Users {
		if conf.Users[i].Uin == uin {
			conf.Users[i].QQMusicKey = newKey
			content, err := yaml.Marshal(conf)
			if err != nil {
				log.Errorf("", "Failed to marshal config: %v", err)
				return
			}
			err = os.WriteFile(configPath, content, os.ModePerm)
			if err != nil {
				log.Errorf("", "Failed to write config: %v", err)
				return
			}
			log.Info("", "uin: "+uin+" refresh success")
			return
		}
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
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("Config file changed: %s\n", event.Name)
					initConf()
				}
			case err = <-watcher.Errors:
				if err != nil {
					fmt.Println("Error:", err)
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
