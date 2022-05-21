package infrastructure

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"strconv"
	"time"

	"tungnt/emmployee_manage/pkg/share/utils"
)

func NewLogOutput() io.Writer {
	maxSize, err := strconv.Atoi(os.Getenv("MAX_LOG_FILE_SIZE"))
	if err != nil {
		maxSize = 3
	}

	maxAge, err := strconv.Atoi(os.Getenv("MAX_AGE"))
	if err != nil {
		maxAge = 3
	}

	maxBackup, err := strconv.Atoi(os.Getenv("MAX_BACKUP"))
	if err != nil {
		maxBackup = 2
	}

	rollingOuput := lumberjack.Logger{
		Filename:   os.Getenv("LOG_FILE_NAME"),
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
		LocalTime:  false,
		Compress:   false,
	}

	return io.MultiWriter(os.Stdout, &rollingOuput)
}

func NewLogger() *logrus.Logger {
	logusLogger := logrus.New()
	if os.Getenv("STAGE") == "DEV" {
		logusLogger.SetOutput(os.Stdout)
		logusLogger.SetLevel(logrus.DebugLevel)
		logusLogger.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	} else {
		logusLogger.SetOutput(NewLogOutput())
		logusLogger.SetLevel(logrus.InfoLevel)
		logusLogger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logusLogger
}

type ContextHook struct {
	fields []string
	level  logrus.Level
}

func NewContextHook(level logrus.Level, keys ...string) *ContextHook {
	ctxHook := ContextHook{}
	ctxHook.AddKey(keys)
	ctxHook.SetLevel(level)
	return &ctxHook
}

func (c *ContextHook) AddKey(keys []string) {
	for _, key := range keys {
		if _, ok := utils.IsContain(c.fields, key); !ok {
			c.fields = append(c.fields, key)
		}
	}
}

func (c *ContextHook) SetLevel(level logrus.Level) {
	c.level = level
}

func (c *ContextHook) Levels() []logrus.Level {
	for i, level := range logrus.AllLevels {
		if level == c.level {
			return logrus.AllLevels[:i+1]
		}
	}

	return logrus.AllLevels
}

func (c *ContextHook) getValues(entry *logrus.Entry) {
	loggerCtx := entry.Context
	if loggerCtx != nil {
		for _, field := range c.fields {
			if v := loggerCtx.Value(field); v != "" {
				entry.WithField(field, v)
			}
		}
	}
}

func (c *ContextHook) Fire(entry *logrus.Entry) error {
	c.getValues(entry)
	return nil
}

type GormLogger struct {
	logger *logrus.Logger
}

func NewGormLogger(logger *logrus.Logger) logger.Interface {
	return &GormLogger{logger}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Infof(s, i)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Warnf(s, i)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.WithContext(ctx).Errorf(s, i)
}

func (g *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	sql, numberRows := fc()
	durations := time.Since(begin)
	entry := g.logger.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"Durations":         durations.String(),
			"SQL":               sql,
			"Number_row_effect": numberRows,
		})

	if err != nil {
		entry.WithField("Error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			entry.Info("Performed SQL")
		} else {
			entry.Error("Failed to performed sql")
		}
	} else {
		entry.Info("Performed SQL")
	}
}

type MigrateLogger struct {
	logger *logrus.Logger
}

func NewMigrateLogger(logger *logrus.Logger) *MigrateLogger {
	return &MigrateLogger{logger}
}

func (m *MigrateLogger) Printf(format string, v ...interface{}) {
	m.logger.WithField("service", "database").Infof(format, v)
}

func (m *MigrateLogger) Verbose() bool {
	return true
}
