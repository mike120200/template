package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MySQLDB *gorm.DB

// InitMySQL 初始化 MySQL 数据库连接（使用 GORM）
func InitMySQL() error {
	host := viper.GetString("MySQL.host")
	port := viper.GetString("MySQL.port")
	user := viper.GetString("MySQL.user")
	password := viper.GetString("MySQL.password")
	database := viper.GetString("MySQL.database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("打开 MySQL 连接失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取 MySQL 底层连接失败: %w", err)
	}
	if maxOpenConns := viper.GetInt("MySQL.maxOpenConns"); maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns := viper.GetInt("MySQL.maxIdleConns"); maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL 连接测试失败: %w", err)
	}
	MySQLDB = db
	zap.L().Info("MySQL 数据库初始化成功", zap.String("host", host), zap.String("database", database))
	return nil
}

func GetMySql() *gorm.DB {
	return MySQLDB
}
