package db

import (
	"github.com/game/server/database"
	"time"
)
//操作表
type ProductExcelParser struct {
	Id        string    `gorm:"primary_key;column:id" form:"id" json:"id"`
	ProductId   string  `json:"product_id"`
	ExcelName  string   `json:"excel_name"` //名称
	ExcelPath  string   `json:"excel_path"`
	SetInfo	   string   `json:"set_info"`
	GetInfo    string   `json:"get_info"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

//func CheckInvoke(name string) error {
//	return database.Model(Db, &Invoke{}).CheckExist("\"name\" = ?", name)
//}

//func InsertInvoke(invoke *Invoke) error {
//	return database.InsertAndGetDb(Db, invoke)
//}

func GetProductExcelParsers() ([]*ProductExcelParser, error) {
	var productExcelParser []*ProductExcelParser
	err := database.Model(Db, &ProductExcelParser{}).Find(&productExcelParser)
	if err != nil {
		return nil, err
	}
	return productExcelParser, nil
}

func GetProductExcelParserByID(id int) (*ProductExcelParser, error) {
	productExcelParser := new(ProductExcelParser)
	err := database.Model(Db, &ProductExcelParser{}).Find(&productExcelParser)
	if err != nil {
		return nil, err
	}
	return productExcelParser, nil
}

//func GetInvokeByName(name string) (*Invoke, error) {
//	if name == "" {
//		return nil, errors.New("The user name is empty!")
//	}
//	var invoke Invoke
//	err := database.GetOne(Db, &invoke, "\"name\" = ?", name)
//	if err != nil {
//		return nil, err
//	}
//	return &invoke, nil
//}
//
//func GetInvokesByRoom(room int) ([]*Invoke, error) {
//	var invoke []*Invoke
//	err := database.Model(Db, &Invoke{}).FilterBy("room", room).Find(&invoke)
//	if err != nil {
//		return nil, err
//	}
//	return invoke, nil
//}
//
//// UpdateChannel 更新状态
//func UpdateInvokePosByname(name string, pos int) error {
//	return database.Model(Db, &Invoke{}).Where("\"name\" = ?", name).Update("current_pos", pos)
//}
//// UpdateChannel 更新状态
//func UpdateInvokeWinnerByname(name string, state bool) error {
//	return database.Model(Db, &Invoke{}).Where("\"name\" = ?", name).Update("is_win", state)
//}
//
//
//func DeleteInvokeByName(name string) error {
//	if name == "" {
//		return errors.New("The user name is empty")
//	}
//
//	err := database.DeleteOne(Db, &Invoke{}, "name = ?", name)
//	return err
//}
