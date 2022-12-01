package db

import (
	"errors"
	"github.com/game/server/crypto"
	"github.com/game/server/database"
	"time"
	"github.com/game/server/config"
)

//用户表
type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"primary_key;column:name" form:"name" json:"name"` //名称
	Password  string    `json:"password"`
	Mail      string    `json:"mail"`
	PhoneNum  string    `json:"phone_num"`
	Company   string    `json:"company"`
	Type      int       `json:"type"` //用户类型(0:管理员， 1:普通用户)
	IsValid   bool      `json:"is_valid"`
}

const (
	SUPER_ADMIN = 0
)

func CheckUser(user *User) error {
	return database.Model(Db, &User{}).CheckExist("\"name\" = ?", user.Name)
}

func InsertUser(user *User) error {
	user.Password = crypto.PasswdEncryMD5(user.Password)
	return database.InsertAndGetDb(Db, user)
}

func GetUserByName(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("The user name is empty!")
	}
	var user User
	err := database.GetOne(Db, &user, "\"name\" = ?", name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByName(name string) error {
	if name == "" {
		return errors.New("The user name is empty")
	}

	err := database.DeleteOne(Db, &User{}, "name = ?", name)
	return err
}

func GetUsers() (users []*User, err error) {
	cond := "\"created_at\""
	err = database.Model(Db, &User{}).Order(cond, 1).Find(&users)
	return
}

func InitDbUser() error {
	var err error
	if err = Db.Migrate(&User{}); err != nil {
		return err
	}
	if err = Db.Migrate(&Token{}); err != nil {
		return err
	}
	userName := config.GetServerConfig().Config.AdminUser
	user := &User{Name: userName, Password: "88888888", Type: 0, IsValid: true}
	if CheckUser(user) == nil {
		err := InsertUser(user)
		logger.Errorf("Insert user root err : %#v\n", err)
	}
	var ErrorUserExist = errors.New("user already exists !")

	if err != ErrorUserExist && err != nil {
		return err
	}

	userName1 := "test"
	user1 := &User{Name: userName1, Password: "66666666", Type: 1, IsValid: true}
	if CheckUser(user1) == nil {
		err := InsertUser(user1)
		logger.Errorf("Insert user1 test err : %#v\n", err)
	}
	var ErrorUser1Exist = errors.New("user1 test already exists !")

	if err != ErrorUser1Exist && err != nil {
		return err
	}

	return nil
}
