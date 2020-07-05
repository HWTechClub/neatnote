// Neat Note. A notes sharing platform for university students.
// Copyright (C) 2020 Humaid AlQassimi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package models

import (
	"fmt"
	"git.sr.ht/~humaid/neatnote/modules/settings"
	_ "github.com/go-sql-driver/mysql" // MySQL driver support
	_ "github.com/mattn/go-sqlite3"    // SQLite driver support
	"log"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables,
		new(User),
		new(Course),
		new(Post),
		new(Comment),
	)
}

// SetupEngine sets up an XORM engine according to the database configuration
// and syncs the schema.
func SetupEngine() *xorm.Engine {
	var err error
	dbConf := &settings.Config.DBConfig

	address := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Name)

	switch dbConf.Type {
	case settings.MySQL:
		engine, err = xorm.NewEngine("mysql", address)
	case settings.SQLite:
		engine, err = xorm.NewEngine("sqlite3", dbConf.Path)
	}

	if err != nil {
		log.Fatal("Unable to connect/load the database! ", err)
	}

	engine.SetMapper(core.GonicMapper{}) // So ID becomes 'id' instead of 'i_d'
	err = engine.Sync(tables...)         // Sync the schema of tables

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 2000)
	engine.SetDefaultCacher(cacher)

	if err != nil {
		log.Fatal("Unable to sync schema! ", err)
	}

	return engine
}
