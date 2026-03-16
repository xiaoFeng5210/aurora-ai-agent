package database

import (
	"aurora-agent/utils"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func DBConnect() (*gorm.DB, error) {
	zapLogger := utils.InitZap("log/zap")
	dbOnce.Do(func() {
		// 从环境变量获取数据库配置
		conf := utils.InitViper("conf", "db", "yaml")
		host := conf.GetString("postgres.host")
		user := conf.GetString("postgres.user")
		password := conf.GetString("postgres.password")
		dbname := conf.GetString("postgres.dbname")
		port := conf.GetString("postgres.port")

		zapLogger.Info("Database configuration loaded successfully.", zap.String("host", host), zap.String("user", user), zap.String("password", password), zap.String("dbname", dbname), zap.String("port", port))

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host, user, password, dbname, port, "disable", "Asia/Shanghai")

		logFile, _ := os.OpenFile("log/db.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		dbCopy, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{ //覆盖默认的NamingStrategy来更改命名约定
				SingularTable: true, //表名映射时不加复数，仅是驼峰-->蛇形
			},
			Logger: logger.New(
				log.New(logFile, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             500 * time.Millisecond, //耗时超过此值认定为慢查询
					LogLevel:                  logger.Info,            // LogLevel的最低阈值，Silent为不输出日志
					IgnoreRecordNotFoundError: true,                   // 忽略RecordNotFound这种错误日志
					ParameterizedQueries:      false,                  // true代表SQL日志里不包含参数
					Colorful:                  false,                  // 禁用颜色
				},
			),
		})
		if err != nil {
			zapLogger.Error("Failed to connect to database.", zap.Error(err))
			panic(err)
		}

		db = dbCopy

		zapLogger.Info("Database connected successfully.")
	})

	sqlDB, _ := db.DB()
	//池子里空闲连接的数量上限（超出此上限就把相应的连接关闭掉）
	sqlDB.SetMaxIdleConns(10)
	//最多开这么多连接
	sqlDB.SetMaxOpenConns(100)
	//一个连接最多可使用这么长时间，超时后连接会自动关闭（因为数据库本身可能也对NoActive连接设置了超时时间，我们的应对办法：定期ping，或者SetConnMaxLifetime）
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
