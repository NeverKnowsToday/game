package db

import (
	"github.com/game/server/config"
	"github.com/game/server/database"
	"github.com/game/server/logger/logging"
)

var (
	Db     *database.MyDb
	logger = logging.GetLogger("db", logging.DEFAULT_LEVEL)
)

func InitDb() error {
	var err error
	dbInfo := config.GetServerConfig().Db
	host := dbInfo.Address
	port := dbInfo.Port
	dbuser := dbInfo.User
	dbname := "baoxian"
	password := dbInfo.Password
	Type := dbInfo.Type
	Db, err = database.OpenDatabase(port, host, dbuser, dbname, password, Type)
	if err != nil {
		return err
	}
	err = Migrate(Db)
	if err != nil {
		return err
	}

	InitDbUser()

	return nil
}

// Migrate is auto create tables
func Migrate(Db *database.MyDb) error {
	logger.Info("Mysql Migrate database structs.")
	if err := Db.Migrate(&User{}); err != nil {
		return err
	}

	if err := Db.Migrate(&Token{}); err != nil {
		return err
	}
	if err := Db.Migrate(&Invoke{}); err != nil {
		return err
	}
	return nil
}
