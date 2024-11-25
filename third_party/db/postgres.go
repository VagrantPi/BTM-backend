package db

import (
	"BTM-backend/configs"
	"BTM-backend/pkg/error_code"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var dbDialector gorm.Dialector

// ConnectToDatabase 這個只會給 di 使用，在一開始 di 就會連線了，如果拿去給地方使用有可能會錯誤
func ConnectToDatabase() *gorm.DB {
	if db == nil {
		var mu sync.Mutex
		mu.Lock()
		if db == nil {
			StartDatabase()
		}
		mu.Unlock()
	}
	return db
}

func StartDatabase() {
	var err error

	dbDialector = newPgDialector()
	db, err = setupGORM(
		dbDialector, newGormConfig(newGormLog()),
	) // like setupBun or setupMysql, other libraries.
	if err != nil {
		log.Fatalf("database 啟用錯誤...: %v", err)
	}

	log.Println("database 啟用成功")
}

// setupGORM 連線後會做健康檢查，並返回連線池實例
// - 閒置 10 秒自動斷線
// - 最多 100 連線數
func setupGORM(dialector gorm.Dialector, gormConfig *gorm.Config) (gormDb *gorm.DB, err error) {
	// 連線啟動，套件有做健康檢查
	gormDb, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "setupGORM: gorm.Open").WithCause(err)
		return nil, err
	}

	if isDebug() {
		gormDb = gormDb.Debug()
	}

	d, err := gormDb.DB()
	if err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "get GORM connect setting failed").WithCause(err)
		return nil, err
	}

	d.SetMaxOpenConns(configs.C.Db.MaxOpenConns)
	d.SetMaxIdleConns(configs.C.Db.MaxIdleConns)
	d.SetConnMaxLifetime(time.Duration(configs.C.Db.ConnMaxLifetimeSec) * time.Second)

	return gormDb, nil
}

func isDebug() bool {
	return configs.C.Db.Debug
}

func newGormConfig(logSetting logger.Interface) *gorm.Config {
	return &gorm.Config{Logger: logSetting}
}

func newGormLog() logger.Interface {
	_, isInCloudRun := os.LookupEnv("K_SERVICE")

	switch {
	case !isDebug():
		return logger.Discard
	case isInCloudRun:
		return &gormLogger{logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)}
	default:
		return &gormLogger{logger.Default}
	}
}

func newPgDialector() gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		// "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True",
		configs.C.Db.Host,
		configs.C.Db.Username,
		configs.C.Db.Password,
		configs.C.Db.Database,
		configs.C.Db.Port,
		configs.C.Db.SslMode, // https://www.postgresql.org/docs/current/libpq-ssl.html#:~:text=Section%C2%A019.9.5.-,34.19.3.%C2%A0Protection%20Provided%20in%20Different%20Modes,-The%20different%20values
	)

	return postgres.New(
		postgres.Config{
			DSN: dsn,
			// PreferSimpleProtocol: true, // disables implicit prepared statement usage
		},
	)
}

type gormLogger struct {
	log logger.Interface
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{l.log.LogMode(level)}
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	traceID := getTraceIDFromContext(ctx)
	var traceMsg string
	if traceID != "" {
		traceMsg = fmt.Sprintf("[trace_id: %s] ", traceID)
	}
	l.log.Info(ctx, fmt.Sprintf("%s%s", traceMsg, msg), data...)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	traceID := getTraceIDFromContext(ctx)
	var traceMsg string
	if traceID != "" {
		traceMsg = fmt.Sprintf("[trace_id: %s] ", traceID)
	}
	l.log.Warn(ctx, fmt.Sprintf("%s%s", traceMsg, msg), data...)
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	traceID := getTraceIDFromContext(ctx)
	var traceMsg string
	if traceID != "" {
		traceMsg = fmt.Sprintf("[trace_id: %s] ", traceID)
	}
	l.log.Error(ctx, fmt.Sprintf("%s%s", traceMsg, msg), data...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	traceID := getTraceIDFromContext(ctx)
	var traceMsg string
	if traceID != "" {
		traceMsg = fmt.Sprintf("[trace_id: %s] ", traceID)
	}

	fcWithTraceMsg := func() (string, int64) {
		sql, rows := fc()
		return traceMsg + sql, rows
	}

	l.log.Trace(ctx, begin, fcWithTraceMsg, err)
}

func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func ProvideDb() *gorm.DB {
	return ConnectToDatabase()
}
