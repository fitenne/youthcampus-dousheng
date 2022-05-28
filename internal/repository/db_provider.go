package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Driver         string
	Host, Port     string
	User, Password string
	DBname         string
	LogLevel       int
}

type MysqlProdiver struct {
	db          *gorm.DB
	connectOnce sync.Once
}

type DBProvider interface {
	Connect(c DBConfig) error
	GetDB() *gorm.DB
}

var dbProvider DBProvider
var initOnce sync.Once

//可以在db_provider.go定义这个变量，顺带在init函数初始化db
var db *gorm.DB

// 连接到 DBConfig 制定的数据库，忽略 DBConfig 中的 Driver 字段
func (p *MysqlProdiver) Connect(c DBConfig) error {
	err := errors.New("already connected")
	p.connectOnce.Do(func() {
		template := "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf(template, c.User, c.Password, c.Host, c.Port, c.DBname)
		dialector := mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        dsn,
		})
		p.db, err = gorm.Open(dialector, &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				LogLevel: logger.LogLevel(c.LogLevel),
			}),
		})
	})
	return err
}

func (p *MysqlProdiver) GetDB() *gorm.DB {
	return p.db
}

// 初始化数据库，只有第一次调用有效
func Init(c DBConfig) error {
	err := errors.New("Init called twice")
	initOnce.Do(func() {
		switch c.Driver {
		case "mysql":
			dbProvider = &MysqlProdiver{}
		default:
			err = errors.New("db driver not supported")
		}
		err = dbProvider.Connect(c)
		if err != nil {
			return
		}
		//在init中顺便对db进行了初始化
		db = dbProvider.GetDB()
		// 创建表user
		//if !db.Migrator().HasTable(&model.User{}) {
		//	if err := db.AutoMigrate(new(model.User)).Error; err != nil {
		//		panic(err)
		//	}
		//}
		// 创建表video
		//if !db.Migrator().HasTable(&model.Video{}) {
		//	if err := db.AutoMigrate(new(model.Video)).Error; err != nil {
		//		panic(err)
		//	}
		//}


		//// 创建表video
		//if !db.Migrator().HasTable(&model.Video{}) {
		//	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(model.Video{}).Error; err != nil {
		//		panic(err)
		//	}
		//}
		err = nil
	})
	return err
}
