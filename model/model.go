package model

import (
	"FileLake/common/logger"
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

//var DB *gorm.DB
//
//func init() {
//	var conn = config.Config.GetString("sqlConn")
//	db, err := gorm.Open("mysql", conn)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//设置连接池
//	//空闲
//	db.DB().SetMaxIdleConns(50)
//	//打开
//	db.DB().SetMaxOpenConns(100)
//	//超时
//	db.DB().SetConnMaxLifetime(time.Second * 30)
//	DB = db
//	fmt.Println("db init succeed")
//	DB.AutoMigrate(&File{})
//
//}


// 数据库实例
var (
	xOrmInstance *xorm.Engine
)

// Init 初始化
func Init(conn string) {

	logger.Debug("DB Init start>>>>>>>>")
	logger.Debug(conn)
	xOrmInstance = GetDBEngine(conn)
	logger.Debug("DB Init end  <<<<<<<<")
}

func GetXOrmInstance() *xorm.Engine {
	return xOrmInstance
}

func GetDBEngine(conn string) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", conn)
	if err != nil {
		logger.Fatal("数据库初始化错误", zap.Any("error", err))
		return nil
	}

	engine.SetLogger(logger.GetXormLogger())

	engine.ShowSQL()

	return engine
}