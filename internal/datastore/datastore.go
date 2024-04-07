// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package datastore

import (
	"fmt"
	"runtime"
	"ss-assistant/internal/config"
	"strconv"
	"time"

	gorm_sql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SyntSugar/ss-infra-go/datastore/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type DataStore struct {
	DB    *gorm.DB
	Cache *redis.Client
}

var _dataStore *DataStore

// transform2MysqlConf transform config.Config to mysql.Config
func transform2MysqlConf(cfg *config.Config) (mysqlConf mysql.Config) {
	address := fmt.Sprintf("%s:%s", cfg.DbConfig.Host, strconv.Itoa(cfg.DbConfig.Port))
	return mysql.Config{
		Addr:         address,
		DBName:       cfg.DbConfig.DbName,
		User:         cfg.DbConfig.User,
		Password:     cfg.DbConfig.Pwd,
		MaxOpenConns: cfg.DbConfig.MaxConn,
	}
}

// transform2RedisOptions transform config.Config to redis.Options
func transform2RedisOptions(cfg *config.Config) (redisOption *redis.Options) {
	address := fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, strconv.Itoa(cfg.RedisConfig.Port))
	return &redis.Options{
		Addr:     address,
		Password: cfg.RedisConfig.Pwd,

		DialTimeout:  1200 * time.Millisecond,
		ReadTimeout:  1500 * time.Millisecond,
		WriteTimeout: time.Second,
		MinIdleConns: runtime.NumCPU(),
	}
}

// Init init datastore
func Init(cfg *config.Config) {
	mysqlConfig := transform2MysqlConf(cfg)
	mysqlClient, err := mysql.NewClient(&mysqlConfig)
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(gorm_sql.New(gorm_sql.Config{
		Conn: mysqlClient,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	redisOption := transform2RedisOptions(cfg)
	redisClient := redis.NewClient(redisOption)

	_dataStore = &DataStore{
		DB:    gormDB,
		Cache: redisClient,
	}
}

// Get get datastore
func Get() *DataStore {
	return _dataStore
}

// Close close datastore
func Close() {
	_dataStore.Cache.Close()
}
