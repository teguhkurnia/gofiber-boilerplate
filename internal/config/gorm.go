package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type logrusWriter struct {
	Logger *logrus.Logger
}

func (w *logrusWriter) Printf(format string, args ...any) {
	w.Logger.Tracef(format, args...)
}

func NewDatabase(config *viper.Viper, log *logrus.Logger, test bool) *gorm.DB {
	username := config.GetString("database.username")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	dbname := config.GetString("database.dbname")
	if test {
		dbname = config.GetString("database.test_dbname")
	}
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{
			Logger: log,
		}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %+v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifeTimeConnection) * time.Second)

	return db
}
