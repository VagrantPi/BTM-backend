package db

import (
	"BTM-backend/configs"
	"BTM-backend/pkg/error_code"
	pkgLogger "BTM-backend/pkg/logger"
	"BTM-backend/pkg/tools"
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var dbDialector gorm.Dialector

func ConnectToDatabase() (*gorm.DB, error) {
	if db == nil {
		var mu sync.Mutex
		mu.Lock()
		if db == nil {
			StartDatabase()
		}
		mu.Unlock()
	}
	return db, nil
}

func StartDatabase() error {
	var err error
	dbDialector = newPgDialector()
	db, err = setupGORM(
		dbDialector, newGormConfig(newGormLog()),
	)
	if err != nil {
		return err
	}

	log.Println("database 啟用成功")
	return nil
}

func ConnectToMockDatabase() (*gorm.DB, error) {
	if db == nil {
		var mu sync.Mutex
		mu.Lock()
		if db == nil {
			GenMockDB()
		}
		mu.Unlock()
	}
	return db, nil
}

func GenMockDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("mock.db"), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return err
	}

	log.Println("mock database 啟用成功")
	return nil
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
	switch {
	case !isDebug():
		return logger.Discard
	default:
		return sLogger{
			slog.Default(),
			pkgLogger.Zap(),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
			}}
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

type sLogger struct {
	log       *slog.Logger
	pkgLogger *pkgLogger.Logger
	logger.Config
}

func (l sLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func (l sLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, data...)
	l.log.Info(m)
}

func (l sLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, data...)
	l.log.Warn(m)
}

func (l sLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, data...)
	l.log.Error(m)
}

func (l sLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	var logMsg string
	traceID := getTraceIDFromContext(ctx)

	_logger, pkgLoggerOk := ctx.Value("log").(*pkgLogger.Logger)
	if pkgLoggerOk {
		l.pkgLogger = _logger
	}

	caller := tools.GetCallerInfo(4)
	logMsg = logMsg + caller

	elapsed := time.Since(begin)
	sql, rows := fc()

	// If there was an error, add it to the log message, or if the query was slow, add slow query log
	errOrSlowMsg := ""
	if err != nil {
		errOrSlowMsg = fmt.Sprintf("ERROR: %v", err)
	} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		errOrSlowMsg = fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)

	}

	num := int64(0)
	if rows != -1 {
		num = rows
	}
	if pkgLoggerOk {
		l.pkgLogger.Info(logMsg,
			zap.String("err or slowMsg", errOrSlowMsg),
			zap.String("times and rows", fmt.Sprintf("[%.3fms][rows:%d]", float64(elapsed.Microseconds())/1000.0, num)),
			zap.String("caller", caller),
			zap.String("sql", sql),
			zap.String("traceID", traceID),
			zap.Any("trace_id", traceID),
		)
	} else {
		l.log.Info(fmt.Sprintf("%s[%.3fms] [rows:%d] %s", logMsg, float64(elapsed.Microseconds())/1000.0, num, sql))
	}
}

func ProvideDatabase(isMock bool) (*gorm.DB, error) {
	if isMock {
		return ConnectToMockDatabase()
	}
	return ConnectToDatabase()
}
