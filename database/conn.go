// Copyright 2023 @falcucci
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package database

import (
	"fmt"

	_ "log"

	"github.com/jinzhu/gorm"

	// postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB = SetupDB()
)

// Method to create a connection and returns the database built-in methods
func SetupDB() *gorm.DB {
	connStr := fmt.Sprintf(
		(`dbname=%s user=%s password=%s host=%s port=%s sslmode=disable
		application_name=%s`),
		Env.DatabaseName,
		Env.UserName,
		Env.Password,
		Env.Host,
		Env.Port,
		Env.ApplicationName,
	)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		logger.LogIt(
			logger.LevelError,
			"database",
			"SetupDB",
			"sql.Open",
			err.Error(),
		)
	}

	db.LogMode(Env.ShowSql)
	err = db.DB().Ping()
	if err != nil {
		logger.LogIt(
			logger.LevelError,
			"database",
			"SetupDB",
			"db.DB().Ping",
			err.Error(),
		)
	}

	return db
}
