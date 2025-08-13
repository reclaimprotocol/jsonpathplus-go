package jsonpathplus

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger interface for pluggable logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// DefaultLogger implements a simple logger
type DefaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stderr, "[JSONPath] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *DefaultLogger) log(level LogLevel, levelStr, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = " |"
		for _, field := range fields {
			fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	l.logger.Printf("[%s] %s%s", levelStr, msg, fieldStr)
}

func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	l.log(LogLevelDebug, "DEBUG", msg, fields...)
}

func (l *DefaultLogger) Info(msg string, fields ...Field) {
	l.log(LogLevelInfo, "INFO", msg, fields...)
}

func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	l.log(LogLevelWarn, "WARN", msg, fields...)
}

func (l *DefaultLogger) Error(msg string, fields ...Field) {
	l.log(LogLevelError, "ERROR", msg, fields...)
}

// NoOpLogger implements a logger that does nothing
type NoOpLogger struct{}

func (l *NoOpLogger) Debug(msg string, fields ...Field) {}
func (l *NoOpLogger) Info(msg string, fields ...Field)  {}
func (l *NoOpLogger) Warn(msg string, fields ...Field)  {}
func (l *NoOpLogger) Error(_ string, _ ...Field)        {}

// Metrics collects performance metrics
type Metrics struct {
	QueriesExecuted      int64
	TotalExecutionTime   time.Duration
	AverageExecutionTime time.Duration
	ErrorCount           int64
	MemoryUsage          int64
}

// MetricsCollector collects and tracks metrics
type MetricsCollector struct {
	metrics Metrics
	enabled bool
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(enabled bool) *MetricsCollector {
	return &MetricsCollector{
		enabled: enabled,
	}
}

// RecordQuery records a query execution
func (m *MetricsCollector) RecordQuery(duration time.Duration, err error) {
	if !m.enabled {
		return
	}

	m.metrics.QueriesExecuted++
	m.metrics.TotalExecutionTime += duration
	m.metrics.AverageExecutionTime = time.Duration(int64(m.metrics.TotalExecutionTime) / m.metrics.QueriesExecuted)

	if err != nil {
		m.metrics.ErrorCount++
	}
}

// UpdateMemoryUsage updates the current memory usage
func (m *MetricsCollector) UpdateMemoryUsage(usage int64) {
	if !m.enabled {
		return
	}
	m.metrics.MemoryUsage = usage
}

// GetMetrics returns a copy of the current metrics
func (m *MetricsCollector) GetMetrics() Metrics {
	return m.metrics
}

// Reset resets all metrics
func (m *MetricsCollector) Reset() {
	m.metrics = Metrics{}
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Error(key string, err error) Field {
	return Field{Key: key, Value: err}
}
