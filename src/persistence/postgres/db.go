package postgres

import (
	"chatbot_be_go/src/persistence/config"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
)

type IDB interface {
	GetDb() *gorm.DB
	GetSqlDb() *sql.DB
}

type appDb struct {
	DB    *gorm.DB
	SqlDB *sql.DB
}

func (a *appDb) GetDb() *gorm.DB {
	return a.DB
}

func (a *appDb) GetSqlDb() *sql.DB {
	return a.SqlDB
}

func New(conf config.SqlDbConf, logger *logrus.Logger) IDB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable client_encoding=utf8",
		conf.Host,
		conf.Username,
		conf.Password,
		conf.Name,
		conf.Port,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: log.Default.LogMode(log.Warn),
	})
	if err != nil {
		panic("Failed to connect to database!")
	}

	if conf.IsLogged {
		db.Config.Logger = log.Default.LogMode(log.Info)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("database err: %s", err)
	}

	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifeTimeConnSeconds) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(conf.MaxIdleTimeConnSeconds) * time.Second)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)

	return &appDb{
		DB:    db,
		SqlDB: sqlDB,
	}
}
