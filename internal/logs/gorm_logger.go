package logs

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct {
	logger.Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(writer logger.Writer, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &GormLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (this *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *this
	newlogger.LogLevel = level
	return &newlogger
}

// Info write info messages
func (this *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if this.LogLevel < logger.Info {
		return
	}
	Write(Log{
		Message:   fmt.Sprintf(this.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...),
		Level:     "info",
		CreatedAt: time.Now(),
	})
}

// Warn write waning messages
func (this *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if this.LogLevel < logger.Error {
		return
	}
	Write(Log{
		Message:   fmt.Sprintf(this.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...),
		Level:     LogLevelWarn,
		CreatedAt: time.Now(),
	})
}

// write error messages
func (this *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if this.LogLevel < logger.Error {
		return
	}
	Write(Log{
		Message:   fmt.Sprintf(this.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...),
		Level:     LogLevelError,
		CreatedAt: time.Now(),
	})
}

// Trace write trace messages
func (this *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if this.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && this.LogLevel >= logger.Error:
		sql, rows := fc()
		if rows == -1 {
			Write(Log{
				Message:   fmt.Sprintf(this.traceErrStr+utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql),
				Level:     LogLevelError,
				CreatedAt: time.Now(),
			})
		} else {
			Write(Log{
				Message:   fmt.Sprintf(this.traceErrStr+utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql),
				Level:     LogLevelError,
				CreatedAt: time.Now(),
			})
		}
	case elapsed > this.SlowThreshold && this.SlowThreshold != 0 && this.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", this.SlowThreshold)
		if rows == -1 {
			Write(Log{
				Message:   fmt.Sprintf(this.traceWarnStr+utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql),
				Level:     LogLevelWarn,
				CreatedAt: time.Now(),
			})
		} else {
			Write(Log{
				Message:   fmt.Sprintf(this.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql),
				Level:     LogLevelWarn,
				CreatedAt: time.Now(),
			})
		}
	case this.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			Write(Log{
				Message:   fmt.Sprintf(this.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql),
				Level:     LogLevelInfo,
				CreatedAt: time.Now(),
			})
		} else {
			Write(Log{
				Message:   fmt.Sprintf(this.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql),
				Level:     LogLevelInfo,
				CreatedAt: time.Now(),
			})
		}
	}
}
