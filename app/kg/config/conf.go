package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Log        Log           `yaml:"log"`
	Quality    []string      `yaml:"quality"`
	ExpireTime time.Duration `yaml:"expireTime"`
	Users      []User        `yaml:"users"`
}
type Log struct {
	Level    string `yaml:"level"`
	Filepath string `yaml:"filepath"`
}
type User struct {
	Token        string       `yaml:"token"`
	UserId       string       `yaml:"userid"`
	Mid          string       `yaml:"mid"`
	RefreshLogin RefreshLogin `yaml:"refreshLogin"`
}
type RefreshLogin struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

func Init() {
	go WatchConfig()
}
func GetConf() *Config {
	once.Do(initConf)
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
