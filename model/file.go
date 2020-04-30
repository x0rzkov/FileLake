package model

import (
	"FileLake/common/logger"
	"fmt"
	"go.uber.org/zap"
)

//******************************************************************//
//**File 类型定义
//******************************************************************//
type File struct {
	Id             int64  `xorm:"pk autoincr comment('自增ID') BIGINT(20)" json:"id"`
	ServiceLink    string `xorm:"service_link char(64) not null" json:"service_link"`
	SeaweedfsId    string `xorm:"seaweedfs_id varchar(64)" json:"seaweedfs_id"`
	AccountAddress string `xorm:"account_address varchar(64)" json:"account_address"`
	ExpireTime     int    `xorm:"expire_time  BIGINT(20)" json:"expire_time"`
}

//******************************************************************//
//** 同步表结构
//******************************************************************//

func (t *File) SyncTable() error {
	db := GetXOrmInstance()
	if db == nil {
		return fmt.Errorf("XOrmInstance is null.")
	}
	return db.Sync2(t)
}

//******************************************************************//
//** 追加一条记录
//******************************************************************//
func (t *File) Add() error {
	db := GetXOrmInstance()
	if db == nil {
		return fmt.Errorf("XOrmInstance is null.")
	}
	t.Id = 0
	_, err := db.InsertOne(t)

	return err
}

//******************************************************************//
//** 取得一条记录
//******************************************************************//
func (t *File) GetByServiceLink() (*File, error) {
	db := GetXOrmInstance()
	if db == nil {
		return nil, fmt.Errorf("XOrmInstance is null.")
	}
	out := File{}
	has, err := db.Where("service_link = ?", t.ServiceLink).Get(&out)
	if err != nil {
		logger.Error("err", zap.Any("err", err))
		return nil, err
	}
	if !has {
		err = fmt.Errorf("No proper record")
		logger.Error("err", zap.Any("err", err))
		return nil, err
	}
	return &out, nil
}

//******************************************************************//
//** 删除一条记录
//******************************************************************//
func (t *File) DeleteByServiceLink() error {
	db := GetXOrmInstance()
	if db == nil {
		return fmt.Errorf("XOrmInstance is null.")
	}
	out := File{}
	has, err := db.Where("service_link = ?", t.ServiceLink).Get(&out)
	if err != nil {
		logger.Error("err", zap.Any("err", err))
		return err
	}
	if !has {
		err = fmt.Errorf("No proper record")
		logger.Error("err", zap.Any("err", err))
		return err
	}
	_, err = db.Delete(out)
	if err != nil {
		logger.Error("err", zap.Any("err", err))
		return err
	}
	return nil
}
//******************************************************************//
//** 取得多条记录
//******************************************************************//
func (t *File) GetbyFid() ([]*File, error) {
	out := make([]*File, 0)
	db := GetXOrmInstance()
	if db == nil {
		return out, fmt.Errorf("XOrmInstance is null.")
	}
	err := db.Where("service_link = ?", t.ServiceLink).Find(&out)
	if err != nil {
		logger.Error("err", zap.Any("err", err))
	}
	return out, err
}
