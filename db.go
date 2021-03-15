package pkg

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type DB struct {
	driver *gorm.DB
}

// NewDB return a new DB instance
func NewDB(debug bool, dsn string) (*DB, error) {
	driver, err := initGorm(debug, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{driver: driver}, nil
}

// NewDBWithMockForTest return a new DB instance
// JUST FOR TESTS
// JUST FOR TESTS
// JUST FOR TESTS
func NewDBWithMockForTest(debug bool, conn *sql.DB) (*DB, error) {
	driver, err := initGormWithConn(debug, conn)
	if err != nil {
		return nil, err
	}
	return &DB{driver: driver}, nil
}

func initGormLog(debug bool) logger.Interface {
	var config = logger.Config{
		SlowThreshold: time.Second * 1,
		Colorful:      false,
		LogLevel:      logger.Info,
	}
	if !debug {
		config.LogLevel = logger.Warn
	}
	return logger.New(log.New(WithSampleLog(), "", 0), config)
}

// initGormWithConn init gorm DB instance with sql.DB (for tests)
func initGormWithConn(debug bool, conn *sql.DB) (*gorm.DB, error) {
	return gorm.Open(
		mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}),
		&gorm.Config{Logger: initGormLog(debug)},
	)
}

// InitDB init gorm DB instance with dsn
func initGorm(debug bool, dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: initGormLog(debug)})
}

func (db *DB) GetDriver() *gorm.DB {
	return db.driver
}

// QueryOne
func (db *DB) QueryOne(dest interface{}, sql string, args ...interface{}) error {
	return db.driver.Raw(sql, args).Scan(dest).Error
}

// QueryMore 查询多个结果，返回的结果需要配合 ScanRows 使用
// example:
// rows, _ := db.QueryMore("SELECT * FROM table LIMIT 10")
// for rows.Next(){
//    _ =  db.ScanRows(rows, &dest)
// }
func (db *DB) QueryMore(sql string, args ...interface{}) (*sql.Rows, error) {
	return db.driver.Raw(sql, args...).Rows()
}

// ScanRows 扫描 rows 到 dest
func (db *DB) ScanRows(rows *sql.Rows, dest interface{}) error {
	return db.driver.ScanRows(rows, dest)
}
