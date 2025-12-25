package db

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PGDB *gorm.DB

// InitPostgreSQL 初始化 PostgreSQL 数据库连接（使用 GORM）

func InitPostgreSQL() error {
	host := viper.GetString("PostgresDB.host")
	port := viper.GetString("PostgresDB.port")
	user := viper.GetString("PostgresDB.user")
	password := viper.GetString("PostgresDB.password")
	database := viper.GetString("PostgresDB.database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("打开 PostgreSQL 连接失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取 PostgreSQL 底层连接失败: %w", err)
	}
	if maxOpenConns := viper.GetInt("PostgresDB.maxOpenConns"); maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns := viper.GetInt("PostgresDB.maxIdleConns"); maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("PostgreSQL 连接测试失败: %w", err)
	}
	PGDB = db
	zap.L().Info("PostgreSQL 数据库初始化成功", zap.String("host", host), zap.String("database", database))
	return nil
}

func GetPg() *gorm.DB {
	return PGDB

}
