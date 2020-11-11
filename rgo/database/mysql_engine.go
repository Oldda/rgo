package database

import(
	"rgo/rgo/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlEngine *gorm.DB

func NewMysqlEngine()error{
	var err error
	config := viper.NewConfig("database","json")
	db_user := config.GetString("mysql.db_user")
	db_host := config.GetString("mysql.db_host")
	db_port := config.GetString("mysql.db_port")
	db_password := config.GetString("mysql.db_password")
	db_name := config.GetString("mysql.db_name")

	dns := db_user + ":"+ db_password + "@tcp("+ db_host +":" + db_port + ")/" + db_name + "?charset=utf8&parseTime=True&loc=Local"
	MysqlEngine, err = gorm.Open(mysql.New(mysql.Config{
		DSN:dns,
		DefaultStringSize: 256, // string 类型字段的默认长度
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	 	DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	  	DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	  	SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}),&gorm.Config{})
  	return err
}