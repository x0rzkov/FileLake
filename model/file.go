package model

type File struct {
	Id             int64  `xorm:"pk autoincr comment('自增ID') BIGINT(20)" json:"id"`
	ServiceLink    string `xorm:"pk name:service_link char(64) not null" json:"service_link"`
	SeaweedfsId    string `xorm:"name:seaweedfs_id varchar(64)" json:"seaweedfs_id"`
	AccountAddress string `xorm:"name:account_address varchar(64)" json:"account_address"`
	ExpireTime     int    `xorm:"name:expire_time  BIGINT(20)" json:"expire_time"`
}


func (File) TableName() string {
	return "file"
}

func (file *File) Add() {
	DB.Create(file)
}

func (file *File) Update() {
	DB.Save(file)
}

func (file *File) GetUserFile(seaweedfsId string, accountAddress string) (File, error) {
	result := DB.Where("seaweedfs_id = ? AND AccountAddress = ?", seaweedfsId, accountAddress).Find(file)
	return *file, result.Error
}

func (file *File) GetFileByLink(serviceLink string) (File, error) {
	result := DB.First(file, serviceLink)
	return *file, result.Error
}

func (file *File) DeleteUser(serviceLink string) {
	DB.Delete(file, serviceLink)
}
