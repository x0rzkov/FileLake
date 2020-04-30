package db

import (
	"HIBL/common/config"
	"HIBL/common/logger"

	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestGetMenuList(t *testing.T) {
	defer logger.Sync()
	logger.Init("test")
	logger.SetLogLevel(zapcore.DebugLevel)
	logger.Debug("TestGetByList start>>>>>>>>>>>>>>>>>>>>>>>>")
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.DBuser, "123456", "192.168.0.210:3306", config.DBname)
	logger.Info("conn", zap.Any("conn", conn))
	Init(conn)

	mylist, err := GetMenuList(0)
	if err != nil {
		logger.Debug("e", zap.Any("err", err))
		return
	}
	for i, menu := range mylist {
		fmt.Println("i = ", i)
		fmt.Println("menu = ", menu)
	}
	logger.Debug("TestGetByList end<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
func TestGetRoles(t *testing.T) {
	defer logger.Sync()
	logger.Init("test")
	logger.SetLogLevel(zapcore.DebugLevel)
	logger.Debug("TestGetByList start>>>>>>>>>>>>>>>>>>>>>>>>")
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.DBuser, "123456", "192.168.0.210:3306", config.DBname)
	logger.Info("conn", zap.Any("conn", conn))
	Init(conn)
	out, _ := GetRoles("11111111", 0, "")
	fmt.Println(out)
	out, _ = GetRoles("22222222", 0, "")
	fmt.Println(out)
	out, _ = GetRoles("33333333", 0, "")
	fmt.Println(out)
}
