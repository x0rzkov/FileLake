package model

type File struct {
	ServiceLink    string `gorm:"primary_key;column:service_link;type:char(64);not null" json:"service_link"`
	SeaweedfsId    string `gorm:"column:seaweedfs_id" json:"seaweedfs_id"`
	AccountAddress string `gorm:"column:account_address" json:"account_address"`
	ExpireTime     int    `gorm:"column:expire_time" json:"expire_time"`
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
