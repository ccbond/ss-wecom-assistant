// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package datastore

import (
	"fmt"
	"ss-wecom-assistant/internal/config"
	"ss-wecom-assistant/internal/model"
	"strconv"

	gorm_sql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SyntSugar/ss-infra-go/datastore/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type DataStore struct {
	DB *gorm.DB
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

	// auto migrate
	gormDB.AutoMigrate(&model.SessionInfo{})
	gormDB.AutoMigrate(&model.User{})

	_dataStore = &DataStore{
		DB: gormDB,
	}
}

// Get get datastore
func Get() *DataStore {
	return _dataStore
}
