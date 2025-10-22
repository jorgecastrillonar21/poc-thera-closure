package monitoring

import (
	"time"

	"gorm.io/gorm"
)

// DatabaseMonitor wraps GORM to provide monitoring capabilities
type DatabaseMonitor struct {
	db      *gorm.DB
	metrics *Metrics
}

// NewDatabaseMonitor creates a new database monitor
func NewDatabaseMonitor(db *gorm.DB, metrics *Metrics) *DatabaseMonitor {
	monitor := &DatabaseMonitor{
		db:      db,
		metrics: metrics,
	}

	// Register GORM plugin for automatic monitoring
	db.Use(&gormMetricsPlugin{metrics: metrics})

	return monitor
}

// GetDB returns the underlying database connection
func (dm *DatabaseMonitor) GetDB() *gorm.DB {
	return dm.db
}

// UpdateConnectionMetrics updates database connection pool metrics
func (dm *DatabaseMonitor) UpdateConnectionMetrics() {
	if sqlDB, err := dm.db.DB(); err == nil {
		stats := sqlDB.Stats()
		dm.metrics.DatabaseConnectionsActive.Set(float64(stats.InUse))
		dm.metrics.DatabaseConnectionsIdle.Set(float64(stats.Idle))
	}
}

// gormMetricsPlugin implements the GORM plugin interface for metrics collection
type gormMetricsPlugin struct {
	metrics *Metrics
}

func (p *gormMetricsPlugin) Name() string {
	return "metrics"
}

func (p *gormMetricsPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for different operations
	db.Callback().Create().Before("gorm:create").Register("metrics:before_create", p.beforeCreate)
	db.Callback().Create().After("gorm:create").Register("metrics:after_create", p.afterCreate)

	db.Callback().Query().Before("gorm:query").Register("metrics:before_query", p.beforeQuery)
	db.Callback().Query().After("gorm:query").Register("metrics:after_query", p.afterQuery)

	db.Callback().Update().Before("gorm:update").Register("metrics:before_update", p.beforeUpdate)
	db.Callback().Update().After("gorm:update").Register("metrics:after_update", p.afterUpdate)

	db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", p.beforeDelete)
	db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", p.afterDelete)

	return nil
}

func (p *gormMetricsPlugin) beforeCreate(db *gorm.DB) {
	db.Set("metrics:start_time", time.Now())
	db.Set("metrics:operation", "create")
}

func (p *gormMetricsPlugin) afterCreate(db *gorm.DB) {
	p.recordMetrics(db, "create")
}

func (p *gormMetricsPlugin) beforeQuery(db *gorm.DB) {
	db.Set("metrics:start_time", time.Now())
	db.Set("metrics:operation", "query")
}

func (p *gormMetricsPlugin) afterQuery(db *gorm.DB) {
	p.recordMetrics(db, "query")
}

func (p *gormMetricsPlugin) beforeUpdate(db *gorm.DB) {
	db.Set("metrics:start_time", time.Now())
	db.Set("metrics:operation", "update")
}

func (p *gormMetricsPlugin) afterUpdate(db *gorm.DB) {
	p.recordMetrics(db, "update")
}

func (p *gormMetricsPlugin) beforeDelete(db *gorm.DB) {
	db.Set("metrics:start_time", time.Now())
	db.Set("metrics:operation", "delete")
}

func (p *gormMetricsPlugin) afterDelete(db *gorm.DB) {
	p.recordMetrics(db, "delete")
}

func (p *gormMetricsPlugin) recordMetrics(db *gorm.DB, operation string) {
	startTime, exists := db.Get("metrics:start_time")
	if !exists {
		return
	}

	start, ok := startTime.(time.Time)
	if !ok {
		return
	}

	duration := time.Since(start)
	table := db.Statement.Table
	if table == "" {
		table = "unknown"
	}

	status := "success"
	if db.Error != nil {
		status = "error"
		p.metrics.DatabaseErrors.WithLabelValues(operation, db.Error.Error()).Inc()
	}

	p.metrics.RecordDatabaseOperation(operation, table, status, duration)
}

// RedisMonitor provides monitoring for Redis operations
type RedisMonitor struct {
	metrics *Metrics
}

// NewRedisMonitor creates a new Redis monitor
func NewRedisMonitor(metrics *Metrics) *RedisMonitor {
	return &RedisMonitor{
		metrics: metrics,
	}
}

// RecordRedisCommand records metrics for a Redis command
func (rm *RedisMonitor) RecordRedisCommand(command string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
		rm.metrics.RedisErrors.WithLabelValues(command, err.Error()).Inc()
	}

	rm.metrics.RedisCommandsTotal.WithLabelValues(command, status).Inc()
	rm.metrics.RedisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// StripeMonitor provides monitoring for Stripe API calls
type StripeMonitor struct {
	metrics *Metrics
}

// NewStripeMonitor creates a new Stripe monitor
func NewStripeMonitor(metrics *Metrics) *StripeMonitor {
	return &StripeMonitor{
		metrics: metrics,
	}
}

// RecordAPICall records metrics for a Stripe API call
func (sm *StripeMonitor) RecordAPICall(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	sm.metrics.RecordStripeAPICall(operation, status, duration)
}
