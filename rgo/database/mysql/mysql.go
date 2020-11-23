package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"rgo/rgo/viper"
	"sync"
)

var (
	mdb       *gorm.DB
	mysqlOnce sync.Once
	mErr      error
)

func Mysql() (*gorm.DB, error) {
	//单例
	mysqlOnce.Do(func() {
		config := viper.NewConfig("database", "json")
		user := config.GetString("mysql.db_user")
		host := config.GetString("mysql.db_host")
		port := config.GetString("mysql.db_port")
		password := config.GetString("mysql.db_password")
		name := config.GetString("mysql.db_name")
		//链接mysql
		dns := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"
		mdb, mErr = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dns,
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
	})

	return mdb, mErr
}
