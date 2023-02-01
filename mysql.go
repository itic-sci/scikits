package scikits

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func NewMysqlClient(label string) *gorm.DB {
	var _db *gorm.DB
	//配置MySQL连接参数
	username := MyViper.GetString(fmt.Sprintf("%s.user", label))
	password := MyViper.GetString(fmt.Sprintf("%s.pass", label))
	host := MyViper.GetString(fmt.Sprintf("%s.host", label))
	port := MyViper.GetString(fmt.Sprintf("%s.port", label))
	Dbname := MyViper.GetString(fmt.Sprintf("%s.db", label))

	//MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local", username, password, host, port, Dbname)

	//连接MYSQL
	var err error
	_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	sqlDB, _ := _db.DB()
	sqlDB.SetMaxOpenConns(100)          //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)           //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	sqlDB.SetConnMaxLifetime(time.Hour) // SetConnMaxLifetime 设置了连接可复用的最大时间。

	return _db
}

//定义全局的db对象，我们执行数据库操作主要通过他实现。
//获取gorm db对象，其他包需要执行数据库查询的时候，只要通过tools.getDB()获取db对象即可。
//不用担心协程并发使用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接
//var _db *gorm.DB
//func init() {
//	_db = NewMysqlClient()
//}
//func GetMysqlDB() *gorm.DB {
//	return _db
//}
