package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/game/server/logger/logging"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

type NumCounter struct {
	ID    string `gorm:"primary_key"` //后端识别名
	Index uint   //该序列
}

var logger = logging.GetLogger("database", logging.DEFAULT_LEVEL)

const LEN_UUID = 8
const DEFAULT_ID uint = 0

var index_lock sync.Mutex

//check integer and string or pointer
func IsNull(val interface{}) bool {
	if val == nil {
		return true
	}
	value := reflect.ValueOf(val)
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	valType := value.Type()
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Bool:
		return value.Bool() == false
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		return value.IsNil()
	case reflect.Array, reflect.String, reflect.Slice, reflect.Map:
		return value.Len() == 0
	default:
		return val == reflect.New(valType).Elem().Interface()
	}
}

type MyDb struct {
	*gorm.DB
	whereConditions string
}

func Model(db *MyDb, data interface{}) *MyDb {
	if data == nil {
		return &MyDb{DB: db.DB}
	} else {
		return &MyDb{DB: db.DB.Model(data)}
	}
}

func (m *MyDb) Where(query string, args ...interface{}) *MyDb {
	condition := m.whereConditions
	if condition != "" {
		condition += " and "
	}
	if m.Dialect().GetName() != "oci8" {
		query = strings.Replace(query, "\"", " ", -1)
	}
	condition += fmt.Sprintf(strings.Replace(query, "?", "%v", -1), args...)
	m.DB = m.DB.Where(query, args...)
	m.whereConditions = condition
	return m
}

func (m *MyDb) WhereClear() *MyDb {
	m.whereConditions = ""
	return m
}

func (m *MyDb) OrWhere(query string, args ...interface{}) *MyDb {
	if query != "" {
		if m.Dialect().GetName() != "oci8" {
			query = strings.Replace(query, "\"", " ", -1)
		}
		m.DB = m.DB.Where(query, args...)
	}
	return m
}
func (m *MyDb) CheckExist(query string, args ...interface{}) error {
	var count int
	if m.Dialect().GetName() != "oci8" {
		query = strings.Replace(query, "\"", " ", -1)
	}
	r := m.Where(query, args...)
	err := r.Count(&count)
	if err != nil {
		return err
	} else if count > 0 {
		return errors.Errorf("The table %T with value{%s} already exist", r.Value, r.whereConditions)
	}
	return nil
}

//filter integer and string or pointer
func (m *MyDb) FilterBy(key string, val interface{}) *MyDb {
	if !IsNull(val) {
		if m.Dialect().GetName() != "oci8" {
			return m.Where(fmt.Sprintf(" %s  = ?", key), val)
		}
		return m.Where(fmt.Sprintf("\"%s\" = ?", key), val)
	}
	return m
}

func (m *MyDb) FilterByLike(key, val string) *MyDb {
	var query string
	if val != "" {
		if strings.Contains(val, "%") {
			val = strings.Replace(val, "%", "\\%", -1)
		} else if strings.Contains(val, "_") {
			val = strings.Replace(val, "_", "\\_", -1)
		}
		if m.Dialect().GetName() != "oci8" {
			query = fmt.Sprintf(" %s  LIKE ?", key)
		} else {
			query = fmt.Sprintf("\"%s\" LIKE ?", key)
		}
		return m.Where(query, "%"+val+"%")
	}
	return m
}

func (m *MyDb) OrFilterByLike(key, val string) *MyDb {
	var query string
	if val != "" {
		if strings.Contains(val, "%") {
			val = strings.Replace(val, "%", "\\%", -1)
		} else if strings.Contains(val, "_") {
			val = strings.Replace(val, "_", "\\_", -1)
		}
		if m.Dialect().GetName() != "oci8" {
			query = fmt.Sprintf(" %s  LIKE ?", key)
		} else {
			query = fmt.Sprintf("\"%s\" LIKE ?", key)
		}
		m.DB = m.DB.Or(query, "%"+val+"%")
	}
	return m
}

func (m *MyDb) FilterByList(key, name string) *MyDb {
	if name != "" {
		return m.FilterByLike(key, fmt.Sprintf(`"%s"`, name))
	}
	return m
}

func getError(method, condition string, err error) error {
	if err != nil {
		return errors.Wrap(err, method+" "+condition)
	}
	return nil
}

func (m *MyDb) Last(out interface{}) error {
	return getError("Last", m.whereConditions, m.DB.Last(out).Error)
}

func (m *MyDb) Offset(offset int) *MyDb {
	if offset != 0 {
		m.DB = m.DB.Offset(offset)
	}
	return m
}

func (m *MyDb) Limit(limit int) *MyDb {
	if limit != 0 {
		m.DB = m.DB.Limit(limit)
	}
	return m
}

func (m *MyDb) Order(cond string, isAsc int) *MyDb {
	if cond != "" {
		if m.Dialect().GetName() != "oci8" {
			cond = strings.Replace(cond, "\"", " ", -1)
		}
		if isAsc == 0 {
			cond += " asc"
		} else {
			cond += " desc"
		}
		m.DB = m.DB.Order(cond)
	}
	return m
}

func (m *MyDb) Find(out interface{}) error {
	return getError("Find", m.whereConditions, m.DB.Find(out).Error)
}

