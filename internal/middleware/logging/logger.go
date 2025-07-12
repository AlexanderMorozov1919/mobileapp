package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Enabled    bool   // Включено ли логирование
	Level      string // DEBUG, INFO, WARN, ERROR
	LogsDir    string // Директория для логов
	SavingDays uint   // Сколько дней хранить логи
}

type Logger struct {
	config  *Config
	logger  *log.Logger
	file    *os.File
	prefix  string
	version string
}

func NewLogger(cfg *Config, prefix, version string) *Logger {
	l := &Logger{
		config:  cfg,
		prefix:  prefix,
		version: version,
	}

	var output io.Writer = os.Stdout
	if cfg.Enabled && cfg.LogsDir != "" {
		if err := os.MkdirAll(cfg.LogsDir, 0755); err == nil {
			logFile := filepath.Join(cfg.LogsDir, time.Now().Format("2006-01-02")+".log")
			if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				l.file = file
				output = io.MultiWriter(os.Stdout, file)
			}
		}
	}

	l.logger = log.New(output, "", log.LstdFlags)

	if cfg.SavingDays > 0 {
		go l.cleanOldLogs()
	}

	return l
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		config:  l.config,
		logger:  l.logger,
		file:    l.file,
		prefix:  l.prefix + "[" + prefix + "] ",
		version: l.version,
	}
}

func (l *Logger) cleanOldLogs() {
	for range time.Tick(24 * time.Hour) {
		files, err := os.ReadDir(l.config.LogsDir)
		if err != nil {
			l.Error("Failed to read logs directory", "error", err)
			continue
		}

		cutoff := time.Now().AddDate(0, 0, int(-l.config.SavingDays))
		for _, file := range files {
			if info, err := file.Info(); err == nil && !file.IsDir() && info.ModTime().Before(cutoff) {
				if err := os.Remove(filepath.Join(l.config.LogsDir, file.Name())); err != nil {
					l.Error("Failed to delete old log file", "file", file.Name(), "error", err)
				}
			}
		}
	}
}

func (l *Logger) log(level, msg string, fields ...interface{}) {
	if !l.ShouldLog(level) {
		return
	}

	message := fmt.Sprintf("[%s] %s%s", level, l.prefix, msg)

	var builder strings.Builder
	builder.WriteString(message)
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprint(fields[i])
		val := "?"
		if i+1 < len(fields) {
			val = fmt.Sprint(fields[i+1])
		}
		builder.WriteString(fmt.Sprintf(" %s=%s", key, val))
	}

	l.logger.Println(builder.String())
}

func (l *Logger) ShouldLog(level string) bool {
	if !l.config.Enabled {
		return false
	}

	levels := map[string]int{
		"DEBUG": 4,
		"INFO":  3,
		"WARN":  2,
		"ERROR": 1,
	}

	currentLevel := levels[strings.ToUpper(l.config.Level)]
	if currentLevel == 0 {
		currentLevel = 3 // INFO по умолчанию
	}

	return levels[strings.ToUpper(level)] <= currentLevel
}

func (l *Logger) Debug(msg string, fields ...interface{}) { l.log("DEBUG", msg, fields...) }
func (l *Logger) Info(msg string, fields ...interface{})  { l.log("INFO", msg, fields...) }
func (l *Logger) Warn(msg string, fields ...interface{})  { l.log("WARN", msg, fields...) }
func (l *Logger) Error(msg string, fields ...interface{}) { l.log("ERROR", msg, fields...) }

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func InitLogger(l *Logger) {
	l.Info("Logger initialized",
		"version", l.version,
		"level", l.config.Level,
		"logsDir", l.config.LogsDir,
		"savingDays", l.config.SavingDays,
	)
}
