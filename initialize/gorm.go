package initialize

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"server/global"
)

func InitGorm() *gorm.DB {
	mysqlCfg := global.Config.Mysql
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()),
	})
	if err != nil {
		global.Log.Error("Failed to connect to mysql", zap.Error(err))
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)

	return db
}
