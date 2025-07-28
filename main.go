package main

import (
	"fmt"
	"log"

	"high-concurrency-api/dao"
	"high-concurrency-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化MySQL连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %s", err)
	}

	// 设置数据库连接池
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))

	// 初始化Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	// 初始化DAO和Handler
	dataDAO := dao.NewDataDAO(db, rdb)
	dataHandler := handlers.NewDataHandler(dataDAO)

	// 设置Gin模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建路由
	r := gin.Default()

	// 注册路由
	r.POST("/api/data", dataHandler.Create)
	r.PUT("/api/data/:id", dataHandler.Update)
	r.DELETE("/api/data/:id", dataHandler.Delete)
	r.GET("/api/data/:id", dataHandler.Get)

	// 启动服务器
	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
} 