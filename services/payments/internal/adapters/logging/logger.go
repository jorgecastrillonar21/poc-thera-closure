package logging

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Fields is a collection of Field items
type Fields []Field

// NewLogger creates a new structured logger instance
func NewLogger(level string) *Logger {
	logger := logrus.New()

	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	logger.SetOutput(os.Stdout)

	return &Logger{logger}
}

// WithContext adds request context information to logs
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.Logger.WithContext(ctx)

	// Add request ID if available
	if reqID := ctx.Value("request_id"); reqID != nil {
		entry = entry.WithField("request_id", reqID)
	}

	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	return entry
}

// WithFields adds structured fields to the log entry
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return l.Logger.WithFields(logrusFields)
}

// WithError adds error information to the log entry
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// Business operation logging methods
func (l *Logger) LogBusinessOperation(ctx context.Context, operation string, entityID string, success bool, duration time.Duration, fields Fields) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"operation":   operation,
		"entity_id":   entityID,
		"success":     success,
		"duration_ms": duration.Milliseconds(),
		"service":     "payments",
		"category":    "business",
	})

	// Add custom fields
	for _, field := range fields {
		entry = entry.WithField(field.Key, field.Value)
	}

	if success {
		entry.Info("Business operation completed successfully")
	} else {
		entry.Error("Business operation failed")
	}
}

// Database operation logging
func (l *Logger) LogDatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, err error) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
		"service":     "payments",
		"category":    "database",
	})

	if err != nil {
		entry.WithError(err).Error("Database operation failed")
	} else {
		entry.Debug("Database operation completed")
	}
}

// External API logging
func (l *Logger) LogExternalAPICall(ctx context.Context, service string, endpoint string, method string, statusCode int, duration time.Duration, err error) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"external_service": service,
		"endpoint":         endpoint,
		"method":           method,
		"status_code":      statusCode,
		"duration_ms":      duration.Milliseconds(),
		"service":          "payments",
		"category":         "external_api",
	})

	if err != nil {
		entry.WithError(err).Error("External API call failed")
	} else if statusCode >= 400 {
		entry.Warn("External API returned error status")
	} else {
		entry.Debug("External API call completed")
	}
}

// Security event logging
func (l *Logger) LogSecurityEvent(ctx context.Context, eventType string, severity string, details Fields) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"event_type": eventType,
		"severity":   severity,
		"service":    "payments",
		"category":   "security",
	})

	// Add custom details
	for _, field := range details {
		entry = entry.WithField(field.Key, field.Value)
	}

	switch severity {
	case "critical":
		entry.Error("Critical security event detected")
	case "high":
		entry.Warn("High priority security event")
	case "medium":
		entry.Info("Medium priority security event")
	default:
		entry.Debug("Security event logged")
	}
}

// Performance metrics logging
func (l *Logger) LogPerformanceMetric(ctx context.Context, metricName string, value float64, unit string, fields Fields) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"metric_name": metricName,
		"value":       value,
		"unit":        unit,
		"service":     "payments",
		"category":    "performance",
	})

	// Add custom fields
	for _, field := range fields {
		entry = entry.WithField(field.Key, field.Value)
	}

	entry.Info("Performance metric recorded")
}

// Audit logging for compliance
func (l *Logger) LogAuditEvent(ctx context.Context, action string, resource string, outcome string, fields Fields) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"action":   action,
		"resource": resource,
		"outcome":  outcome,
		"service":  "payments",
		"category": "audit",
	})

	// Add custom fields
	for _, field := range fields {
		entry = entry.WithField(field.Key, field.Value)
	}

	entry.Info("Audit event logged")
}
