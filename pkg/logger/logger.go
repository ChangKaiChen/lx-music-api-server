package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type LogEntry struct {
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Path      string    `json:"path,omitempty"`
	Token     string    `json:"token,omitempty"`
	SongId    string    `json:"song_id,omitempty"`
	Quality   string    `json:"quality,omitempty"`
	URL       string    `json:"url,omitempty"`
	Message   string    `json:"message"`
}
type Logger struct {
	file    *os.File
	m       sync.Mutex
	level   string
	service string
}

var (
	log  *Logger
	once sync.Once
)

func Init(serviceName, level, filepath string) {
	once.Do(func() {
		log = NewLogger(serviceName, level, filepath)
	})
}
func NewLogger(serviceName, level, filepath string) *Logger {
	if serviceName == "" {
		panic("server should not be empty")
	}
	if filepath == "" {
		panic("log filepath should not be empty")
	}
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &Logger{
		file:    file,
		level:   level,
		service: serviceName,
	}
}
func GetLogger() *Logger {
	if log == nil {
		time.Sleep(time.Second)
		if log == nil {
			panic("log not initialized")
		}
	}
	return log
}
func (l *Logger) Log(path, token, songId, quality, url, msg string) {
	entry := LogEntry{
		Level:     l.level,
		Timestamp: time.Now(),
		Service:   l.service,
		Path:      path,
		Token:     token,
		SongId:    songId,
		Quality:   quality,
		URL:       url,
		Message:   msg,
	}
	l.writeEntry(entry)
}
func (l *Logger) Info(path, message string) {
	entry := LogEntry{
		Level:     "INFO",
		Timestamp: time.Now(),
		Service:   l.service,
		Path:      path,
		Message:   message,
	}
	l.writeEntry(entry)
}
func (l *Logger) Error(path, err string) {
	entry := LogEntry{
		Level:     "ERROR",
		Timestamp: time.Now(),
		Service:   l.service,
		Path:      path,
		Message:   err,
	}
	l.writeEntry(entry)
}
func (l *Logger) Errorf(path, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	entry := LogEntry{
		Level:     "ERROR",
		Timestamp: time.Now(),
		Service:   l.service,
		Path:      path,
		Message:   message,
	}
	l.writeEntry(entry)
}
func (l *Logger) writeEntry(entry LogEntry) {
	l.m.Lock()
	defer l.m.Unlock()
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	if _, err = l.file.Write(data); err != nil {
		return
	}
	if _, err = l.file.Write([]byte("\n")); err != nil {

	}
}