func (m *MyDb) First(out interface{}) error {

	return getError("First", m.whereConditions, m.DB.First(out).Error)
}

func (m *MyDb) Count(value interface{}) error {
	return getError("Count", m.whereConditions, m.DB.Count(value).Error)
}

func (m *MyDb) OrmCount(value interface{}) *MyDb {
	m.DB = m.DB.Count(value)
	return m
}

func (m *MyDb) Delete(value interface{}) error {
	return getError("Delete", m.whereConditions, m.DB.Delete(value).Error)
}

func (m *MyDb) Updates(values interface{}) error {
	return getError("Updates", m.whereConditions, m.DB.Updates(values).Error)
}

func (m *MyDb) Update(key string, value interface{}) error {
	return getError("Update", m.whereConditions, m.DB.Update(key, value).Error)
}

func GetOne(db *MyDb, out interface{}, query string, args ...interface{}) error {
	if db.Dialect().GetName() != "oci8" {
		query = strings.Replace(query, "\"", " ", -1)
	}
	return Model(db, nil).Where(query, args...).First(out)
}

func DeleteOne(db *MyDb, value interface{}, query string, args ...interface{}) error {
	if db.Dialect().GetName() != "oci8" {
		query = strings.Replace(query, "\"", " ", -1)
	}
	return Model(db, nil).Where(query, args...).Delete(value)
}

func InsertAndGetDb(db *MyDb, data interface{}) error {
	return db.Create(data).First(data).Error
}

func InsertDb(db *MyDb, data interface{}) error {
	return db.Create(data).Error
}

//根据data类型的主键，更新某个指定的属性
func UpdateAttr(db *MyDb, data interface{}, key string, value interface{}) error {
	return Model(db, data).Update(key, value)
}

func UpdateAttrs(db *MyDb, data interface{}, values interface{}) error {
	return Model(db, data).Updates(values)
}

func SaveAttrs(db *MyDb, data interface{}) error {
	return db.Save(data).Error
}

//生成随机字符串
func GetRandomString(size int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成SubNet token
func GenerateToken(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func CreatUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func RecordNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error())
}

//获取计数(increment大于0时，每次获取加increment大小的数字)
func getCounter(db *MyDb, key string, increment uint) (uint, error) {
	index_lock.Lock()
	defer index_lock.Unlock()
	var nc NumCounter
	err := GetOne(db, &nc, "\"id\" = ?", key)
	if increment == 0 {
		return nc.Index, nil
	}
	if RecordNotFound(err) {
		nc.ID = key
		nc.Index = increment
		if err = InsertDb(db, &nc); err != nil {
			err = errors.Wrap(err, "Insert num counter failed")
		}
	} else if err != nil {
		return 0, errors.Wrap(err, "Get num counter failed")
	} else {
		nc.Index += increment
		err = UpdateAttr(db, &nc, "index", nc.Index)
	}
	return nc.Index, err
}

func UpdateTable(db *MyDb, keys ...string) error {
	_, err := getCounter(db, "table_"+strings.Join(keys, "_"), 1)
	return err
}

func GetTableNum(db *MyDb, keys ...string) (uint, error) {
	return getCounter(db, "table_"+strings.Join(keys, "_"), 0)
}

func GetMsp(uuid string) string {
	return strings.Replace(uuid, ".", "MSP", -1)
}

func ReSetUUID(mspId string) string {
	return strings.Replace(mspId, "MSP", ".", -1)
}

// OpenDatabase create connect with mysql
func OpenDatabase(port int, host, user, dbname, password, Type string) (*MyDb, error) {
	var err error
	db := &MyDb{}
	var dbDSN string
	if Type == "postgres" {
		dbDSN = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			host, port, user, dbname, password)
	} else if Type == "mysql" {
		dbDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
			user, password, host, port, dbname)
	} else if Type == "oracle" {
		dbDSN = fmt.Sprintf("%s/%s@%s:%d/XE",
			user, password, host, port)
		Type = "oci8"
	} else {
		errInfo := fmt.Sprintf("This database is not supported ,Type = %s", Type)
		logger.Fatal(errInfo)
		return nil, errors.New(errInfo)
	}
	db.DB, err = gorm.Open(Type, dbDSN)
	if err != nil {
		logger.Error(err)
		logger.Fatal("Initlization database connection error.")
		return db, err
	}

	err = db.DB.DB().Ping()
	if err != nil {
		return db, err
	}
	// Disable table name's pluralization globally
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`, table name setted with `TableName` won't be affected
	db.DB.DB().SetMaxIdleConns(10)
	db.DB.DB().SetMaxOpenConns(10)

	return db, nil
}

func (m *MyDb) Migrate(tb interface{}) error {
	return m.DB.AutoMigrate(tb).Error
}

//monitorning
func (m *MyDb) HasTable(tbname string) bool {
	return m.DB.HasTable(tbname)
}

func (m *MyDb) Table(tbname string) *MyDb {
	m.DB = m.DB.Table(tbname)
	return m
}

func (m *MyDb) CreateTable(t interface{}) error {
	return m.DB.CreateTable(t).Error
}

func (m *MyDb) DropTable(t interface{}) error {
	return m.DB.DropTable(t).Error
}
