package service

import (
	"FileLake/model"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/linxGnu/goseaweedfs"
	"io"
	"log"
	"net/http"
)

type Result struct {
	Code    int
	Message string
	Data    interface{}
}

var FailResult Result = Result{Code: 2, Message: "错误", Data: nil}

// @Summary 新增文件
// @Success 200 object service.Result 成功后返回值
// @Router /add/{accountaddress} [post]
func Add(c *gin.Context) {
	accountAddress := c.Param("accountaddress")
	file, _ := c.FormFile("file")
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		println(err.Error())
	}
	// step 1, save file on seaweedfs
	fp, err := model.SaveSeaweedfs(fileReader, file.Filename, file.Size, "col","")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		println(err.Error())
	}

	// step 2, save fs in DB
	fileId := fp.FileID
	b := make([]byte, 6)
	rand.Read(b)
	serviceLink := base64.StdEncoding.EncodeToString(b)

	fileRecord := model.File{serviceLink, fileId, accountAddress, 9999999999}
	model.DB.Create(&fileRecord)

	// step 3, return post json
	jsonString, _ := json.Marshal(fileRecord)
	c.JSON(http.StatusOK, jsonString)
}




// @Summary 更新文件信息
// @Success 200 object service.Result 成功后返回值
// @Router /user [put]
func Update(c *gin.Context, serviceLink string, seaweedfsId string, accountAddress string, expireTime int) {
	file := model.File{serviceLink, seaweedfsId, accountAddress, expireTime}
	if err := c.ShouldBindJSON(&file); err == nil {
		log.Println(file)
		file.Update()
		c.JSON(http.StatusOK, file)
	} else {
		log.Println(file)
		c.JSON(http.StatusBadRequest, FailResult)
	}
}

// @Summary 获取用户信息
// @Success 200 object service.Result 成功后返回值
// @Router /getfilebylink/{servicelink}
func GetFileByLink(c *gin.Context) {
	file := model.File{}
	serviceLink := c.Param("serviceLink")
	f, err := file.GetFileByLink(serviceLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{2, "参数错误", err.Error()})
		return file, err
	}
	c.JSON(http.StatusOK, f)
	return f, nil
}

// @Summary 获取用户信息
// @Success 200 object service.Result 成功后返回值
// @Router /user [get]
func GetUserFile(c *gin.Context, seaweedfsId string, accountAddress string) (model.File, error) {
	file := model.File{}
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, Result{2, "参数错误", err.Error()})
		return file, err
	}
	f, err := file.GetUserFile(seaweedfsId, accountAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, Result{2, "参数错误", err.Error()})
		return file, err
	}
	c.JSON(http.StatusOK, f)
	return f, nil
}

// @Summary 删除文件信息
// @Success 200 object service.Result 成功后返回值
// @Router /user [delete]
func Delete(context *gin.Context, serviceLink string) {
	file, _ := GetFileByLink(context, serviceLink)
	if err := context.ShouldBindJSON(&file); err == nil {
		log.Println(file)
		file.DeleteUser(serviceLink)
		context.JSON(http.StatusOK, file)
	} else {
		log.Println(file)
		context.JSON(http.StatusBadRequest, FailResult)
	}
}
