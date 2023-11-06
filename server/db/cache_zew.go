package db

import (
	"github.com/game/server/database"
	"time"
)


//操作表
type CacheZew struct {
	Id        string    `gorm:"primary_key;column:id" form:"id" json:"id"`
	ProductId   string  `json:"product_id"`
	CompanyID    string    `json:"company_id"`
	Age        string   `json:"age"`
	Sex        string   `json:"sex"`
	Total      string   `json:"total"`
	Period     string   `json:"period"`
	Hash       string   `json:"hash"`
    ExcelResult string  `json:"excel_result"`
	Deleted    int      `json:"deleted"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

//func CheckInvoke(name string) error {
//	return database.Model(Db, &Invoke{}).CheckExist("\"name\" = ?", name)
//}

func InsertCacheZew(cacheZew *CacheZew) error {
	return database.InsertDb(Db, cacheZew)
}

//func GetProductExcelParsers() ([]*ProductExcelParser, error) {
//	var productExcelParser []*ProductExcelParser
//	err := database.Model(Db, &ProductExcelParser{}).Find(&productExcelParser)
//	if err != nil {
//		return nil, err
//	}
//	return productExcelParser, nil
//}

func GetCacheZewByHash(hash string) (*CacheZew, error) {
	cacheZew := new(CacheZew)
	err := database.Model(Db, &CacheZew{}).FilterBy("hash", hash).Find(&cacheZew)
	if err != nil {
		return nil, err
	}
	return cacheZew, nil
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
