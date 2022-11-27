package db

import (
	"errors"
	"github.com/game/server/database"
	"time"
)

//操作表
type Invoke struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"primary_key;column:name" form:"name" json:"name"` //名称
    CurrentPos int      `json:"current_pos"`
	IsWin     bool      `json:"is_win"`
	Room      int       `json:"room"`
}

func CheckInvoke(invoke *Invoke) error {
	return database.Model(Db, &Invoke{}).CheckExist("\"name\" = ?", invoke.Name)
}

func InsertInvoke(invoke *Invoke) error {
	return database.InsertAndGetDb(Db, invoke)
}

func GetInvokeByName(name string) (*Invoke, error) {
	if name == "" {
		return nil, errors.New("The user name is empty!")
	}
	var invoke Invoke
	err := database.GetOne(Db, &invoke, "\"name\" = ?", name)
	if err != nil {
		return nil, err
	}
	return &invoke, nil
}

// GetChannelByUUID 通过netuuid, uuid查询通道信息
func GetInvokesByRoom(room int) ([]*Invoke, error) {
	var invoke []*Invoke
	err := database.Model(Db, &Invoke{}).FilterBy("room", room).Find(&invoke)
	if err != nil {
		return nil, err
	}
	return invoke, nil
}

func DeleteInvokeByName(name string) error {
	if name == "" {
		return errors.New("The user name is empty")
	}

	err := database.DeleteOne(Db, &Invoke{}, "name = ?", name)
	return err
}
